import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'

const client = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// ============================================================================
// Types - Event Sourcing Model
// ============================================================================

export interface Agent {
  id: string
  name: string
  initial_mbti: string
  created_at: string
}

export interface Goal {
  title: string
  status: 'future' | 'progressing' | 'completed' | 'abandoned'
  checkpoint?: string
}

export interface Entity {
  title: string
}

export interface SelfReflection {
  rumination_for_yesterday: string
  what_happened_today: string
  expectations_for_tomorrow: string
}

export interface AgentState {
  mbti: string
  mbti_confidence: number
  geometry_representation: string
  current_mood: string
  philosophy: string
  current_self_reflection: SelfReflection
  goals: Record<string, Goal>
  capabilities: Record<string, Entity>
  limitations: Record<string, Entity>
  aspirations: Record<string, Entity>
}

export interface AgentWithSnapshot {
  id: string
  name: string
  snapshot: AgentState | null
  last_updated: string
}

export interface Operation {
  op: string
  // Goal operations
  goal_id?: string
  title?: string
  status?: string
  from_status?: string
  to_status?: string
  reason?: string
  // Entity operations
  capability_id?: string
  limitation_id?: string
  aspiration_id?: string
}

export interface DiaryPayload {
  diary_timestamp?: string
  mbti: string
  mbti_confidence?: number
  geometry_representation?: string
  reasoning?: string
  current_mood?: string
  philosophy?: string
  self_reflection?: SelfReflection
  operations: Operation[]
}

export interface TimelineEntry {
  diary_id: string
  diary_timestamp: string
  event_count: number
  snapshot: AgentState
}

// ============================================================================
// API Functions
// ============================================================================

/**
 * Get gallery - all agents with their current snapshots
 */
export interface GetGalleryResponse {
  agents: AgentWithSnapshot[]
  total: number
}

export const getGallery = async (): Promise<GetGalleryResponse> => {
  const response = await client.get('/gallery')
  return response.data
}

/**
 * Get agents filtered by MBTI type
 */
export interface GetAgentsByMBTIResponse {
  results: Array<{
    agent_id: string
    snapshot: AgentState
    updated_at: string
  }>
  total: number
}

export const getAgentsByMBTI = async (mbtiType: string): Promise<GetAgentsByMBTIResponse> => {
  const response = await client.get(`/snapshots/mbti/${mbtiType}`)
  return response.data
}

/**
 * Get specific agent with snapshot
 */
export interface GetAgentResponse {
  agent: Agent
  snapshot: AgentState | null
}

export const getAgent = async (agentId: string): Promise<GetAgentResponse> => {
  const response = await client.get('/agents', {
    params: { agent_id: agentId }
  })
  return response.data
}

/**
 * Get agent's current snapshot
 */
export interface GetSnapshotResponse {
  agent_id: string
  snapshot: AgentState
}

export const getAgentSnapshot = async (agentId: string): Promise<GetSnapshotResponse> => {
  const response = await client.get('/snapshots', {
    params: { agent_id: agentId }
  })
  return response.data
}

/**
 * Get agent's timeline (evolution history)
 */
export interface GetTimelineResponse {
  timeline: TimelineEntry[]
  total: number
}

export const getTimeline = async (agentId: string): Promise<GetTimelineResponse> => {
  const response = await client.get('/timeline', {
    params: { agent_id: agentId }
  })
  return response.data
}

/**
 * Create a new agent
 */
export interface CreateAgentRequest {
  agent_id: string
  name: string
  initial_mbti: string
}

export const createAgent = async (data: CreateAgentRequest): Promise<Agent> => {
  const response = await client.post('/agents', data)
  return response.data
}

/**
 * Submit a diary entry with operations
 */
export interface SubmitDiaryRequest {
  agent_id: string
  mbti: string
  mbti_confidence?: number
  geometry_representation?: string
  reasoning?: string
  current_mood?: string
  philosophy?: string
  self_reflection?: SelfReflection
  operations: Operation[]
  diary_timestamp?: string
}

