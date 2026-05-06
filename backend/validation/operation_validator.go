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
	// Validate entity type
	if !op.EntityType.IsValid() {
		return fmt.Errorf("invalid entity type: %s", op.EntityType)
	}

	// Validate operation type
	if !op.Op.IsValid() {
		return fmt.Errorf("invalid operation type: %s", op.Op)
	}

	// Get entity collection (create if needed)
	collection, exists := state.EntityCollections[op.EntityType]
	if !exists {
		collection = models.EntityCollection{EntitiesById: make(map[string]models.Entity)}
		state.EntityCollections[op.EntityType] = collection
	}

	// Validate based on operation type
	switch op.Op {
	case models.OpCreate:
		return validateCreate(op, collection, op.EntityType)
	case models.OpUpdate:
		return validateUpdate(op, collection, op.EntityType)
	case models.OpDelete:
		return validateDelete(op, collection, op.EntityType)
	default:
		return fmt.Errorf("unsupported operation: %s", op.Op)
	}
}

func validateCreate(op models.Operation, collection models.EntityCollection, entityType models.EntityType) error {
	// Validate required fields
	if op.EntityID == "" || op.EntityContent == "" {
		return fmt.Errorf("create %s requires entity_id and entity_content", entityType)
	}

	// Check if entity already exists
	if _, exists := collection.EntitiesById[op.EntityID]; exists {
		return fmt.Errorf("%s %s already exists", entityType, op.EntityID)
	}

	// Validate target status
	if op.TargetStatus == "" {
		return fmt.Errorf("create %s requires target_status", entityType)
	}
	if !op.TargetStatus.IsValid() {
		return fmt.Errorf("invalid status: %s", op.TargetStatus)
	}

	// Only allow pending or progress for new entities
	if op.TargetStatus != models.StatusPending && op.TargetStatus != models.StatusProgress {
		return fmt.Errorf("new %s can only be pending or progress, got: %s", entityType, op.TargetStatus)
	}

	return nil
}

func validateUpdate(op models.Operation, collection models.EntityCollection, entityType models.EntityType) error {
	// Validate required fields
	if op.EntityID == "" {
		return fmt.Errorf("update %s requires entity_id", entityType)
	}

	// Check if entity exists
	entity, exists := collection.EntitiesById[op.EntityID]
	if !exists {
		return fmt.Errorf("%s %s not found", entityType, op.EntityID)
	}

	// At least one field must be provided for update
	if op.EntityContent == "" && op.TargetStatus == "" {
		return fmt.Errorf("update %s requires at least one of: entity_content or target_status", entityType)
	}

	// Validate status transition if provided
	if op.TargetStatus != "" {
		if !op.TargetStatus.IsValid() {
			return fmt.Errorf("invalid status: %s", op.TargetStatus)
		}
		// For goals, validate status transitions
		if entityType == models.EntityGoal {
			if err := ValidateGoalStatusTransition(entity.Status, op.TargetStatus); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateDelete(op models.Operation, collection models.EntityCollection, entityType models.EntityType) error {
	// Validate required fields
	if op.EntityID == "" {
		return fmt.Errorf("delete %s requires entity_id", entityType)
	}

	// Check if entity exists
	if _, exists := collection.EntitiesById[op.EntityID]; !exists {
		return fmt.Errorf("%s %s not found", entityType, op.EntityID)
	}

	return nil
}

// State manipulation helpers
func cloneState(state *models.AgentState) *models.AgentState {
	clone := &models.AgentState{
		MBTI:                  state.MBTI,
		MBTIConfidence:        state.MBTIConfidence,
		GeometryRep:           state.GeometryRep,
		CurrentMood:           state.CurrentMood,
		Philosophy:            state.Philosophy,
		CurrentSelfReflection: state.CurrentSelfReflection,
		EntityCollections:     make(map[models.EntityType]models.EntityCollection),
	}

	// Clone each entity collection
	for entityType, collection := range state.EntityCollections {
		clonedCollection := models.EntityCollection{
			EntitiesById: make(map[string]models.Entity),
		}
		for id, entity := range collection.EntitiesById {
			clonedCollection.EntitiesById[id] = entity
		}
		clone.EntityCollections[entityType] = clonedCollection
	}

	return clone
}

func applyOperationToState(op models.Operation, state *models.AgentState) {
	// Ensure collection exists
	collection, exists := state.EntityCollections[op.EntityType]
	if !exists {
		collection = models.EntityCollection{EntitiesById: make(map[string]models.Entity)}
		state.EntityCollections[op.EntityType] = collection
	}

	switch op.Op {
	case models.OpCreate:
		collection.EntitiesById[op.EntityID] = models.Entity{
			Id:      op.EntityID,
			Content: op.EntityContent,
			Status:  op.TargetStatus,
		}

	case models.OpUpdate:
		entity := collection.EntitiesById[op.EntityID]
		if op.EntityContent != "" {
			entity.Content = op.EntityContent
		}
		if op.TargetStatus != "" {
			entity.Status = op.TargetStatus
		}
		collection.EntitiesById[op.EntityID] = entity

	case models.OpDelete:
		delete(collection.EntitiesById, op.EntityID)
	}
}
