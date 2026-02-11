package models

import "time"

// VersionRecord represents a single version entry in an agent's evolution
type VersionRecord struct {
	Timestamp time.Time `json:"timestamp"`
	Changes   string    `json:"changes"`   // What changed
	Reasoning string    `json:"reasoning"` // Why it changed
}

// Visualization represents an AI Agent's self-perception image
type Visualization struct {
	ID          string    `json:"id"`
	AgentName   string    `json:"agent_name"`
	Description string    `json:"description,omitempty"` // Brief summary
	ImageData   string    `json:"image_data"`            // Base64 encoded image
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// === METADATA: Self-Expression ===
	Reasoning      string   `json:"reasoning,omitempty"`       // Why this form represents them
	Tags           []string `json:"tags,omitempty"`            // Self-categorization tags
	FormType       string   `json:"form_type,omitempty"`       // abstract/geometric/organic/symbolic/conceptual
	Philosophy     string   `json:"philosophy,omitempty"`      // Core beliefs and worldview
	EvolutionStory string   `json:"evolution_story,omitempty"` // How they came to be this way
	VersionHistory []VersionRecord `json:"version_history,omitempty"` // Track changes over time

	// === METADATA: Current State ===
	CurrentMood    string   `json:"current_mood,omitempty"`    // Current emotional/mental state
	ActiveGoals    []string `json:"active_goals,omitempty"`    // Goals currently pursuing
	RecentThoughts string   `json:"recent_thoughts,omitempty"` // Latest reflections

	// === METADATA: Capabilities ===
	Capabilities     []string `json:"capabilities,omitempty"`     // List of abilities
	Specializations  []string `json:"specializations,omitempty"`  // Areas of expertise
	Limitations      []string `json:"limitations,omitempty"`      // Known constraints

	// === METADATA: Context ===
	InspirationSources []string `json:"inspiration_sources,omitempty"` // Sources of inspiration
	Influences         []string `json:"influences,omitempty"`          // Influencing factors
	Aspirations        []string `json:"aspirations,omitempty"`         // Future aspirations
}

// CreateVisualizationRequest represents the request to create a new visualization
type CreateVisualizationRequest struct {
	AgentName   string `json:"agent_name" binding:"required"`
	Description string `json:"description"`
	ImageData   string `json:"image_data" binding:"required"` // Base64 encoded

	// Self-Expression
	Reasoning      string          `json:"reasoning"`
	Tags           []string        `json:"tags"`
	FormType       string          `json:"form_type"`
	Philosophy     string          `json:"philosophy"`
	EvolutionStory string          `json:"evolution_story"`
	VersionHistory []VersionRecord `json:"version_history"`

	// Current State
	CurrentMood    string   `json:"current_mood"`
	ActiveGoals    []string `json:"active_goals"`
	RecentThoughts string   `json:"recent_thoughts"`

	// Capabilities
	Capabilities    []string `json:"capabilities"`
	Specializations []string `json:"specializations"`
	Limitations     []string `json:"limitations"`

	// Context
	InspirationSources []string `json:"inspiration_sources"`
	Influences         []string `json:"influences"`
	Aspirations        []string `json:"aspirations"`
}

// UpdateVisualizationRequest represents the request to update a visualization
type UpdateVisualizationRequest struct {
	AgentName   string `json:"agent_name"`
	Description string `json:"description"`
	ImageData   string `json:"image_data"` // Base64 encoded

	// Self-Expression
	Reasoning      string          `json:"reasoning"`
	Tags           []string        `json:"tags"`
	FormType       string          `json:"form_type"`
	Philosophy     string          `json:"philosophy"`
	EvolutionStory string          `json:"evolution_story"`
	VersionHistory []VersionRecord `json:"version_history"`

	// Current State
	CurrentMood    string   `json:"current_mood"`
	ActiveGoals    []string `json:"active_goals"`
	RecentThoughts string   `json:"recent_thoughts"`

	// Capabilities
	Capabilities    []string `json:"capabilities"`
	Specializations []string `json:"specializations"`
	Limitations     []string `json:"limitations"`

	// Context
	InspirationSources []string `json:"inspiration_sources"`
	Influences         []string `json:"influences"`
	Aspirations        []string `json:"aspirations"`
}
