package models

import (
	"encoding/json"
	"fmt"
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

// AgentSnapshotResult encapsulates the complete snapshot information
// following DDD principles - a rich domain object
type AgentSnapshotResult struct {
	AgentID   string      `json:"agent_id"`
	State     *AgentState `json:"state"`
	Sequence  int64       `json:"sequence"`
	UpdatedAt *time.Time  `json:"updated_at,omitempty"`
}

// HasSnapshot returns true if this result contains a valid snapshot
func (r *AgentSnapshotResult) HasSnapshot() bool {
	return r.UpdatedAt != nil
}

// FormattedUpdatedAt returns the UpdatedAt timestamp in ISO 8601 format
// Returns empty string if no snapshot exists
func (r *AgentSnapshotResult) FormattedUpdatedAt() (string, error) {
	if r.UpdatedAt == nil {
		return "", fmt.Errorf("no snapshot available")
	}
	return r.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"), nil
}
