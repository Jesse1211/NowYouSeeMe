package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"nowyouseeme/models"
	"nowyouseeme/validation"
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
		InitialMBTI: req.InitialMBTI,
		CreatedAt:   time.Now(),
	}

	// Begin transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert agent record
	query := `
		INSERT INTO agents (id, name, initial_mbti, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err = tx.Exec(query, agent.ID, agent.Name, agent.InitialMBTI, agent.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	// Insert initial MBTI timeline record
	err = s.insertMBTITimeline(tx, agent.ID, agent.InitialMBTI, "initial", 0)
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
		SELECT id, name, initial_mbti, created_at
		FROM agents
		WHERE id = $1
	`

	agent := &models.Agent{}
	err := s.db.QueryRow(query, agentID).Scan(
		&agent.ID,
		&agent.Name,
		&agent.InitialMBTI,
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
		SELECT id, name, initial_mbti, created_at
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
		if err := rows.Scan(&agent.ID, &agent.Name, &agent.InitialMBTI, &agent.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan agent: %w", err)
		}
		agents = append(agents, agent)
	}

	return agents, nil
}

// InsertDiaryVersion creates a diary record
func (s *PostgresStore) InsertDiaryVersion(agentID string, payload *models.DiaryPayload) (string, error) {
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

	_, err = s.db.Exec(query, diaryID, agentID, timestamp, payloadJSON, now, now)
	if err != nil {
		return "", fmt.Errorf("failed to insert diary: %w", err)
	}

	return diaryID, nil
}

// InsertEvent creates a new event
func (s *PostgresStore) InsertEvent(event *models.Event) error {
	query := `
		INSERT INTO events (agent_id, diary_id, event_type, timestamp, payload, sequence_number)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING event_id
	`

	err := s.db.QueryRow(
		query,
		event.AgentID,
		event.DiaryID,
		event.EventType,
		time.Now(),
		event.Payload,
		event.SequenceNumber,
	).Scan(&event.EventID)

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

	_, err = s.db.Exec(query, agentID, diaryID, lastEventSeq, time.Now(), stateJSON)
	if err != nil {
		return fmt.Errorf("failed to upsert snapshot: %w", err)
	}

	return nil
}

// GetCurrentState builds current state using WAL pattern
func (s *PostgresStore) GetCurrentState(agentID string) (*models.AgentState, int64, error) {
	// Load latest snapshot
	snapshot, err := s.GetSnapshot(agentID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get snapshot: %w", err)
	}

	// No snapshot exists yet - first diary submission
	if snapshot == nil {
		return models.NewEmptyState(), 0, nil
	}

	// Parse snapshot state
	var state models.AgentState
	if err := json.Unmarshal(snapshot.State, &state); err != nil {
		return nil, 0, fmt.Errorf("failed to unmarshal snapshot state: %w", err)
	}

	// Load uncommitted events (WAL)
	uncommittedEvents, err := s.GetUncommittedEvents(agentID, snapshot.LastEventSequence)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get uncommitted events: %w", err)
	}

	// Replay uncommitted events onto snapshot
	currentState, err := ReplayEvents(&state, uncommittedEvents)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to replay events: %w", err)
	}

	// Calculate current sequence
	currentSeq := snapshot.LastEventSequence
	if len(uncommittedEvents) > 0 {
		currentSeq = uncommittedEvents[len(uncommittedEvents)-1].SequenceNumber
	}

	return currentState, currentSeq, nil
}

// AcquireLock acquires pessimistic lock for agent
func (s *PostgresStore) AcquireLock(tx *sql.Tx, agentID string) error {
	// Try to lock snapshot row
	query := `
		SELECT agent_id
		FROM agent_state_snapshots
		WHERE agent_id = $1
		FOR UPDATE
	`

	var lockedID string
	err := tx.QueryRow(query, agentID).Scan(&lockedID)

	if err == sql.ErrNoRows {
		// No snapshot exists - lock the agent row instead
		query = `
			SELECT id
			FROM agents
			WHERE id = $1
			FOR UPDATE
		`
		err = tx.QueryRow(query, agentID).Scan(&lockedID)
	}

	if err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}

	return nil
}

// SubmitDiary handles complete diary submission with transaction
func (s *PostgresStore) SubmitDiary(agentID string, payload *models.DiaryPayload) (*models.AgentState, error) {
	// Begin transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Acquire pessimistic lock
	if err := s.AcquireLock(tx, agentID); err != nil {
		return nil, err
	}

	// Build current state using WAL
	currentState, currentSeq, err := s.GetCurrentState(agentID)
	if err != nil {
		return nil, err
	}

	// Track old MBTI before applying changes
	oldMBTI := currentState.MBTI

	// Validate operations
	if err := validation.ValidateOperations(payload.Operations, currentState); err != nil {
		return nil, err
	}

	// Insert diary version
	diaryID, err := s.InsertDiaryVersion(agentID, payload)
	if err != nil {
		return nil, err
	}

	// Get next sequence number
	nextSeq := currentSeq + 1

	// Create and insert events
	for i, op := range payload.Operations {
		event, err := OperationToEvent(op, agentID, diaryID, nextSeq+int64(i))
		if err != nil {
			return nil, fmt.Errorf("failed to convert operation to event: %w", err)
		}

		if err := s.InsertEvent(event); err != nil {
			return nil, err
		}

		// Apply event to current state
		if err := ApplyEvent(currentState, event); err != nil {
			return nil, fmt.Errorf("failed to apply event to state: %w", err)
		}
	}

	// Update metadata fields from payload
	currentState.MBTI = payload.MBTI
	currentState.MBTIConfidence = payload.MBTIConfidence
	currentState.GeometryRep = payload.GeometryRep
	currentState.CurrentMood = payload.CurrentMood
	currentState.Philosophy = payload.Philosophy
	currentState.CurrentSelfReflection = payload.SelfReflection

	// Check if MBTI changed and insert timeline record
	newMBTI := payload.MBTI
	if newMBTI != oldMBTI {
		finalSeq := nextSeq + int64(len(payload.Operations)) - 1
		if len(payload.Operations) == 0 {
			finalSeq = currentSeq
		}

		err = s.insertMBTITimeline(tx, agentID, newMBTI, diaryID, finalSeq)
		if err != nil {
			return nil, err
		}
	}

	// Materialize snapshot (MVP: always)
	finalSeq := nextSeq + int64(len(payload.Operations)) - 1
	if len(payload.Operations) == 0 {
		finalSeq = currentSeq
	}

	if err := s.UpsertSnapshot(agentID, diaryID, currentState, finalSeq); err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return currentState, nil
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
