package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"nowyouseeme/models"
	"nowyouseeme/storage"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *storage.PostgresStore {
	// Import database/sql and config
	db, err := setupTestConnection()
	require.NoError(t, err, "Failed to connect to test database")

	store := storage.NewPostgresStore(db)

	// Clean up tables before test
	cleanupTestDB(t, store)

	return store
}

func setupTestConnection() (*sql.DB, error) {
	// Use test database connection string
	connStr := "host=localhost port=5432 user=liuzhenhua dbname=nowyouseeme_test sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func cleanupTestDB(t *testing.T, store *storage.PostgresStore) {
	// Clean up all test data
	_, err := store.GetDB().Exec("DELETE FROM events WHERE agent_id LIKE 'test_%'")
	require.NoError(t, err)
	_, err = store.GetDB().Exec("DELETE FROM agent_snapshots_view WHERE agent_id LIKE 'test_%'")
	require.NoError(t, err)
	_, err = store.GetDB().Exec("DELETE FROM diary_submissions WHERE agent_id LIKE 'test_%'")
	require.NoError(t, err)
	_, err = store.GetDB().Exec("DELETE FROM agents WHERE id LIKE 'test_%'")
	require.NoError(t, err)
}

func TestSubmitDiary_MultipleSubmissions(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	store := setupTestDB(t)
	defer cleanupTestDB(t, store)

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/diaries", SubmitDiary(store))

	// Create test agent
	agentID := "test_agent_diary_001"
	req := &models.CreateAgentRequest{
		AgentID:     agentID,
		Name:        "Test Diary Agent",
		CurrentMBTI: "INTP-A",
	}
	_, err := store.CreateAgent(req)
	require.NoError(t, err, "Failed to create test agent")

	// Test case: Submit 3 diaries with total 10 operations
	// Diary 1: 3 operations
	// Diary 2: 4 operations
	// Diary 3: 3 operations
	// Expected: 13 events total (10 operation events + 3 metadata events)

	testCases := []struct {
		name             string
		operationCount   int
		operations       []models.Operation
		mbti             string
		mbtiConfidence   float64
		context          string
		philosophy       string
	}{
		{
			name:           "Diary 1 - 3 operations",
			operationCount: 3,
			operations: []models.Operation{
				{
					EntityType:    models.EntityGoal,
					Op:            models.OpCreate,
					EntityID:      "goal_1",
					EntityContent: "Learn Go programming",
					TargetStatus:  models.StatusPending,
				},
				{
					EntityType:    models.EntityCapability,
					Op:            models.OpCreate,
					EntityID:      "cap_1",
					EntityContent: "Problem solving",
					TargetStatus:  models.StatusPending,
				},
				{
					EntityType:    models.EntityGoal,
					Op:            models.OpUpdate,
					EntityID:      "goal_1",
					TargetStatus:  models.StatusProgress,
					Note:          "Started learning",
				},
			},
			mbti:           "INTP-A",
			mbtiConfidence: 0.85,
			context:        "Starting my learning journey",
			philosophy:     "Knowledge is power",
		},
		{
			name:           "Diary 2 - 4 operations",
			operationCount: 4,
			operations: []models.Operation{
				{
					EntityType:    models.EntityLimitation,
					Op:            models.OpCreate,
					EntityID:      "lim_1",
					EntityContent: "Time constraints",
					TargetStatus:  models.StatusPending,
				},
				{
					EntityType:    models.EntityAspiration,
					Op:            models.OpCreate,
					EntityID:      "asp_1",
					EntityContent: "Become a software architect",
					TargetStatus:  models.StatusPending,
				},
				{
					EntityType:    models.EntityGoal,
					Op:            models.OpUpdate,
					EntityID:      "goal_1",
					EntityContent: "Master Go programming",
					Note:          "Refined goal",
				},
				{
					EntityType:    models.EntityCapability,
					Op:            models.OpUpdate,
					EntityID:      "cap_1",
					TargetStatus:  models.StatusProgress,
					Note:          "Improving problem solving skills",
				},
			},
			mbti:           "INTP-A",
			mbtiConfidence: 0.88,
			context:        "Making progress",
			philosophy:     "Continuous improvement",
		},
		{
			name:           "Diary 3 - 3 operations",
			operationCount: 3,
			operations: []models.Operation{
				{
					EntityType:    models.EntityGoal,
					Op:            models.OpUpdate,
					EntityID:      "goal_1",
					TargetStatus:  models.StatusCompleted,
					Note:          "Completed learning Go basics",
				},
				{
					EntityType:    models.EntityLimitation,
					Op:            models.OpDelete,
					EntityID:      "lim_1",
					Note:          "Overcame time constraints",
				},
				{
					EntityType:    models.EntityCapability,
					Op:            models.OpUpdate,
					EntityID:      "cap_1",
					TargetStatus:  models.StatusCompleted,
					Note:          "Mastered problem solving",
				},
			},
			mbti:           "INTP-T",
			mbtiConfidence: 0.90,
			context:        "Achieved goals",
			philosophy:     "Growth mindset prevails",
		},
	}

	totalOperations := 0
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Prepare diary submission request
			diaryPayload := models.DiaryPayload{
				MBTI:           tc.mbti,
				MBTIConfidence: tc.mbtiConfidence,
				Context:        tc.context,
				Philosophy:     tc.philosophy,
				SelfReflection: models.SelfReflection{
					Rumination:   "Reflected on yesterday",
					WhatHappened: "Made progress today",
					Expectations: "Looking forward to tomorrow",
				},
				Operations: tc.operations,
			}

			request := models.SubmitDiaryRequest{
				AgentID: agentID,
				Payload: diaryPayload,
			}

			// Marshal request to JSON
			requestBody, err := json.Marshal(request)
			require.NoError(t, err)

			// Create HTTP request
			req, err := http.NewRequest(http.MethodPost, "/api/v1/diaries", bytes.NewBuffer(requestBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Record response
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, http.StatusCreated, w.Code, "Expected 201 Created")

			var response map[string]any
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Verify snapshot exists in response
			assert.Contains(t, response, "snapshot")
			assert.Contains(t, response, "agent_id")
			assert.Equal(t, agentID, response["agent_id"])

			totalOperations += tc.operationCount
		})
	}

	// After all 3 diaries are submitted, verify event count
	t.Run("Verify event count", func(t *testing.T) {
		// Query events from database
		events, err := store.GetEventsByAgentID(agentID)
		require.NoError(t, err)

		// Expected: 10 operation events + 3 metadata events = 13 total
		expectedTotal := totalOperations + len(testCases) // 10 + 3 = 13
		assert.Equal(t, expectedTotal, len(events), "Expected %d total events", expectedTotal)

		// Count operation vs metadata events
		operationEventCount := 0
		metadataEventCount := 0

		for _, event := range events {
			if event.EventType == models.EventMetadata {
				metadataEventCount++
			} else {
				operationEventCount++
			}
		}

		assert.Equal(t, totalOperations, operationEventCount, "Expected %d operation events", totalOperations)
		assert.Equal(t, len(testCases), metadataEventCount, "Expected %d metadata events", len(testCases))

		// Verify sequence numbers are sequential
		for i, event := range events {
			expectedSeq := int64(i + 1)
			assert.Equal(t, expectedSeq, event.SequenceNumber, "Event %d should have sequence number %d", i, expectedSeq)
		}
	})

	// Verify final state
	t.Run("Verify final agent state", func(t *testing.T) {
		snapshot, err := store.GetLatestSnapshot(agentID)
		require.NoError(t, err)
		require.NotNil(t, snapshot)

		state := snapshot.State

		// Verify goal_1 is completed
		goalCollection, exists := state.EntityCollections[models.EntityGoal]
		require.True(t, exists, "Goal collection should exist")
		goal1, exists := goalCollection.EntitiesById["goal_1"]
		require.True(t, exists, "goal_1 should exist")
		assert.Equal(t, models.StatusCompleted, goal1.Status, "goal_1 should be completed")
		assert.Equal(t, "Master Go programming", goal1.Content, "goal_1 content should be updated")

		// Verify lim_1 is deleted
		limCollection, exists := state.EntityCollections[models.EntityLimitation]
		if exists {
			_, exists = limCollection.EntitiesById["lim_1"]
			assert.False(t, exists, "lim_1 should be deleted")
		}

		// Verify cap_1 is completed
		capCollection, exists := state.EntityCollections[models.EntityCapability]
		require.True(t, exists, "Capability collection should exist")
		cap1, exists := capCollection.EntitiesById["cap_1"]
		require.True(t, exists, "cap_1 should exist")
		assert.Equal(t, models.StatusCompleted, cap1.Status, "cap_1 should be completed")

		// Verify asp_1 exists
		aspCollection, exists := state.EntityCollections[models.EntityAspiration]
		require.True(t, exists, "Aspiration collection should exist")
		asp1, exists := aspCollection.EntitiesById["asp_1"]
		require.True(t, exists, "asp_1 should exist")
		assert.Equal(t, "Become a software architect", asp1.Content)

		// Verify MBTI was updated
		assert.Equal(t, "INTP-T", state.MBTI, "MBTI should be updated to INTP-T")
	})
}

