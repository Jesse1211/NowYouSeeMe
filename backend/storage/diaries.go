package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"nowyouseeme/models"
	"nowyouseeme/validation"
	"time"

	"github.com/google/uuid"
)

func (s *PostgresStore) InsertDiaryVersion(agentID string, payload *models.DiaryPayload) (string, error) {
	return s.insertDiaryVersionTx(nil, agentID, payload)
}

// insertDiaryVersionTx creates a diary record within a transaction or using default connection
func (s *PostgresStore) insertDiaryVersionTx(tx *sql.Tx, agentID string, payload *models.DiaryPayload) (string, error) {
	diaryID := "diary_" + uuid.New().String()
	now := time.Now()

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	query := `
		INSERT INTO agent_diary_versions (id, agent_id, raw_payload, created_at)
		VALUES ($1, $2, $3, $4)
	`

	var execErr error
	if tx != nil {
		_, execErr = tx.Exec(query, diaryID, agentID, payloadJSON, now)
	} else {
		_, execErr = s.db.Exec(query, diaryID, agentID, payloadJSON, now)
	}

	if execErr != nil {
		return "", fmt.Errorf("failed to insert diary: %w", execErr)
	}

	return diaryID, nil
}

// InsertEvent creates a new event
func (s *PostgresStore) AcquireLock(tx *sql.Tx, agentID string) error {
	// Always lock the agents table row for consistency
	// This ensures uniform locking behavior regardless of snapshot existence
	query := `
		SELECT id
		FROM agents
		WHERE id = $1
		FOR UPDATE
	`

	var lockedID string
	err := tx.QueryRow(query, agentID).Scan(&lockedID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("agent not found: %s", agentID)
		}
		return fmt.Errorf("failed to acquire lock: %w", err)
	}

	return nil
}

// SubmitDiary handles complete diary submission with transaction
func (s *PostgresStore) SubmitDiary(agentID string, payload *models.DiaryPayload) (*models.AgentState, error) {
	// Begin transaction with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Acquire pessimistic lock
	if err := s.AcquireLock(tx, agentID); err != nil {
		return nil, err
	}

	// Build current state using WAL
	latestSnapshot, err := s.GetLatestSnapshot(agentID)
	if err != nil {
		return nil, err
	}

	// Handle first diary submission - no snapshot exists yet
	if latestSnapshot == nil {
		latestSnapshot = &models.AgentSnapshotResult{
			AgentID:           agentID,
			State:             models.NewEmptyState(),
			LastEventSequence: 0,
			UpdatedAt:         nil,
		}
	}

	// Track old MBTI before applying changes
	latestMBTI := latestSnapshot.State.MBTI

	// Validate operations
	if err := validation.ValidateOperations(payload.Operations, latestSnapshot.State); err != nil {
		return nil, err
	}

	// Insert diary version
	diaryID, err := s.insertDiaryVersionTx(tx, agentID, payload)
	if err != nil {
		return nil, err
	}

	nextSeq := latestSnapshot.LastEventSequence

	// Insert metadata event first (doesn't participate in AgentState replay)
	nextSeq++
	metadataEvent, err := MetadataToEventConverter(payload, agentID, diaryID, nextSeq)
	if err != nil {
		return nil, fmt.Errorf("failed to create metadata event: %w", err)
	}
	if err := s.insertEventTx(tx, metadataEvent); err != nil {
		return nil, fmt.Errorf("failed to insert metadata event: %w", err)
	}

	// Create and insert operation events (participate in AgentState replay)
	for _, op := range payload.Operations {
		nextSeq++
		event, err := OperationToEventConverter(op, agentID, diaryID, nextSeq)
		if err != nil {
			return nil, fmt.Errorf("failed to convert operation to event: %w", err)
		}

		if err := s.insertEventTx(tx, event); err != nil {
			return nil, err
		}

		// Apply event to current state (metadata events are skipped in ApplyEventToSnapshot)
		if err := ApplyEventToSnapshot(latestSnapshot, event); err != nil {
			return nil, fmt.Errorf("failed to apply event to state: %w", err)
		}
	}

	// Check if MBTI changed and update agent record + insert timeline record
	newMBTI := payload.MBTI
	if newMBTI != latestMBTI {
		// Update agents.current_mbti
		_, err = tx.Exec(`UPDATE agents SET current_mbti = $1 WHERE id = $2`, newMBTI, agentID)
		if err != nil {
			return nil, fmt.Errorf("failed to update agent current_mbti: %w", err)
		}

		// Insert MBTI timeline record
		err = s.insertMBTITimeline(tx, agentID, newMBTI, diaryID, nextSeq)
		if err != nil {
			return nil, err
		}
	}

	// Materialize snapshot (Current strategy: Always materialize after each diary submission)
	// This ensures minimal WAL replay overhead on reads. Alternative strategies:
	// - Conditional: Only materialize every N diary submissions (trades write performance for read performance)
	// - Threshold: Only materialize if uncommitted events exceed threshold
	// Current "always materialize" is optimal for read-heavy workloads.
	if err := s.upsertSnapshotTx(tx, agentID, latestSnapshot.State, nextSeq); err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return latestSnapshot.State, nil
}

// insertMBTITimeline inserts a new MBTI timeline record
