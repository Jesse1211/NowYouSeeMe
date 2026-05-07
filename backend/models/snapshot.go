package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// AgentSnapshot represents materialized current state
type AgentSnapshot struct {
	AgentID           string          `json:"agent_id"`
	LastEventSequence int64           `json:"last_event_sequence"`
	UpdatedAt         time.Time       `json:"updated_at"`
	State             json.RawMessage `json:"state"`
}

func (r *AgentSnapshot) FromAgentSnapshotToAgentSnapshotResult() (*AgentSnapshotResult, error) {
	var state AgentState
	if err := json.Unmarshal(r.State, &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal snapshot state: %w", err)
	}

	return &AgentSnapshotResult{
		AgentID:           r.AgentID,
		LastEventSequence: r.LastEventSequence,
		UpdatedAt:         &r.UpdatedAt,
		State:             &state,
	}, nil
}

// AgentState represents the JSONB state structure
type AgentState struct {
	MBTI                  string                          `json:"mbti"`
	MBTIConfidence        float64                         `json:"mbti_confidence"`
	GeometryRep           string                          `json:"geometry_representation"`
	CurrentMood           string                          `json:"current_mood"`
	Philosophy            string                          `json:"philosophy"`
	CurrentSelfReflection SelfReflection                  `json:"current_self_reflection"`
	EntityCollections     map[EntityType]EntityCollection `json:"entity_collections"`
}

// EntityCollection represents a collection of entities of a specific type
type EntityCollection struct {
	EntitiesById map[string]Entity `json:"entities_by_id"`
}

// Entity represents any entity (goal, capability, limitation, aspiration)
type Entity struct {
	Id      string `json:"id"`
	Content string `json:"content"`
	Status  Status `json:"status"`
}

// NewEmptyState creates an empty AgentState
func NewEmptyState() *AgentState {
	return &AgentState{
		EntityCollections: map[EntityType]EntityCollection{
			EntityGoal:       {EntitiesById: make(map[string]Entity)},
			EntityCapability: {EntitiesById: make(map[string]Entity)},
			EntityLimitation: {EntitiesById: make(map[string]Entity)},
			EntityAspiration: {EntitiesById: make(map[string]Entity)},
		},
	}
}

// AgentSnapshotResult encapsulates the complete snapshot information
// following DDD principles - a rich domain object
type AgentSnapshotResult struct {
	AgentID           string      `json:"agent_id"`
	LastEventSequence int64       `json:"sequence"`
	UpdatedAt         *time.Time  `json:"updated_at,omitempty"`
	State             *AgentState `json:"state"`
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
