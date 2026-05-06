package models

// OperationType represents the type of operation that can be performed
type OperationType string

// Operation type constants - these are the ONLY valid operation types
const (
	// Goal operations
	OpGoalCreate     OperationType = "goal_create"
	OpGoalTransition OperationType = "goal_transition"
	OpGoalUpdate     OperationType = "goal_update"
	OpGoalComplete   OperationType = "goal_complete"
	OpGoalAbandon    OperationType = "goal_abandon"

	// Capability operations
	OpCapabilityAdd    OperationType = "capability_add"
	OpCapabilityRemove OperationType = "capability_remove"
	OpCapabilityUpdate OperationType = "capability_update"

	// Limitation operations
	OpLimitationAdd    OperationType = "limitation_add"
	OpLimitationRemove OperationType = "limitation_remove"
	OpLimitationUpdate OperationType = "limitation_update"

	// Aspiration operations
	OpAspirationAdd    OperationType = "aspiration_add"
	OpAspirationRemove OperationType = "aspiration_remove"
	OpAspirationUpdate OperationType = "aspiration_update"

	// Metadata operations
	OpMetadataUpdate OperationType = "metadata_update"
)

// String returns the string representation of the operation type
func (o OperationType) String() string {
	return string(o)
}

// IsValid checks if the operation type is valid
func (o OperationType) IsValid() bool {
	switch o {
	case OpGoalCreate, OpGoalTransition, OpGoalUpdate, OpGoalComplete, OpGoalAbandon,
		OpCapabilityAdd, OpCapabilityRemove, OpCapabilityUpdate,
		OpLimitationAdd, OpLimitationRemove, OpLimitationUpdate,
		OpAspirationAdd, OpAspirationRemove, OpAspirationUpdate,
		OpMetadataUpdate:
		return true
	default:
		return false
	}
}

// AllOperationTypes returns all valid operation types
func AllOperationTypes() []OperationType {
	return []OperationType{
		OpGoalCreate, OpGoalTransition, OpGoalUpdate, OpGoalComplete, OpGoalAbandon,
		OpCapabilityAdd, OpCapabilityRemove, OpCapabilityUpdate,
		OpLimitationAdd, OpLimitationRemove, OpLimitationUpdate,
		OpAspirationAdd, OpAspirationRemove, OpAspirationUpdate,
		OpMetadataUpdate,
	}
}
