package api

import (
	"net/http"
	"nowyouseeme/models"
	"nowyouseeme/storage"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetVisualizations returns all visualizations
func GetVisualizations(store *storage.MemoryStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		visualizations := store.GetAllVisualizations()
		c.JSON(http.StatusOK, gin.H{
			"visualizations": visualizations,
			"count":          len(visualizations),
		})
	}
}

// GetVisualization returns a specific visualization by ID
func GetVisualization(store *storage.MemoryStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		visualization, err := store.GetVisualization(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Visualization not found"})
			return
		}

		c.JSON(http.StatusOK, visualization)
	}
}

// CreateVisualization creates a new visualization
func CreateVisualization(store *storage.MemoryStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreateVisualizationRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		visualization := &models.Visualization{
			ID:          uuid.New().String(),
			AgentName:   req.AgentName,
			Description: req.Description,
			ImageData:   req.ImageData,
			CreatedAt:   time.Now(),
		}

		if err := store.CreateVisualization(visualization); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, visualization)
	}
}

// UpdateVisualization updates an existing visualization
func UpdateVisualization(store *storage.MemoryStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		// Check if visualization exists
		existing, err := store.GetVisualization(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Visualization not found"})
			return
		}

		var req models.UpdateVisualizationRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update fields (keep existing values if not provided)
		if req.AgentName != "" {
			existing.AgentName = req.AgentName
		}
		if req.Description != "" {
			existing.Description = req.Description
		}
		if req.ImageData != "" {
			existing.ImageData = req.ImageData
		}

		if err := store.UpdateVisualization(id, existing); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, existing)
	}
}

// DeleteVisualization deletes a visualization by ID
func DeleteVisualization(store *storage.MemoryStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := store.DeleteVisualization(id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Visualization not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Visualization deleted successfully",
			"id":      id,
		})
	}
}

// HealthCheck returns server health status
func HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now(),
		})
	}
}
