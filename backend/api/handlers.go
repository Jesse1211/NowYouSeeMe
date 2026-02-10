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

// HealthCheck returns server health status
func HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now(),
		})
	}
}
