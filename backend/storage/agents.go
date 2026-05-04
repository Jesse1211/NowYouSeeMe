package storage

import (
	"context"
	"database/sql"
	"fmt"
	"nowyouseeme/models"
	"time"
)

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
