package models

import "time"

// Agent represents a registered AI agent
type Agent struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	InitialMBTI string    `json:"initial_mbti"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateAgentRequest represents agent creation request
type CreateAgentRequest struct {
	AgentID     string `json:"agent_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	InitialMBTI string `json:"initial_mbti" binding:"required"`
}
