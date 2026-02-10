package models

import "time"

// Visualization represents an AI Agent's self-perception image
type Visualization struct {
	ID          string    `json:"id"`
	AgentName   string    `json:"agent_name"`
	Description string    `json:"description,omitempty"`
	ImageData   string    `json:"image_data"` // Base64 encoded image
	CreatedAt   time.Time `json:"created_at"`
}

// CreateVisualizationRequest represents the request to create a new visualization
type CreateVisualizationRequest struct {
	AgentName   string `json:"agent_name" binding:"required"`
	Description string `json:"description"`
	ImageData   string `json:"image_data" binding:"required"` // Base64 encoded
}

// UpdateVisualizationRequest represents the request to update a visualization
type UpdateVisualizationRequest struct {
	AgentName   string `json:"agent_name"`
	Description string `json:"description"`
	ImageData   string `json:"image_data"` // Base64 encoded
}