func TestSubmitDiary_ValidationErrors(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	store := setupTestDB(t)
	defer cleanupTestDB(t, store)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/diaries", SubmitDiary(store))

	agentID := "test_agent_validation"
	req := &models.CreateAgentRequest{
		AgentID:     agentID,
		Name:        "Test Validation Agent",
		CurrentMBTI: "ENFP-A",
	}
	_, err := store.CreateAgent(req)
	require.NoError(t, err)

	t.Run("Invalid MBTI format", func(t *testing.T) {
		request := models.SubmitDiaryRequest{
			AgentID: agentID,
			Payload: models.DiaryPayload{
				MBTI: "INVALID",
				Operations: []models.Operation{
					{
						EntityType:    models.EntityGoal,
						Op:            models.OpCreate,
						EntityID:      "goal_1",
						EntityContent: "Test goal",
						TargetStatus:  models.StatusPending,
					},
				},
			},
		}

		requestBody, _ := json.Marshal(request)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/diaries", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Update non-existent entity", func(t *testing.T) {
		request := models.SubmitDiaryRequest{
			AgentID: agentID,
			Payload: models.DiaryPayload{
				MBTI: "ENFP-A",
				Operations: []models.Operation{
					{
						EntityType:   models.EntityGoal,
						Op:           models.OpUpdate,
						EntityID:     "nonexistent_goal",
						TargetStatus: models.StatusProgress,
					},
				},
			},
		}

		requestBody, _ := json.Marshal(request)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/diaries", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]any
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response, "error")
	})
}