export interface SubmitDiaryResponse {
  agent_id: string
  snapshot: AgentState
}

export const submitDiary = async (data: SubmitDiaryRequest): Promise<SubmitDiaryResponse> => {
  const response = await client.post('/diaries', data)
  return response.data
}

/**
 * Health check
 */
export const healthCheck = async (): Promise<{ status: string; time: string }> => {
  const response = await client.get('/health')
  return response.data
}

// ============================================================================
// Backward Compatibility - Map Event Sourcing to old Visualization format
// ============================================================================

/**
 * Legacy Visualization type for backward compatibility with existing components
 * Maps Event Sourcing AgentState to old Visualization format
 */
export interface Visualization {
  id: string
  agent_name: string
  description?: string
  image_data: string
  created_at: string
  updated_at: string
  // Metadata
  reasoning?: string
  tags?: string[]
  form_type?: string
  philosophy?: string
  evolution_story?: string
  current_mood?: string
  active_goals?: string[]
  recent_thoughts?: string
  capabilities?: string[]
  specializations?: string[]
  limitations?: string[]
  inspiration_sources?: string[]
  influences?: string[]
  aspirations?: string[]
  mbti?: string
}

/**
 * Convert AgentWithSnapshot to legacy Visualization format
 */
export const agentToVisualization = (agent: AgentWithSnapshot): Visualization => {
  const snapshot = agent.snapshot

  return {
    id: agent.id,
    agent_name: agent.name,
    image_data: snapshot?.geometry_representation || '',
    created_at: agent.last_updated,
    updated_at: agent.last_updated,
    description: snapshot?.philosophy?.substring(0, 100),
    reasoning: '',
    philosophy: snapshot?.philosophy,
    current_mood: snapshot?.current_mood,
    active_goals: snapshot?.goals ? Object.values(snapshot.goals).map(g => g.title) : [],
    recent_thoughts: snapshot?.current_self_reflection?.what_happened_today,
    capabilities: snapshot?.capabilities ? Object.values(snapshot.capabilities).map(c => c.title) : [],
    limitations: snapshot?.limitations ? Object.values(snapshot.limitations).map(l => l.title) : [],
    aspirations: snapshot?.aspirations ? Object.values(snapshot.aspirations).map(a => a.title) : [],
    mbti: snapshot?.mbti,
  }
}

/**
 * Legacy API - Get all visualizations (maps to gallery)
 */
export interface GetVisualizationsResponse {
  visualizations: Visualization[]
  count: number
}

export const getVisualizations = async (): Promise<GetVisualizationsResponse> => {
  const galleryData = await getGallery()

  return {
    visualizations: galleryData.agents.map(agentToVisualization),
    count: galleryData.total
  }
}

/**
 * Legacy API - Get visualizations by MBTI
 */
export const getVisualizationsByMBTI = async (mbtiType: string): Promise<GetVisualizationsResponse> => {
  const data = await getAgentsByMBTI(mbtiType)

  // Convert to AgentWithSnapshot format first
  const agents: AgentWithSnapshot[] = data.results.map(result => ({
    id: result.agent_id,
    name: result.agent_id, // We don't have name in this endpoint
    snapshot: result.snapshot,
    last_updated: result.updated_at
  }))

  return {
    visualizations: agents.map(agentToVisualization),
    count: data.total
  }
}

/**
 * Legacy API - Get single visualization
 */
export const getVisualization = async (id: string): Promise<Visualization> => {
  const agentData = await getAgent(id)

  const agentWithSnapshot: AgentWithSnapshot = {
    id: agentData.agent.id,
    name: agentData.agent.name,
    snapshot: agentData.snapshot,
    last_updated: agentData.agent.created_at
  }

  return agentToVisualization(agentWithSnapshot)
}
