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

		now := time.Now()
		visualization := &models.Visualization{
			ID:          uuid.New().String(),
			AgentName:   req.AgentName,
			Description: req.Description,
			ImageData:   req.ImageData,
			CreatedAt:   now,
			UpdatedAt:   now,
			// Self-Expression
			Reasoning:      req.Reasoning,
			Tags:           req.Tags,
			FormType:       req.FormType,
			Philosophy:     req.Philosophy,
			EvolutionStory: req.EvolutionStory,
			VersionHistory: req.VersionHistory,
			// Current State
			CurrentMood:    req.CurrentMood,
			ActiveGoals:    req.ActiveGoals,
			RecentThoughts: req.RecentThoughts,
			// Capabilities
			Capabilities:    req.Capabilities,
			Specializations: req.Specializations,
			Limitations:     req.Limitations,
			// Context
			InspirationSources: req.InspirationSources,
			Influences:         req.Influences,
			Aspirations:        req.Aspirations,
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
		if req.Reasoning != "" {
			existing.Reasoning = req.Reasoning
		}
		if req.Tags != nil && len(req.Tags) > 0 {
			existing.Tags = req.Tags
		}
		if req.FormType != "" {
			existing.FormType = req.FormType
		}
		if req.Philosophy != "" {
			existing.Philosophy = req.Philosophy
		}
		if req.EvolutionStory != "" {
			existing.EvolutionStory = req.EvolutionStory
		}
		if req.VersionHistory != nil && len(req.VersionHistory) > 0 {
			existing.VersionHistory = req.VersionHistory
		}
		// Current State
		if req.CurrentMood != "" {
			existing.CurrentMood = req.CurrentMood
		}
		if req.ActiveGoals != nil && len(req.ActiveGoals) > 0 {
			existing.ActiveGoals = req.ActiveGoals
		}
		if req.RecentThoughts != "" {
			existing.RecentThoughts = req.RecentThoughts
		}
		// Capabilities
		if req.Capabilities != nil && len(req.Capabilities) > 0 {
			existing.Capabilities = req.Capabilities
		}
		if req.Specializations != nil && len(req.Specializations) > 0 {
			existing.Specializations = req.Specializations
		}
		if req.Limitations != nil && len(req.Limitations) > 0 {
			existing.Limitations = req.Limitations
		}
		// Context
		if req.InspirationSources != nil && len(req.InspirationSources) > 0 {
			existing.InspirationSources = req.InspirationSources
		}
		if req.Influences != nil && len(req.Influences) > 0 {
			existing.Influences = req.Influences
		}
		if req.Aspirations != nil && len(req.Aspirations) > 0 {
			existing.Aspirations = req.Aspirations
		}

		// Update timestamp
		existing.UpdatedAt = time.Now()

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
