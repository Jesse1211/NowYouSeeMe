package storage

import (
	"database/sql"
	"fmt"
	"nowyouseeme/models"
	"time"
)

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
	return s.getUncommittedEventsTx(nil, agentID, afterSequence)
}

// getUncommittedEventsTx retrieves events within a transaction or using default connection
func (s *PostgresStore) getUncommittedEventsTx(tx *sql.Tx, agentID string, afterSequence int64) ([]*models.Event, error) {
	query := `
		SELECT event_id, agent_id, diary_id, event_type, timestamp, payload, sequence_number
		FROM events
		WHERE agent_id = $1 AND sequence_number > $2
		ORDER BY sequence_number ASC
	`

	var rows *sql.Rows
	var err error

	if tx != nil {
		rows, err = tx.Query(query, agentID, afterSequence)
	} else {
		rows, err = s.db.Query(query, agentID, afterSequence)
	}

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
