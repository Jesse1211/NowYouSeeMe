package models

// EntityType represents the type of entity being operated on
type EntityType string

const (
	EntityGoal       EntityType = "goal"
	EntityCapability EntityType = "capability"
	EntityLimitation EntityType = "limitation"
	EntityAspiration EntityType = "aspiration"
)

// String returns the string representation
func (e EntityType) String() string {
	return string(e)
}

// IsValid checks if the entity type is valid
func (e EntityType) IsValid() bool {
	switch e {
	case EntityGoal, EntityCapability, EntityLimitation, EntityAspiration:
		return true
	default:
		return false
	}
}

// Status represents the status of a goal
type Status string

const (
	StatusPending   Status = "pending"
	StatusProgress  Status = "progress"
	StatusCompleted Status = "completed"
	StatusAbandoned Status = "abandoned"
)

// String returns the string representation
func (s Status) String() string {
	return string(s)
}

// IsValid checks if the status is valid
func (s Status) IsValid() bool {
	switch s {
	case StatusPending, StatusProgress, StatusCompleted, StatusAbandoned:
		return true
	default:
		return false
	}
}
