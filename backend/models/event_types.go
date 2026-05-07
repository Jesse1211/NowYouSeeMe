package models

// EventType represents the type of CRUD Event
type EventType string

// Event type constants - simplified to CRUD Events
const (
	EventCreate   EventType = "create"
	EventUpdate   EventType = "update"
	EventDelete   EventType = "delete"
	EventMetadata EventType = "metadata_submission"
)

// String returns the string representation of the Event type
func (o EventType) String() string {
	return string(o)
}

// IsValid checks if the Event type is valid
func (o EventType) IsValid() bool {
	switch o {
	case EventCreate, EventUpdate, EventDelete, EventMetadata:
		return true
	default:
		return false
	}
}

// AllEventTypes returns all valid Event types
func AllEventTypes() []EventType {
	return []EventType{EventCreate, EventUpdate, EventDelete, EventMetadata}
}

func OperationTypeToEventType(operationType OperationType) EventType {
	switch operationType {
	case OpCreate:
		return EventCreate
	case OpUpdate:
		return EventUpdate
	case OpDelete:
		return EventDelete
	default:
		return ""
	}
}
