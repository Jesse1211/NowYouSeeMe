package validation

import (
	"fmt"
	"nowyouseeme/models"
)

// ValidateGoalStatusTransition validates if a status transition is allowed
func ValidateGoalStatusTransition(from, to models.Status) error {
	// Allow same status (no-op)
	if from == to {
		return nil
	}

	validTransitions := map[models.Status][]models.Status{
		models.StatusPending: {
			models.StatusProgress,
			models.StatusAbandoned,
		},
		models.StatusProgress: {
			models.StatusCompleted,
			models.StatusAbandoned,
			models.StatusPending, // Allow moving back to pending
		},
		models.StatusCompleted: {
			// Completed is terminal - no transitions allowed
		},
		models.StatusAbandoned: {
			models.StatusPending,  // Allow reactivation
			models.StatusProgress, // Allow reactivation directly to progress
		},
	}

	allowed, exists := validTransitions[from]
	if !exists {
		return fmt.Errorf("unknown status: %s", from)
	}

	for _, allowedStatus := range allowed {
		if to == allowedStatus {
			return nil
		}
	}

	return fmt.Errorf("invalid status transition: %s -> %s", from, to)
}
