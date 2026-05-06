package models

// OperationType represents the type of CRUD operation
type OperationType string

// Operation type constants - simplified to CRUD operations
const (
	OpCreate OperationType = "create"
	OpUpdate OperationType = "update"
	OpDelete OperationType = "delete"
)

// String returns the string representation of the operation type
func (o OperationType) String() string {
	return string(o)
}

// IsValid checks if the operation type is valid
func (o OperationType) IsValid() bool {
	switch o {
	case OpCreate, OpUpdate, OpDelete:
		return true
	default:
		return false
	}
}

// AllOperationTypes returns all valid operation types
func AllOperationTypes() []OperationType {
	return []OperationType{OpCreate, OpUpdate, OpDelete}
}
