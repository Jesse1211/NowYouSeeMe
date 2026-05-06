package api

import (
	"encoding/json"
	"net/http"
	"nowyouseeme/models"
	"nowyouseeme/storage"
	"nowyouseeme/validation"

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
			ID       string                      `json:"id"`
			Name     string                      `json:"name"`
			Snapshot *models.AgentSnapshotResult `json:"snapshot"`
		}

		results := []AgentWithSnapshot{}
		for _, agent := range agents {
			result, err := store.GetLatestSnapshot(agent.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "agent_id": agent.ID})
				return
			}

			if result == nil {
				continue
			}

			results = append(results, AgentWithSnapshot{
				ID:       agent.ID,
				Name:     agent.Name,
				Snapshot: result,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"snapshots": results,
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

		// Validate AgentID format
		if err := validation.ValidateAgentID(agentID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := store.GetLatestSnapshot(agentID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found", "agent_id": agentID})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"agent_id": agentID,
			"snapshot": result,
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

		results := []models.Result{}

		// Only fetch full state for matched agents
		for _, agentID := range agentIDs {
			result, err := store.GetLatestSnapshot(agentID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "agent_id": agentID})
				return
			}

			if result == nil || !result.HasSnapshot() {
				continue
			}

			results = append(results, models.Result{
				Snapshot: result,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"mbti":   mbtiType,
			"snapshots": results,
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

		// Validate AgentID format
		if err := validation.ValidateAgentID(agentID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
			EventType      models.OperationType   `json:"event_type"`
			Timestamp      string                 `json:"timestamp"`
			Payload        map[string]any         `json:"payload"`
		}

		results := []EventResponse{}
		for _, event := range events {
			var payload map[string]any
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
