package models

import (
	"encoding/json"
	"time"
)

// AgentStateSnapshot represents materialized current state
type AgentStateSnapshot struct {
	AgentID            string          `json:"agent_id"`
	DerivedFromDiaryID string          `json:"derived_from_diary_id"`
	LastEventSequence  int64           `json:"last_event_sequence"`
	UpdatedAt          time.Time       `json:"updated_at"`
	State              json.RawMessage `json:"state"`
}

// AgentState represents the JSONB state structure
type AgentState struct {
	MBTI                  string         `json:"mbti"`
	MBTIConfidence        float64        `json:"mbti_confidence"`
	GeometryRep           string         `json:"geometry_representation"`
	CurrentMood           string         `json:"current_mood"`
	Philosophy            string         `json:"philosophy"`
	CurrentSelfReflection SelfReflection `json:"current_self_reflection"`

	Goals        map[string]Goal   `json:"goals"`
	Capabilities map[string]Entity `json:"capabilities"`
	Limitations  map[string]Entity `json:"limitations"`
	Aspirations  map[string]Entity `json:"aspirations"`
}

// Goal represents a goal entity
type Goal struct {
	Title      string  `json:"title"`
	Status     string  `json:"status"`
	Checkpoint *string `json:"checkpoint,omitempty"`
}

// Entity represents capability/limitation/aspiration
type Entity struct {
	Title string `json:"title"`
}

// NewEmptyState creates an empty AgentState
func NewEmptyState() *AgentState {
	return &AgentState{
		Goals:        make(map[string]Goal),
		Capabilities: make(map[string]Entity),
		Limitations:  make(map[string]Entity),
		Aspirations:  make(map[string]Entity),
	}
}
