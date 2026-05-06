package models

import (
	"encoding/json"
	"time"
)

// Event represents a single state-changing event
type Event struct {
	EventID        int64           `json:"event_id"`
	AgentID        string          `json:"agent_id"`
	DiaryID        string          `json:"diary_id"`
	EventType      OperationType   `json:"event_type"`
	Timestamp      time.Time       `json:"timestamp"`
	Payload        json.RawMessage `json:"payload"`
	SequenceNumber int64           `json:"sequence_number"`
}

// EventPayload represents typed event payloads
type EventPayload struct {
	// Goal events
	GoalID     string `json:"goal_id,omitempty"`
	Title      string `json:"title,omitempty"`
	Status     string `json:"status,omitempty"`
	FromStatus string `json:"from_status,omitempty"`
	ToStatus   string `json:"to_status,omitempty"`
	Checkpoint string `json:"checkpoint,omitempty"`
	Reason     string `json:"reason,omitempty"`

	// Entity events
	CapabilityID string `json:"capability_id,omitempty"`
	LimitationID string `json:"limitation_id,omitempty"`
	AspirationID string `json:"aspiration_id,omitempty"`
}
