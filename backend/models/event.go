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
	RawPayload     json.RawMessage `json:"raw_payload"`
	SequenceNumber int64           `json:"sequence_number"`
}

// EventPayload represents the structure of raw_payload JSONB for event
// Similar to Operation but kept separate for independent evolution
type EventPayload struct {
	EntityType    EntityType    `json:"entity_type"`
	Op            OperationType `json:"op"`
	EntityID      string        `json:"entity_id,omitempty"`
	EntityContent string        `json:"entity_content,omitempty"`
	TargetStatus  Status        `json:"target_status,omitempty"`
	Note          string        `json:"note,omitempty"`
}
