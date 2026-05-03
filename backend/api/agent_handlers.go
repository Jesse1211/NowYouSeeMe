package api

import (
	"net/http"
	"nowyouseeme/models"
	"nowyouseeme/storage"
	"nowyouseeme/validation"

	"github.com/gin-gonic/gin"
)

// CreateAgent creates a new agent
func CreateAgent(store *storage.PostgresStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreateAgentRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate MBTI format
		if err := validation.ValidateMBTI(req.InitialMBTI); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		agent, err := store.CreateAgent(&req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, agent)
	}
}

// GetAgents lists all agents with optional filtering
func GetAgents(store *storage.PostgresStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		agentID := c.Query("agent_id")

		// Get specific agent
		if agentID != "" {
			agent, err := store.GetAgent(agentID)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found", "agent_id": agentID})
				return
			}

			// Get current snapshot with WAL
			snapshot, _, err := store.GetCurrentState(agentID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"agent":    agent,
				"snapshot": snapshot,
			})
			return
		}

		// List all agents
		agents, err := store.GetAllAgents()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"agents": agents,
			"count":  len(agents),
		})
	}
}
