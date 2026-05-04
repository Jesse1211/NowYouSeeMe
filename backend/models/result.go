package models

// Result represents an API response containing agent snapshot information
type Result struct {
	Snapshot *AgentSnapshotResult `json:"snapshot"`
}
