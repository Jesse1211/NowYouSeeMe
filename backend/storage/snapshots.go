package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"nowyouseeme/models"
	"time"
)

func (s *PostgresStore) GetSnapshot(agentID string) (*models.AgentSnapshotResult, error) {
	return s.getSnapshotTx(nil, agentID)
}

// getSnapshotTx retrieves snapshot within a transaction or using default connection
func (s *PostgresStore) getSnapshotTx(tx *sql.Tx, agentID string) (*models.AgentSnapshotResult, error) {
	query := `
		SELECT agent_id, last_event_sequence, updated_at, state
		FROM agent_state_snapshots
		WHERE agent_id = $1
	`

	snapshot := &models.AgentSnapshot{}
	var err error

	if tx != nil {
		err = tx.QueryRow(query, agentID).Scan(
			&snapshot.AgentID,
			&snapshot.LastEventSequence,
			&snapshot.UpdatedAt,
			&snapshot.State,
		)
	} else {
		err = s.db.QueryRow(query, agentID).Scan(
			&snapshot.AgentID,
			&snapshot.LastEventSequence,
			&snapshot.UpdatedAt,
			&snapshot.State,
		)
	}

	if err == sql.ErrNoRows {
		return nil, nil // No snapshot exists yet
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshot: %w", err)
	}

	return snapshot.FromAgentSnapshotToAgentSnapshotResult()
}

// UpsertSnapshot creates or updates a snapshot
func (s *PostgresStore) UpsertSnapshot(agentID string, state *models.AgentState, lastEventSeq int64) error {
	return s.upsertSnapshotTx(nil, agentID, state, lastEventSeq)
}

// upsertSnapshotTx creates or updates a snapshot within a transaction or using default connection
func (s *PostgresStore) upsertSnapshotTx(tx *sql.Tx, agentID string, state *models.AgentState, lastEventSeq int64) error {
	stateJSON, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	query := `
		INSERT INTO agent_state_snapshots (agent_id, last_event_sequence, updated_at, state)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (agent_id)
		DO UPDATE SET
			last_event_sequence = EXCLUDED.last_event_sequence,
			updated_at = EXCLUDED.updated_at,
			state = EXCLUDED.state
	`

	var execErr error
	if tx != nil {
		_, execErr = tx.Exec(query, agentID, lastEventSeq, time.Now(), stateJSON)
	} else {
		_, execErr = s.db.Exec(query, agentID, lastEventSeq, time.Now(), stateJSON)
	}

	if execErr != nil {
		return fmt.Errorf("failed to upsert snapshot: %w", execErr)
	}

	return nil
}

// GetLatestSnapshot builds latest snapshot using WAL pattern with transactional consistency
// Returns a rich domain object encapsulating all snapshot information
func (s *PostgresStore) GetLatestSnapshot(agentID string) (*models.AgentSnapshotResult, error) {
	// Use a read transaction to ensure consistent view of snapshot + uncommitted events
	// This prevents race conditions where snapshot might be updated between the two reads
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{
		ReadOnly: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin read transaction: %w", err)
	}
	defer tx.Rollback() // Safe to call even after commit

	// Load latest snapshot within transaction
	snapshot, err := s.getSnapshotTx(tx, agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshot: %w", err)
	}

	// No snapshot exists yet - just created account without any diary submission
	if snapshot == nil {
		tx.Commit() // Explicit commit for read-only transaction
		return nil, nil
	}

	// Load uncommitted events (WAL) within same transaction for consistency
	// Note: In current implementation, snapshots are materialized after every diary submission,
	// so uncommittedEvents should typically be empty. WAL replay only occurs during:
	// 1. Concurrent reads while another transaction is committing
	// 2. Recovery scenarios if snapshot update failed but events were persisted
	uncommittedEvents, err := s.getUncommittedEventsTx(tx, agentID, snapshot.LastEventSequence)
	if err != nil {
		return nil, fmt.Errorf("failed to get uncommitted events: %w", err)
	}

	// Commit read transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit read transaction: %w", err)
	}

	// Replay uncommitted events onto snapshot (updates both state and metadata)
	latestSnapshot := snapshot
	if len(uncommittedEvents) > 0 {
		latestSnapshot, err = ReplayEventsOnSnapshot(snapshot, uncommittedEvents)
		if err != nil {
			return nil, fmt.Errorf("failed to replay events: %w", err)
		}
	}

	return latestSnapshot, nil
}

// AcquireLock acquires pessimistic lock for agent
// Always locks the agents table for consistency and to avoid deadlocks
