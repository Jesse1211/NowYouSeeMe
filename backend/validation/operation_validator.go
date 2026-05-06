package validation

import (
	"fmt"
	"nowyouseeme/models"
)

// ValidationError represents validation failure
type ValidationError struct {
	Errors []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %v", e.Errors)
}

// ValidateOperations validates operations against current state
func ValidateOperations(operations []models.Operation, latestState *models.AgentState) error {
	errors := []string{}
	tempState := cloneState(latestState)

	for i, op := range operations {
		if err := validateOperation(op, tempState); err != nil {
			errors = append(errors, fmt.Sprintf("operation[%d]: %s", i, err.Error()))
			continue
		}

		// Apply operation to temp state for subsequent validations
		applyOperationToState(op, tempState)
	}

	if len(errors) > 0 {
		return &ValidationError{Errors: errors}
	}

	return nil
}

func validateOperation(op models.Operation, state *models.AgentState) error {
	switch op.Op {
	case models.OpGoalCreate:
		return validateGoalCreate(op, state)
	case models.OpGoalTransition:
		return validateGoalTransition(op, state)
	case models.OpGoalUpdate:
		return validateGoalUpdate(op, state)
	case models.OpGoalComplete:
		return validateGoalComplete(op, state)
	case models.OpGoalAbandon:
		return validateGoalAbandon(op, state)

	case models.OpCapabilityAdd:
		return validateEntityAdd(op.CapabilityID, state.Capabilities, "capability")
	case models.OpCapabilityRemove:
		return validateEntityRemove(op.CapabilityID, state.Capabilities, "capability")
	case models.OpCapabilityUpdate:
		return validateEntityUpdate(op.CapabilityID, state.Capabilities, "capability")

	case models.OpLimitationAdd:
		return validateEntityAdd(op.LimitationID, state.Limitations, "limitation")
	case models.OpLimitationRemove:
		return validateEntityRemove(op.LimitationID, state.Limitations, "limitation")
	case models.OpLimitationUpdate:
		return validateEntityUpdate(op.LimitationID, state.Limitations, "limitation")

	case models.OpAspirationAdd:
		return validateEntityAdd(op.AspirationID, state.Aspirations, "aspiration")
	case models.OpAspirationRemove:
		return validateEntityRemove(op.AspirationID, state.Aspirations, "aspiration")
	case models.OpAspirationUpdate:
		return validateEntityUpdate(op.AspirationID, state.Aspirations, "aspiration")

	default:
		return fmt.Errorf("unknown operation type: %s", op.Op)
	}
}

// Goal validation functions
func validateGoalCreate(op models.Operation, state *models.AgentState) error {
	if op.GoalID == "" || op.Title == "" || op.Status == "" {
		return fmt.Errorf("goal_create requires goal_id, title, and status")
	}
	if _, exists := state.Goals[op.GoalID]; exists {
		return fmt.Errorf("goal %s already exists", op.GoalID)
	}
	return ValidateGoalCreateStatus(op.Status)
}

func validateGoalTransition(op models.Operation, state *models.AgentState) error {
	if op.GoalID == "" || op.FromStatus == "" || op.ToStatus == "" {
		return fmt.Errorf("goal_transition requires goal_id, from_status, and to_status")
	}

	goal, exists := state.Goals[op.GoalID]
	if !exists {
		return fmt.Errorf("goal %s not found", op.GoalID)
	}

	if goal.Status != op.FromStatus {
		return fmt.Errorf("goal %s is in status %s, not %s", op.GoalID, goal.Status, op.FromStatus)
	}

	return ValidateGoalTransition(op.FromStatus, op.ToStatus)
}

func validateGoalUpdate(op models.Operation, state *models.AgentState) error {
	if op.GoalID == "" || op.Title == "" {
		return fmt.Errorf("goal_update requires goal_id and title")
	}
	if _, exists := state.Goals[op.GoalID]; !exists {
		return fmt.Errorf("goal %s not found", op.GoalID)
	}
	return nil
}

func validateGoalComplete(op models.Operation, state *models.AgentState) error {
	if op.GoalID == "" {
		return fmt.Errorf("goal_complete requires goal_id")
	}

	goal, exists := state.Goals[op.GoalID]
	if !exists {
		return fmt.Errorf("goal %s not found", op.GoalID)
	}

	if goal.Status != "progressing" {
		return fmt.Errorf("can only complete goals in 'progressing' status, got: %s", goal.Status)
	}

	return nil
}

