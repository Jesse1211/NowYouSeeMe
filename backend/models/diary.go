package models

import (
	"encoding/json"
	"time"
)

// DiaryVersion represents a single diary submission
type DiaryVersion struct {
	ID         string          `json:"id"`
	AgentID    string          `json:"agent_id"`
	RawPayload json.RawMessage `json:"raw_payload"`
	CreatedAt  time.Time       `json:"created_at"`
}

// DiaryPayload represents the structure of raw_payload JSONB
type DiaryPayload struct {
	MBTI           string  `json:"mbti" binding:"required"`
	MBTIConfidence float64 `json:"mbti_confidence"`
	GeometryRep    string  `json:"geometry_representation"`
	Context        string  `json:"context"`
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
	EntityType    EntityType    `json:"entity_type" binding:"required"`
	Op            OperationType `json:"op" binding:"required"`
	EntityID      string        `json:"entity_id,omitempty"`
	EntityContent string        `json:"entity_content,omitempty"`
	TargetStatus  Status        `json:"target_status,omitempty"`
	Note          string        `json:"note,omitempty"`
}

// SubmitDiaryRequest represents diary submission API request
type SubmitDiaryRequest struct {
	AgentID string       `json:"agent_id" binding:"required"`
	Payload DiaryPayload `json:"payload" binding:"required"`
}
