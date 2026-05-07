package models

import (
	"encoding/json"
	"time"
)

// Event represents a single event (both metadata and operations)
// Event.EventType determines the structure of RawPayload:
//   - EventMetadata: RawPayload is MetadataPayload (doesn't participate in AgentState replay)
//   - EventCreate/EventUpdate/EventDelete: RawPayload is OperationPayload (participates in AgentState replay)
type Event struct {
	EventID        int64           `json:"event_id"`
	AgentID        string          `json:"agent_id"`
	DiaryID        string          `json:"diary_id"`
	EventType      EventType       `json:"event_type"`
	Timestamp      time.Time       `json:"timestamp"`
	RawPayload     json.RawMessage `json:"raw_payload"`
	SequenceNumber int64           `json:"sequence_number"`
}

// OperationPayload represents the payload for operation events (create/update/delete)
// These events participate in AgentState replay
type OperationPayload struct {
	EntityType    EntityType `json:"entity_type"`
	EntityID      string     `json:"entity_id,omitempty"`
	EntityContent string     `json:"entity_content,omitempty"`
	TargetStatus  Status     `json:"target_status,omitempty"`
	Note          string     `json:"note,omitempty"`
}

// MetadataPayload represents the payload for metadata events
// These events do NOT participate in AgentState replay
type MetadataPayload struct {
	MBTI           string         `json:"mbti,omitempty"`
	MBTIConfidence float64        `json:"mbti_confidence,omitempty"`
	GeometryRep    string         `json:"geometry_representation,omitempty"`
	Context        string         `json:"context,omitempty"`
	CurrentMood    string         `json:"current_mood,omitempty"`
	Philosophy     string         `json:"philosophy,omitempty"`
	SelfReflection map[string]any `json:"self_reflection,omitempty"`
}
