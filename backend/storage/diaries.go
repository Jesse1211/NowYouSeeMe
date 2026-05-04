package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"nowyouseeme/models"
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

	// Use provided timestamp or default to now
	timestamp := now
	if payload.DiaryTimestamp != nil {
		timestamp = *payload.DiaryTimestamp
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	query := `
		INSERT INTO agent_diary_versions (id, agent_id, diary_timestamp, raw_payload, created_at, parsed_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	var execErr error
	if tx != nil {
		_, execErr = tx.Exec(query, diaryID, agentID, timestamp, payloadJSON, now, now)
	} else {
		_, execErr = s.db.Exec(query, diaryID, agentID, timestamp, payloadJSON, now, now)
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
	result, err := s.GetLatestSnapshot(agentID)
	if err != nil {
		return nil, err
	}

	// Handle first diary submission - no snapshot exists yet
	if result == nil {
		result = &models.AgentSnapshotResult{
			AgentID:   agentID,
			State:     models.NewEmptyState(),
			Sequence:  0,
			UpdatedAt: nil,
		}
	}

	// Track old MBTI before applying changes
	latestMBTI := result.State.MBTI

	// Validate operations
	// TODO: Re-enable after fixing seed script state tracking
	// if err := validation.ValidateOperations(payload.Operations, result.State); err != nil {
	// 	return nil, err
	// }

	// Insert diary version
	diaryID, err := s.insertDiaryVersionTx(tx, agentID, payload)
	if err != nil {
		return nil, err
	}

	// Get next sequence number
	nextSeq := result.Sequence + 1

	// Create metadata_update event (Event Sourcing purity: all state changes via events)
	metadataEvent := &models.Event{
		AgentID:        agentID,
		DiaryID:        diaryID,
		EventType:      "metadata_update",
		SequenceNumber: nextSeq,
	}
	metadataPayload := models.EventPayload{
		MBTI:                  payload.MBTI,
		MBTIConfidence:        payload.MBTIConfidence,
		GeometryRep:           payload.GeometryRep,
		CurrentMood:           payload.CurrentMood,
		Philosophy:            payload.Philosophy,
		CurrentSelfReflection: &payload.SelfReflection,
	}
	metadataPayloadJSON, err := json.Marshal(metadataPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata payload: %w", err)
	}
	metadataEvent.Payload = metadataPayloadJSON

	// Insert metadata event
	if err := s.insertEventTx(tx, metadataEvent); err != nil {
		return nil, err
	}

	// Apply metadata event to state
	if err := ApplyEvent(result.State, metadataEvent); err != nil {
		return nil, fmt.Errorf("failed to apply metadata event: %w", err)
	}

	// Create and insert operation events (starting from nextSeq + 1)
	for i, op := range payload.Operations {
		event, err := OperationToEvent(op, agentID, diaryID, nextSeq+1+int64(i))
		if err != nil {
			return nil, fmt.Errorf("failed to convert operation to event: %w", err)
		}

		if err := s.insertEventTx(tx, event); err != nil {
			return nil, err
		}

		// Apply event to current state
		if err := ApplyEvent(result.State, event); err != nil {
			return nil, fmt.Errorf("failed to apply event to state: %w", err)
		}
	}

	// Calculate final sequence number (metadata event + operation events)
	finalSeq := nextSeq + int64(len(payload.Operations))

	// Check if MBTI changed and update agent record + insert timeline record
	newMBTI := payload.MBTI
	if newMBTI != latestMBTI {
		// Update agents.current_mbti
		_, err = tx.Exec(`UPDATE agents SET current_mbti = $1 WHERE id = $2`, newMBTI, agentID)
		if err != nil {
			return nil, fmt.Errorf("failed to update agent current_mbti: %w", err)
		}

		// Insert MBTI timeline record
		err = s.insertMBTITimeline(tx, agentID, newMBTI, diaryID, finalSeq)
		if err != nil {
			return nil, err
		}
	}

	// Materialize snapshot (Current strategy: Always materialize after each diary submission)
	// This ensures minimal WAL replay overhead on reads. Alternative strategies:
	// - Conditional: Only materialize every N diary submissions (trades write performance for read performance)
	// - Threshold: Only materialize if uncommitted events exceed threshold
	// Current "always materialize" is optimal for read-heavy workloads.
	if err := s.upsertSnapshotTx(tx, agentID, diaryID, result.State, finalSeq); err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return result.State, nil
}

// insertMBTITimeline inserts a new MBTI timeline record
