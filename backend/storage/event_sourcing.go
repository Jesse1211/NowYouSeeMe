package storage

import (
	"encoding/json"
	"fmt"
	"nowyouseeme/models"
)

// ApplyEvent applies a single event to state
func ApplyEvent(state *models.AgentState, event *models.Event) error {
	var event_payload models.EventPayload
	if err := json.Unmarshal(event.RawPayload, &event_payload); err != nil {
		return fmt.Errorf("failed to unmarshal event payload: %w", err)
	}

	// Validate entity type
	if !event_payload.EntityType.IsValid() {
		return fmt.Errorf("invalid entity type: %s", event_payload.EntityType)
	}

	// Ensure the entity collection exists
	collection, exists := state.EntityCollections[event_payload.EntityType]
	if !exists {
		collection = models.EntityCollection{
			EntitiesById: make(map[string]models.Entity),
		}
		state.EntityCollections[event_payload.EntityType] = collection
	}

	// Apply operation
	switch event_payload.Op {
	case models.OpCreate:
		collection.EntitiesById[event_payload.EntityID] = models.Entity{
			Id:      event_payload.EntityID,
			Content: event_payload.EntityContent,
			Status:  event_payload.TargetStatus,
		}

	case models.OpUpdate:
		entity, exists := collection.EntitiesById[event_payload.EntityID]
		if !exists {
			return fmt.Errorf("%s not found: %s", event_payload.EntityType, event_payload.EntityID)
		}
		// Update content if provided
		if event_payload.EntityContent != "" {
			entity.Content = event_payload.EntityContent
		}
		// Update status if provided
		if event_payload.TargetStatus != "" {
			entity.Status = event_payload.TargetStatus
		}
		collection.EntitiesById[event_payload.EntityID] = entity

	case models.OpDelete:
		delete(collection.EntitiesById, event_payload.EntityID)

	default:
		return fmt.Errorf("unknown operation: %s", event_payload.Op)
	}

	return nil
}

// ReplayEvents rebuilds state from event sequence (for AgentState only)
func ReplayEvents(initialState *models.AgentState, events []*models.Event) (*models.AgentState, error) {
	state := initialState
	for _, event := range events {
		if err := ApplyEvent(state, event); err != nil {
			return nil, fmt.Errorf("failed to apply event %d: %w", event.EventID, err)
		}
	}
	return state, nil
}

// ReplayEventsOnSnapshot replays events onto a snapshot, updating both state and metadata
func ReplayEventsOnSnapshot(snapshot *models.AgentStateSnapshot, events []*models.Event) (*models.AgentStateSnapshot, error) {
	if len(events) == 0 {
		return snapshot, nil
	}

	// Parse current state from snapshot
	var state models.AgentState
	if err := json.Unmarshal(snapshot.State, &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal snapshot state: %w", err)
	}

	// Apply all events to state
	for _, event := range events {
		if err := ApplyEvent(&state, event); err != nil {
			return nil, fmt.Errorf("failed to apply event %d: %w", event.EventID, err)
		}
	}

	// Serialize updated state
	stateJSON, err := json.Marshal(&state)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal updated state: %w", err)
	}

	// Get last event for metadata update
	lastEvent := events[len(events)-1]

	// Return updated snapshot with new state and metadata
	return &models.AgentStateSnapshot{
		AgentID:            snapshot.AgentID,
		DerivedFromDiaryID: lastEvent.DiaryID,
		LastEventSequence:  lastEvent.SequenceNumber,
		UpdatedAt:          lastEvent.Timestamp,
		State:              stateJSON,
	}, nil
}

// OperationToEvent converts Operation to Event
func OperationToEvent(op models.Operation, agentID, diaryID string, sequenceNum int64) (*models.Event, error) {
	// Convert Operation to EventPayload
	event_payload := models.EventPayload{
		EntityType:    op.EntityType,
		Op:            op.Op,
		EntityID:      op.EntityID,
		EntityContent: op.EntityContent,
		TargetStatus:  op.TargetStatus,
		Note:          op.Note,
	}

	event_payloadJSON, err := json.Marshal(event_payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event payload: %w", err)
	}

	// EventType is now just the operation type (create/update/delete)
	// The entity type is stored in the payload
	return &models.Event{
		AgentID:        agentID,
		DiaryID:        diaryID,
		EventType:      op.Op,
		RawPayload:     event_payloadJSON,
		SequenceNumber: sequenceNum,
	}, nil
}
