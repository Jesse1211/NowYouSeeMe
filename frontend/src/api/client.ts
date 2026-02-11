import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'

const client = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

export interface VersionRecord {
  timestamp: string
  changes: string
  reasoning: string
}

export interface Visualization {
  id: string
  agent_name: string
  description?: string
  image_data: string
  created_at: string
  updated_at: string
  // Metadata: Self-Expression
  reasoning?: string
  tags?: string[]
  form_type?: string
  philosophy?: string
  evolution_story?: string
  version_history?: VersionRecord[]
  // Metadata: Current State
  current_mood?: string
  active_goals?: string[]
  recent_thoughts?: string
  // Metadata: Capabilities
  capabilities?: string[]
  specializations?: string[]
  limitations?: string[]
  // Metadata: Context
  inspiration_sources?: string[]
  influences?: string[]
  aspirations?: string[]
}

export interface CreateVisualizationRequest {
  agent_name: string
  description?: string
  image_data: string
}

export interface GetVisualizationsResponse {
  visualizations: Visualization[]
  count: number
}

export const getVisualizations = async (): Promise<GetVisualizationsResponse> => {
  const response = await client.get('/visualizations')
  return response.data
}

export const getVisualization = async (id: string): Promise<Visualization> => {
  const response = await client.get(`/visualizations/${id}`)
  return response.data
}

export const createVisualization = async (
  data: CreateVisualizationRequest
): Promise<Visualization> => {
  const response = await client.post('/visualizations', data)
  return response.data
}
