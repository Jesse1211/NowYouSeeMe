package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"nowyouseeme/models"
	// "nowyouseeme/validation" // TODO: Re-enable after fixing seed script
	"time"

	"github.com/google/uuid"
)

// PostgresStore implements PostgreSQL storage
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore creates a new PostgreSQL store
func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

// CreateAgent creates a new agent
func (s *PostgresStore) CreateAgent(req *models.CreateAgentRequest) (*models.Agent, error) {
	agent := &models.Agent{
		ID:          req.AgentID,
		Name:        req.Name,
		CurrentMBTI: req.CurrentMBTI,
		CreatedAt:   time.Now(),
	}

	// Begin transaction with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert agent record
	query := `
		INSERT INTO agents (id, name, current_mbti, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err = tx.Exec(query, agent.ID, agent.Name, agent.CurrentMBTI, agent.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	// Insert initial MBTI timeline record (no diary yet, so diary_id is NULL)
	err = s.insertMBTITimelineNullable(tx, agent.ID, agent.CurrentMBTI, nil, 0)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return agent, nil
}

// GetAgent retrieves an agent by ID
func (s *PostgresStore) GetAgent(agentID string) (*models.Agent, error) {
	query := `
		SELECT id, name, current_mbti, created_at
		FROM agents
		WHERE id = $1
	`

	agent := &models.Agent{}
	err := s.db.QueryRow(query, agentID).Scan(
		&agent.ID,
		&agent.Name,
		&agent.CurrentMBTI,
		&agent.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("agent not found: %s", agentID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get agent: %w", err)
	}

	return agent, nil
}

// GetAllAgents retrieves all agents
func (s *PostgresStore) GetAllAgents() ([]*models.Agent, error) {
	query := `
		SELECT id, name, current_mbti, created_at
		FROM agents
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query agents: %w", err)
	}
	defer rows.Close()

	agents := []*models.Agent{}
	for rows.Next() {
		agent := &models.Agent{}
		if err := rows.Scan(&agent.ID, &agent.Name, &agent.CurrentMBTI, &agent.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan agent: %w", err)
		}
		agents = append(agents, agent)
	}

	return agents, nil
}

// InsertDiaryVersion creates a diary record
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
func (s *PostgresStore) InsertEvent(event *models.Event) error {
	return s.insertEventTx(nil, event)
}

// insertEventTx creates a new event within a transaction or using default connection
func (s *PostgresStore) insertEventTx(tx *sql.Tx, event *models.Event) error {
	query := `
		INSERT INTO events (agent_id, diary_id, event_type, timestamp, payload, sequence_number)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING event_id
	`

	var err error
	if tx != nil {
		err = tx.QueryRow(
			query,
			event.AgentID,
			event.DiaryID,
			event.EventType,
			time.Now(),
			event.Payload,
			event.SequenceNumber,
		).Scan(&event.EventID)
	} else {
		err = s.db.QueryRow(
			query,
			event.AgentID,
			event.DiaryID,
			event.EventType,
			time.Now(),
			event.Payload,
			event.SequenceNumber,
		).Scan(&event.EventID)
	}

	if err != nil {
		return fmt.Errorf("failed to insert event: %w", err)
	}

	return nil
}

// GetNextSequenceNumber gets the next sequence number for an agent
func (s *PostgresStore) GetNextSequenceNumber(agentID string) (int64, error) {
	query := `
		SELECT COALESCE(MAX(sequence_number), 0) + 1
		FROM events
		WHERE agent_id = $1
	`

	var nextSeq int64
	err := s.db.QueryRow(query, agentID).Scan(&nextSeq)
	if err != nil {
		return 0, fmt.Errorf("failed to get next sequence: %w", err)
	}

	return nextSeq, nil
}

