package api

import (
	"net/http"
	"nowyouseeme/models"
	"nowyouseeme/storage"
	"nowyouseeme/validation"

	"github.com/gin-gonic/gin"
)

// SubmitDiary handles diary submission
func SubmitDiary(store *storage.PostgresStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.SubmitDiaryRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate MBTI format
		if err := validation.ValidateMBTI(req.Payload.MBTI); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check agent exists
		_, err := store.GetAgent(req.AgentID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found", "agent_id": req.AgentID})
			return
		}

		// Submit diary (validates, creates events, materializes snapshot)
		snapshot, err := store.SubmitDiary(req.AgentID, &req.Payload)
		if err != nil {
			// Check if validation error
			if validationErr, ok := err.(*validation.ValidationError); ok {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Validation failed",
					"details": validationErr.Errors,
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"agent_id": req.AgentID,
			"snapshot": snapshot,
		})
	}
}
