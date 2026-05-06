package models

import (
	"encoding/json"
	"time"
)

// DiaryVersion represents a single diary submission
type DiaryVersion struct {
	ID          string          `json:"id"`
	AgentID     string          `json:"agent_id"`
	RawPayload  json.RawMessage `json:"raw_payload"`
	CreatedAt   time.Time       `json:"created_at"`
	ParsedError *string         `json:"parsed_error,omitempty"`
}

// DiaryPayload represents the structure of raw_payload JSONB
type DiaryPayload struct {
	MBTI           string  `json:"mbti" binding:"required"`
	MBTIConfidence float64 `json:"mbti_confidence"`
	GeometryRep    string  `json:"geometry_representation"`
	Reasoning      string  `json:"reasoning"`
	CurrentMood    string  `json:"current_mood"`
	Philosophy     string  `json:"philosophy"`

	SelfReflection SelfReflection `json:"self_reflection"`
	Operations     []Operation    `json:"operations"`
}

// SelfReflection holds agent's reflective thoughts
type SelfReflection struct {
	Rumination   string `json:"rumination_for_yesterday"`
	WhatHappened string `json:"what_happened_today"`
	Expectations string `json:"expectations_for_tomorrow"`
}

// Operation represents a single state-changing operation
type Operation struct {
	Op string `json:"op" binding:"required"`

	// Goal operations
	GoalID     string `json:"goal_id,omitempty"`
	Title      string `json:"title,omitempty"`
	Status     string `json:"status,omitempty"`
	FromStatus string `json:"from_status,omitempty"`
	ToStatus   string `json:"to_status,omitempty"`
	Reason     string `json:"reason,omitempty"`
	Checkpoint string `json:"checkpoint,omitempty"`

	// Capability operations
	CapabilityID string `json:"capability_id,omitempty"`

	// Limitation operations
	LimitationID string `json:"limitation_id,omitempty"`

	// Aspiration operations
	AspirationID string `json:"aspiration_id,omitempty"`
}

// SubmitDiaryRequest represents diary submission API request
type SubmitDiaryRequest struct {
	AgentID string       `json:"agent_id" binding:"required"`
	Payload DiaryPayload `json:"payload" binding:"required"`
}
