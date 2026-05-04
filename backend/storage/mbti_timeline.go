package storage

import (
	"database/sql"
	"fmt"
)

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
		WHERE current_mbti like $1
	`

	rows, err := s.db.Query(query, mbtiType+"%") // Use prefix match for flexibility (e.g. "INTJ" matches "INTJ-A" and "INTJ-T")
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
