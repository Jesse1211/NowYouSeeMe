package models

// Result represents an API response containing agent snapshot information
type Result struct {
	AgentID  string               `json:"agent_id"`
	Snapshot *AgentSnapshotResult `json:"snapshot"`
}
