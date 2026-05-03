package storage

import (
	"encoding/json"
	"fmt"
	"nowyouseeme/models"
)

// ApplyEvent applies a single event to state
func ApplyEvent(state *models.AgentState, event *models.Event) error {
	var payload models.EventPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return fmt.Errorf("failed to unmarshal event payload: %w", err)
	}

	switch event.EventType {
	case "goal_create":
		state.Goals[payload.GoalID] = models.Goal{
			Title:  payload.Title,
			Status: payload.Status,
		}

	case "goal_transition":
		goal := state.Goals[payload.GoalID]
		goal.Status = payload.ToStatus
		state.Goals[payload.GoalID] = goal

	case "goal_update":
		goal := state.Goals[payload.GoalID]
		goal.Title = payload.Title
		state.Goals[payload.GoalID] = goal

	case "goal_complete":
		goal := state.Goals[payload.GoalID]
		goal.Status = "completed"
		if payload.Checkpoint != "" {
			goal.Checkpoint = &payload.Checkpoint
		}
		state.Goals[payload.GoalID] = goal

	case "goal_abandon":
		goal := state.Goals[payload.GoalID]
		goal.Status = "abandoned"
		state.Goals[payload.GoalID] = goal

	case "capability_add":
		state.Capabilities[payload.CapabilityID] = models.Entity{Title: payload.Title}
	case "capability_remove":
		delete(state.Capabilities, payload.CapabilityID)
	case "capability_update":
		state.Capabilities[payload.CapabilityID] = models.Entity{Title: payload.Title}

	case "limitation_add":
		state.Limitations[payload.LimitationID] = models.Entity{Title: payload.Title}
	case "limitation_remove":
		delete(state.Limitations, payload.LimitationID)
	case "limitation_update":
		state.Limitations[payload.LimitationID] = models.Entity{Title: payload.Title}

	case "aspiration_add":
		state.Aspirations[payload.AspirationID] = models.Entity{Title: payload.Title}
	case "aspiration_remove":
		delete(state.Aspirations, payload.AspirationID)
	case "aspiration_update":
		state.Aspirations[payload.AspirationID] = models.Entity{Title: payload.Title}

	default:
		return fmt.Errorf("unknown event type: %s", event.EventType)
	}

	return nil
}

// ReplayEvents rebuilds state from event sequence
func ReplayEvents(initialState *models.AgentState, events []*models.Event) (*models.AgentState, error) {
	state := initialState
	for _, event := range events {
		if err := ApplyEvent(state, event); err != nil {
			return nil, fmt.Errorf("failed to apply event %d: %w", event.EventID, err)
		}
	}
	return state, nil
}

// OperationToEvent converts Operation to Event
func OperationToEvent(op models.Operation, agentID, diaryID string, sequenceNum int64) (*models.Event, error) {
	payload := models.EventPayload{
		GoalID:       op.GoalID,
		Title:        op.Title,
		Status:       op.Status,
		FromStatus:   op.FromStatus,
		ToStatus:     op.ToStatus,
		Reason:       op.Reason,
		CapabilityID: op.CapabilityID,
		LimitationID: op.LimitationID,
		AspirationID: op.AspirationID,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event payload: %w", err)
	}

	return &models.Event{
		AgentID:        agentID,
		DiaryID:        diaryID,
		EventType:      op.Op,
		Payload:        payloadJSON,
		SequenceNumber: sequenceNum,
	}, nil
}
