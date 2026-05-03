package validation

import (
	"fmt"
	"regexp"
)

// agentIDRegex defines valid AgentID format: alphanumeric, underscore, hyphen
var agentIDRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

const (
	// MaxAgentIDLength defines maximum allowed length for AgentID
	MaxAgentIDLength = 100
	// MinAgentIDLength defines minimum allowed length for AgentID
	MinAgentIDLength = 1
)

// ValidateAgentID validates an agent ID string
// Returns nil if valid, error otherwise
//
// Requirements:
// - Length: 1-100 characters
// - Format: alphanumeric characters, underscores, and hyphens only
// - No special characters, emojis, or control characters
func ValidateAgentID(agentID string) error {
	if agentID == "" {
		return fmt.Errorf("agent_id cannot be empty")
	}

	if len(agentID) < MinAgentIDLength {
		return fmt.Errorf("agent_id too short: minimum %d characters, got %d", MinAgentIDLength, len(agentID))
	}

	if len(agentID) > MaxAgentIDLength {
		return fmt.Errorf("agent_id too long: maximum %d characters, got %d", MaxAgentIDLength, len(agentID))
	}

	if !agentIDRegex.MatchString(agentID) {
		return fmt.Errorf("agent_id contains invalid characters: only alphanumeric, underscore (_), and hyphen (-) are allowed, got '%s'", agentID)
	}

	return nil
}
