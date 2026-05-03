package api

import (
	"encoding/json"
	"net/http"
	"nowyouseeme/models"
	"nowyouseeme/storage"

	"github.com/gin-gonic/gin"
)

// GetGallery returns all agents with current snapshots
func GetGallery(store *storage.PostgresStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		agents, err := store.GetAllAgents()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		type AgentWithSnapshot struct {
			ID          string             `json:"id"`
			Name        string             `json:"name"`
			Snapshot    *models.AgentState `json:"snapshot"`
			LastUpdated string             `json:"last_updated"`
		}

		results := []AgentWithSnapshot{}
		for _, agent := range agents {
			snapshot, _, err := store.GetLatestSnapshot(agent.ID)
			if err != nil {
				continue // Skip agents with errors
			}

			snapshotDB, _ := store.GetSnapshot(agent.ID)
			lastUpdated := ""
			if snapshotDB != nil {
				lastUpdated = snapshotDB.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
			}

			results = append(results, AgentWithSnapshot{
				ID:          agent.ID,
				Name:        agent.Name,
				Snapshot:    snapshot,
				LastUpdated: lastUpdated,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"agents": results,
			"total":  len(results),
		})
	}
}

// GetSnapshot returns current snapshot for an agent
func GetSnapshot(store *storage.PostgresStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		agentID := c.Query("agent_id")
		if agentID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "agent_id query parameter required"})
			return
		}

		snapshot, _, err := store.GetLatestSnapshot(agentID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found", "agent_id": agentID})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"agent_id": agentID,
			"snapshot": snapshot,
		})
	}
}

// GetSnapshotsByMBTI filters agents by MBTI type
func GetSnapshotsByMBTI(store *storage.PostgresStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		mbtiType := c.Query("mbti")
		if mbtiType == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "mbti query parameter required"})
			return
		}

		// Only query agents with matching current MBTI
		agentIDs, err := store.GetAgentIDsByCurrentMBTI(mbtiType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		type Result struct {
			AgentID   string             `json:"agent_id"`
			Snapshot  *models.AgentState `json:"snapshot"`
			UpdatedAt string             `json:"updated_at"`
		}

		results := []Result{}

		// Only fetch full state for matched agents
		for _, agentID := range agentIDs {
			snapshot, _, err := store.GetLatestSnapshot(agentID)
			if err != nil {
				continue
			}

			snapshotDB, _ := store.GetSnapshot(agentID)
			updatedAt := ""
			if snapshotDB != nil {
				updatedAt = snapshotDB.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
			}

			results = append(results, Result{
				AgentID:   agentID,
				Snapshot:  snapshot,
				UpdatedAt: updatedAt,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"mbti":   mbtiType,
			"agents": results,
			"count":  len(results),
		})
	}
}

// GetTimeline returns event stream for an agent
func GetTimeline(store *storage.PostgresStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		agentID := c.Query("agent_id")
		if agentID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "agent_id query parameter required"})
			return
		}

		events, err := store.GetEventsByAgent(agentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		type EventResponse struct {
			EventID        int64                  `json:"event_id"`
			SequenceNumber int64                  `json:"sequence_number"`
			EventType      string                 `json:"event_type"`
			Timestamp      string                 `json:"timestamp"`
			Payload        map[string]interface{} `json:"payload"`
		}

		results := []EventResponse{}
		for _, event := range events {
			var payload map[string]interface{}
			json.Unmarshal(event.Payload, &payload)

			results = append(results, EventResponse{
				EventID:        event.EventID,
				SequenceNumber: event.SequenceNumber,
				EventType:      event.EventType,
				Timestamp:      event.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
				Payload:        payload,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"agent_id":     agentID,
			"events":       results,
			"total_events": len(results),
		})
	}
}

// HealthCheck returns server health status
func HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	}
}