func validateGoalAbandon(op models.Operation, state *models.AgentState) error {
	if op.GoalID == "" {
		return fmt.Errorf("goal_abandon requires goal_id")
	}
	if _, exists := state.Goals[op.GoalID]; !exists {
		return fmt.Errorf("goal %s not found", op.GoalID)
	}
	return nil
}

// Entity validation functions
func validateEntityAdd(id string, entities map[string]models.Entity, entityType string) error {
	if id == "" {
		return fmt.Errorf("%s_add requires %s_id", entityType, entityType)
	}
	if _, exists := entities[id]; exists {
		return fmt.Errorf("%s %s already exists", entityType, id)
	}
	return nil
}

func validateEntityRemove(id string, entities map[string]models.Entity, entityType string) error {
	if id == "" {
		return fmt.Errorf("%s_remove requires %s_id", entityType, entityType)
	}
	if _, exists := entities[id]; !exists {
		return fmt.Errorf("%s %s not found", entityType, id)
	}
	return nil
}

func validateEntityUpdate(id string, entities map[string]models.Entity, entityType string) error {
	if id == "" {
		return fmt.Errorf("%s_update requires %s_id", entityType, entityType)
	}
	if _, exists := entities[id]; !exists {
		return fmt.Errorf("%s %s not found", entityType, id)
	}
	return nil
}

// State manipulation helpers
func cloneState(state *models.AgentState) *models.AgentState {
	clone := &models.AgentState{
		MBTI:                 state.MBTI,
		MBTIConfidence:       state.MBTIConfidence,
		GeometryRep:          state.GeometryRep,
		CurrentMood:          state.CurrentMood,
		Philosophy:           state.Philosophy,
		CurrentSelfReflection: state.CurrentSelfReflection,
		Goals:                make(map[string]models.Goal),
		Capabilities:         make(map[string]models.Entity),
		Limitations:          make(map[string]models.Entity),
		Aspirations:          make(map[string]models.Entity),
	}

	for k, v := range state.Goals {
		clone.Goals[k] = v
	}
	for k, v := range state.Capabilities {
		clone.Capabilities[k] = v
	}
	for k, v := range state.Limitations {
		clone.Limitations[k] = v
	}
	for k, v := range state.Aspirations {
		clone.Aspirations[k] = v
	}

	return clone
}

func applyOperationToState(op models.Operation, state *models.AgentState) {
	switch op.Op {
	case models.OpGoalCreate:
		state.Goals[op.GoalID] = models.Goal{Title: op.Title, Status: op.Status}
	case models.OpGoalTransition:
		goal := state.Goals[op.GoalID]
		goal.Status = op.ToStatus
		state.Goals[op.GoalID] = goal
	case models.OpGoalUpdate:
		goal := state.Goals[op.GoalID]
		goal.Title = op.Title
		state.Goals[op.GoalID] = goal
	case models.OpGoalComplete:
		goal := state.Goals[op.GoalID]
		goal.Status = "completed"
		state.Goals[op.GoalID] = goal
	case models.OpGoalAbandon:
		goal := state.Goals[op.GoalID]
		goal.Status = "abandoned"
		state.Goals[op.GoalID] = goal

	case models.OpCapabilityAdd:
		state.Capabilities[op.CapabilityID] = models.Entity{Title: op.Title}
	case models.OpCapabilityRemove:
		delete(state.Capabilities, op.CapabilityID)
	case models.OpCapabilityUpdate:
		state.Capabilities[op.CapabilityID] = models.Entity{Title: op.Title}

	case models.OpLimitationAdd:
		state.Limitations[op.LimitationID] = models.Entity{Title: op.Title}
	case models.OpLimitationRemove:
		delete(state.Limitations, op.LimitationID)
	case models.OpLimitationUpdate:
		state.Limitations[op.LimitationID] = models.Entity{Title: op.Title}

	case models.OpAspirationAdd:
		state.Aspirations[op.AspirationID] = models.Entity{Title: op.Title}
	case models.OpAspirationRemove:
		delete(state.Aspirations, op.AspirationID)
	case models.OpAspirationUpdate:
		state.Aspirations[op.AspirationID] = models.Entity{Title: op.Title}
	}
}
