package validation

import "fmt"

// ValidGoalTransitions defines allowed state transitions
var ValidGoalTransitions = map[string][]string{
	"future":      {"progressing", "abandoned"},
	"progressing": {"completed", "abandoned", "future"},
	"abandoned":   {"future", "progressing"},
	"completed":   {}, // terminal state
}

// ValidateGoalTransition checks if a goal transition is valid
func ValidateGoalTransition(fromStatus, toStatus string) error {
	allowedTransitions, exists := ValidGoalTransitions[fromStatus]
	if !exists {
		return fmt.Errorf("invalid goal status: %s", fromStatus)
	}

	for _, allowed := range allowedTransitions {
		if allowed == toStatus {
			return nil
		}
	}

	return fmt.Errorf("invalid transition: %s → %s", fromStatus, toStatus)
}

// ValidateGoalCreateStatus checks if status is valid for goal creation
func ValidateGoalCreateStatus(status string) error {
	if status != "future" && status != "progressing" {
		return fmt.Errorf("goal_create status must be 'future' or 'progressing' (MVP limitation), got: %s", status)
	}
	return nil
}
