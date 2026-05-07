package storage

import (
	"encoding/json"
	"fmt"
	"nowyouseeme/models"
)

// ApplyEventToSnapshot applies a single event to state
// Metadata events (EventMetadata) are skipped - they don't participate in AgentState replay
// Only operation events (EventCreate/EventUpdate/EventDelete) modify the state
func ApplyEventToSnapshot(snapshot *models.AgentSnapshotResult, event *models.Event) error {
	// Skip metadata events - they don't participate in AgentState replay
	if event.EventType == models.EventMetadata {
		return nil
	}

	state := snapshot.State

	// Parse operation payload
	var operationPayload models.OperationPayload
	if err := json.Unmarshal(event.RawPayload, &operationPayload); err != nil {
		return fmt.Errorf("failed to unmarshal operation payload: %w", err)
	}

	// Validate entity type
	if !operationPayload.EntityType.IsValid() {
		return fmt.Errorf("invalid entity type: %s", operationPayload.EntityType)
	}

	// Ensure the entity collection exists
	collection, exists := state.EntityCollections[operationPayload.EntityType]
	if !exists {
		collection = models.EntityCollection{
			EntitiesById: make(map[string]models.Entity),
		}
		state.EntityCollections[operationPayload.EntityType] = collection
	}

	// Apply operation
	switch event.EventType {
	case models.EventCreate:
		collection.EntitiesById[operationPayload.EntityID] = models.Entity{
			Id:      operationPayload.EntityID,
			Content: operationPayload.EntityContent,
			Status:  operationPayload.TargetStatus,
		}

	case models.EventUpdate:
		entity, exists := collection.EntitiesById[operationPayload.EntityID]
		if !exists {
			return fmt.Errorf("%s not found: %s", operationPayload.EntityType, operationPayload.EntityID)
		}
		// Update content if provided
		if operationPayload.EntityContent != "" {
			entity.Content = operationPayload.EntityContent
		}
		// Update status if provided
		if operationPayload.TargetStatus != "" {
			entity.Status = operationPayload.TargetStatus
		}
		collection.EntitiesById[operationPayload.EntityID] = entity

	case models.EventDelete:
		delete(collection.EntitiesById, operationPayload.EntityID)

	default:
		return fmt.Errorf("unknown operation: %s", event.EventType)
	}

	// Update snapshot metadata
	snapshot.LastEventSequence = event.SequenceNumber
	snapshot.UpdatedAt = &event.Timestamp

	return nil
}

// ReplayEventsOnSnapshot replays events onto a snapshot, updating both state and metadata
func ReplayEventsOnSnapshot(snapshot *models.AgentSnapshotResult, events []*models.Event) (*models.AgentSnapshotResult, error) {
	if len(events) == 0 {
		return snapshot, nil
	}

	// Apply all events to state
	for _, event := range events {
		if err := ApplyEventToSnapshot(snapshot, event); err != nil {
			return nil, fmt.Errorf("failed to apply event %d: %w", event.EventID, err)
		}
	}

	return snapshot, nil
}

// OperationToEventConverter converts Operation to Event with OperationPayload
func OperationToEventConverter(op models.Operation, agentID, diaryID string, sequenceNum int64) (*models.Event, error) {
	// Convert Operation to OperationPayload
	operationPayload := models.OperationPayload{
		EntityType:    op.EntityType,
		EntityID:      op.EntityID,
		EntityContent: op.EntityContent,
		TargetStatus:  op.TargetStatus,
		Note:          op.Note,
	}

	payloadJSON, err := json.Marshal(operationPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal operation payload: %w", err)
	}

	// EventType is the operation type (create/update/delete)
	// The entity type is stored in the payload
	return &models.Event{
		AgentID:        agentID,
		DiaryID:        diaryID,
		EventType:      models.OperationTypeToEventType(op.Op),
		RawPayload:     payloadJSON,
		SequenceNumber: sequenceNum,
	}, nil
}

// MetadataToEventConverter creates a metadata event from DiaryPayload
func MetadataToEventConverter(payload *models.DiaryPayload, agentID, diaryID string, sequenceNum int64) (*models.Event, error) {
	// Extract metadata fields from DiaryPayload
	metadataPayload := models.MetadataPayload{
		MBTI:           payload.MBTI,
		MBTIConfidence: payload.MBTIConfidence,
		GeometryRep:    payload.GeometryRep,
		Context:        payload.Context,
		CurrentMood:    payload.CurrentMood,
		Philosophy:     payload.Philosophy,
		SelfReflection: map[string]any{
			"rumination_for_yesterday":  payload.SelfReflection.Rumination,
			"what_happened_today":       payload.SelfReflection.WhatHappened,
			"expectations_for_tomorrow": payload.SelfReflection.Expectations,
		},
	}

	payloadJSON, err := json.Marshal(metadataPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata payload: %w", err)
	}

	return &models.Event{
		AgentID:        agentID,
		DiaryID:        diaryID,
		EventType:      models.EventMetadata,
		RawPayload:     payloadJSON,
		SequenceNumber: sequenceNum,
	}, nil
}