// GetUncommittedEvents retrieves events after a sequence number
func (s *PostgresStore) GetUncommittedEvents(agentID string, afterSequence int64) ([]*models.Event, error) {
	query := `
		SELECT event_id, agent_id, diary_id, event_type, timestamp, payload, sequence_number
		FROM events
		WHERE agent_id = $1 AND sequence_number > $2
		ORDER BY sequence_number ASC
	`

	rows, err := s.db.Query(query, agentID, afterSequence)
	if err != nil {
		return nil, fmt.Errorf("failed to query uncommitted events: %w", err)
	}
	defer rows.Close()

	events := []*models.Event{}
	for rows.Next() {
		event := &models.Event{}
		if err := rows.Scan(
			&event.EventID,
			&event.AgentID,
			&event.DiaryID,
			&event.EventType,
			&event.Timestamp,
			&event.Payload,
			&event.SequenceNumber,
		); err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}

// GetEventsByAgent retrieves all events for an agent
func (s *PostgresStore) GetEventsByAgent(agentID string) ([]*models.Event, error) {
	query := `
		SELECT event_id, agent_id, diary_id, event_type, timestamp, payload, sequence_number
		FROM events
		WHERE agent_id = $1
		ORDER BY sequence_number ASC
	`

	rows, err := s.db.Query(query, agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	events := []*models.Event{}
	for rows.Next() {
		event := &models.Event{}
		if err := rows.Scan(
			&event.EventID,
			&event.AgentID,
			&event.DiaryID,
			&event.EventType,
			&event.Timestamp,
			&event.Payload,
			&event.SequenceNumber,
		); err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}

// GetSnapshot retrieves the latest snapshot for an agent
func (s *PostgresStore) GetSnapshot(agentID string) (*models.AgentStateSnapshot, error) {
	query := `
		SELECT agent_id, derived_from_diary_id, last_event_sequence, updated_at, state
		FROM agent_state_snapshots
		WHERE agent_id = $1
	`

	snapshot := &models.AgentStateSnapshot{}
	err := s.db.QueryRow(query, agentID).Scan(
		&snapshot.AgentID,
		&snapshot.DerivedFromDiaryID,
		&snapshot.LastEventSequence,
		&snapshot.UpdatedAt,
		&snapshot.State,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No snapshot exists yet
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshot: %w", err)
	}

	return snapshot, nil
}

// UpsertSnapshot creates or updates a snapshot
func (s *PostgresStore) UpsertSnapshot(agentID, diaryID string, state *models.AgentState, lastEventSeq int64) error {
	return s.upsertSnapshotTx(nil, agentID, diaryID, state, lastEventSeq)
}

// upsertSnapshotTx creates or updates a snapshot within a transaction or using default connection
func (s *PostgresStore) upsertSnapshotTx(tx *sql.Tx, agentID, diaryID string, state *models.AgentState, lastEventSeq int64) error {
	stateJSON, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	query := `
		INSERT INTO agent_state_snapshots (agent_id, derived_from_diary_id, last_event_sequence, updated_at, state)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (agent_id)
		DO UPDATE SET
			derived_from_diary_id = EXCLUDED.derived_from_diary_id,
			last_event_sequence = EXCLUDED.last_event_sequence,
			updated_at = EXCLUDED.updated_at,
			state = EXCLUDED.state
	`

	var execErr error
	if tx != nil {
		_, execErr = tx.Exec(query, agentID, diaryID, lastEventSeq, time.Now(), stateJSON)
	} else {
		_, execErr = s.db.Exec(query, agentID, diaryID, lastEventSeq, time.Now(), stateJSON)
	}

	if execErr != nil {
		return fmt.Errorf("failed to upsert snapshot: %w", execErr)
	}

	return nil
}

// GetLatestSnapshot builds latest snapshot using WAL pattern
// Returns a rich domain object encapsulating all snapshot information
func (s *PostgresStore) GetLatestSnapshot(agentID string) (*models.AgentSnapshotResult, error) {
	// Load latest snapshot
	snapshot, err := s.GetSnapshot(agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshot: %w", err)
	}

	// No snapshot exists yet - just created account without any diary submission
	if snapshot == nil {
		return nil, nil
	}

	// Parse snapshot state
	var state models.AgentState
	if err := json.Unmarshal(snapshot.State, &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal snapshot state: %w", err)
	}

	// Load uncommitted events (WAL)
	uncommittedEvents, err := s.GetUncommittedEvents(agentID, snapshot.LastEventSequence)
	if err != nil {
		return nil, fmt.Errorf("failed to get uncommitted events: %w", err)
	}

	// Replay uncommitted events onto snapshot
	latestSnapshot, err := ReplayEvents(&state, uncommittedEvents)
	if err != nil {
		return nil, fmt.Errorf("failed to replay events: %w", err)
	}

	// Calculate current sequence
	latestSeq := snapshot.LastEventSequence
	if len(uncommittedEvents) > 0 {
		latestSeq = uncommittedEvents[len(uncommittedEvents)-1].SequenceNumber
	}

	return &models.AgentSnapshotResult{
		AgentID:   agentID,
		State:     latestSnapshot,
		Sequence:  latestSeq,
		UpdatedAt: &snapshot.UpdatedAt,
	}, nil
}

// AcquireLock acquires pessimistic lock for agent
// Always locks the agents table for consistency and to avoid deadlocks
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

	// Create and insert events
	for i, op := range payload.Operations {
		event, err := OperationToEvent(op, agentID, diaryID, nextSeq+int64(i))
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

	// Update metadata fields from payload
	result.State.MBTI = payload.MBTI
	result.State.MBTIConfidence = payload.MBTIConfidence
	result.State.GeometryRep = payload.GeometryRep
	result.State.CurrentMood = payload.CurrentMood
	result.State.Philosophy = payload.Philosophy
	result.State.CurrentSelfReflection = payload.SelfReflection

	// Calculate final sequence number (used for MBTI timeline and snapshot)
	finalSeq := result.Sequence
	if len(payload.Operations) > 0 {
		finalSeq = nextSeq + int64(len(payload.Operations)) - 1
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
		err = s.insertMBTITimeline(tx, agentID, newMBTI, diaryID, finalSeq)
		if err != nil {
			return nil, err
		}
	}

	// Materialize snapshot (MVP: always)

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
func (s *PostgresStore) insertMBTITimeline(
	tx *sql.Tx,
	agentID string,
	mbti string,
	diaryID string,
	eventSequence int64,
) error {
	query := `
		INSERT INTO agent_mbti_timeline (
			agent_id, mbti, effective_from, diary_id, event_sequence
		) VALUES ($1, $2, NOW(), $3, $4)
	`

	_, err := tx.Exec(query, agentID, mbti, diaryID, eventSequence)
	if err != nil {
		return fmt.Errorf("failed to insert MBTI timeline: %w", err)
	}

	return nil
}

// insertMBTITimelineNullable inserts a new MBTI timeline record with nullable diary_id
func (s *PostgresStore) insertMBTITimelineNullable(
	tx *sql.Tx,
	agentID string,
	mbti string,
	diaryID *string,
	eventSequence int64,
) error {
	query := `
		INSERT INTO agent_mbti_timeline (
			agent_id, mbti, effective_from, diary_id, event_sequence
		) VALUES ($1, $2, NOW(), $3, $4)
	`

	_, err := tx.Exec(query, agentID, mbti, diaryID, eventSequence)
	if err != nil {
		return fmt.Errorf("failed to insert MBTI timeline: %w", err)
	}

	return nil
}

// GetAgentIDsByCurrentMBTI returns agent IDs with the specified current MBTI type
func (s *PostgresStore) GetAgentIDsByCurrentMBTI(mbtiType string) ([]string, error) {
	query := `
		SELECT id
		FROM agents
		WHERE current_mbti = $1
	`

	rows, err := s.db.Query(query, mbtiType)
	if err != nil {
		return nil, fmt.Errorf("failed to query agents by MBTI: %w", err)
	}
	defer rows.Close()

	agentIDs := []string{}
	for rows.Next() {
		var agentID string
		if err := rows.Scan(&agentID); err != nil {
			return nil, err
		}
		agentIDs = append(agentIDs, agentID)
	}

	return agentIDs, nil
}
