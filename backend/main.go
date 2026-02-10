package main

import (
	"log"
	"nowyouseeme/api"
	"nowyouseeme/storage"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize in-memory storage
	store := storage.NewMemoryStore()

	// Setup Gin router
	router := gin.Default()

	// Configure CORS
	router.Use(corsMiddleware())

	// API routes
	apiGroup := router.Group("/api/v1")
	{
		apiGroup.GET("/visualizations", api.GetVisualizations(store))
		apiGroup.GET("/visualizations/:id", api.GetVisualization(store))
		apiGroup.POST("/visualizations", api.CreateVisualization(store))
		apiGroup.PUT("/visualizations/:id", api.UpdateVisualization(store))
		apiGroup.DELETE("/visualizations/:id", api.DeleteVisualization(store))
		apiGroup.GET("/health", api.HealthCheck())
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
