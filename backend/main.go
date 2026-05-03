package main

import (
	"log"
	"nowyouseeme/api"
	"nowyouseeme/config"
	"nowyouseeme/storage"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file (ignore error if file doesn't exist)
	_ = godotenv.Load()
	_ = godotenv.Load("../.env") // Also try root .env

	// Load database configuration
	dbConfig := config.LoadDBConfig()

	// Connect to PostgreSQL
	db, err := config.ConnectDB(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize PostgreSQL storage
	store := storage.NewPostgresStore(db)

	// Setup Gin router
	router := gin.Default()

	// Configure CORS
	router.Use(corsMiddleware())

	// API routes
	apiGroup := router.Group("/api/v1")
	{
		// Agent management
		apiGroup.POST("/agents", api.CreateAgent(store))
		apiGroup.GET("/agents", api.GetAgents(store))

		// Diary submission
		apiGroup.POST("/diaries", api.SubmitDiary(store))

		// Gallery & discovery
		apiGroup.GET("/gallery", api.GetGallery(store))
		apiGroup.GET("/snapshot", api.GetSnapshot(store))
		apiGroup.GET("/snapshots", api.GetSnapshotsByMBTI(store))
		apiGroup.GET("/timeline", api.GetTimeline(store))

		// Health check
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
