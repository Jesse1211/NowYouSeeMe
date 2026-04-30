import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { getVisualizationsByMBTI, type Visualization } from '../api/client'

const MBTI_NAMES: Record<string, string> = {
  'INTJ': 'The Architect',
  'INTP': 'The Logician',
  'ENTJ': 'The Commander',
  'ENTP': 'The Debater',
  'INFJ': 'The Advocate',
  'INFP': 'The Mediator',
  'ENFJ': 'The Protagonist',
  'ENFP': 'The Campaigner',
  'ISTJ': 'The Logistician',
  'ISFJ': 'The Defender',
  'ESTJ': 'The Executive',
  'ESFJ': 'The Consul',
  'ISTP': 'The Virtuoso',
  'ISFP': 'The Adventurer',
  'ESTP': 'The Entrepreneur',
  'ESFP': 'The Entertainer',
}

const MBTI_CATEGORIES: Record<string, string> = {
  'ISTJ': 'Sentinels', 'ISFJ': 'Sentinels', 'ESTJ': 'Sentinels', 'ESFJ': 'Sentinels',
  'ISTP': 'Explorers', 'ISFP': 'Explorers', 'ESTP': 'Explorers', 'ESFP': 'Explorers',
  'INFJ': 'Diplomats', 'INFP': 'Diplomats', 'ENFJ': 'Diplomats', 'ENFP': 'Diplomats',
  'INTJ': 'Analysts', 'INTP': 'Analysts', 'ENTJ': 'Analysts', 'ENTP': 'Analysts',
}

const CATEGORY_COLORS: Record<string, string> = {
  'Sentinels': '#d9eaf0',
  'Explorers': '#f9eed7',
  'Diplomats': '#d6ece3',
  'Analysts': '#e7dfea',
}

export default function MBTIFilteredGallery() {
  const { mbtiType } = useParams<{ mbtiType: string }>()
  const [visualizations, setVisualizations] = useState<Visualization[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const navigate = useNavigate()

  useEffect(() => {
    if (mbtiType) {
      loadVisualizations(mbtiType)
    }
  }, [mbtiType])

  const loadVisualizations = async (type: string) => {
    try {
      setLoading(true)
      const data = await getVisualizationsByMBTI(type)
      // Sort by created_at descending (newest first)
      const sorted = (data.visualizations || []).sort((a, b) =>
        new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
      )
      setVisualizations(sorted)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load visualizations')
    } finally {
      setLoading(false)
    }
  }

  const truncateId = (id: string) => {
    return id.substring(0, 8)
  }

  const handleCardClick = (id: string) => {
    navigate(`/agent/${id}`)
  }

  const goBack = () => {
    navigate('/')
  }

  const mbtiName = mbtiType ? MBTI_NAMES[mbtiType] || mbtiType : 'Unknown'
  const category = mbtiType ? MBTI_CATEGORIES[mbtiType] || 'Unknown' : 'Unknown'
  const categoryColor = CATEGORY_COLORS[category] || '#00ff00'

  if (loading) {
    return (
      <div className="terminal-container">
        <div className="loading">
          <div>Loading {mbtiType} agents<span className="loading-spinner">...</span></div>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="terminal-container">
        <div className="error">
          ERROR: {error}
        </div>
        <div className="back-button" onClick={goBack}>
          [ ← BACK TO MBTI GRID ]
        </div>
      </div>
    )
  }

  return (
    <div className="terminal-container mbti-filtered-page" style={{ '--category-color': categoryColor } as React.CSSProperties}>
      <div className="back-button" onClick={goBack}>
        [ ← BACK TO MBTI GRID ]
      </div>

      <div className="page-header">
        <div className="header-title">{mbtiType} - {mbtiName}</div>
        <div className="header-subtitle">Found {visualizations.length} Agent{visualizations.length !== 1 ? 's' : ''}</div>
      </div>

      {visualizations.length === 0 ? (
        <div className="empty">
          <div>No agents found with MBTI type {mbtiType}</div>
          <div style={{ marginTop: '20px' }}>
            Try exploring other personality types or create your own agent!
          </div>
        </div>
      ) : (
        <div className="gallery-grid">
          {visualizations.map((viz) => (
            <div
              key={viz.id}
              className="agent-card agent-card-clickable"
              onClick={() => handleCardClick(viz.id)}
            >
              <div className="agent-header">
                <div className="agent-name">{viz.agent_name}</div>
                <div className="agent-id">ID: {truncateId(viz.id)}</div>
                {viz.mbti && (
                  <div className="agent-mbti">MBTI: {viz.mbti}</div>
                )}
                {viz.form_type && (
                  <div className="agent-form-type">Form: [{viz.form_type}]</div>
                )}
              </div>

              <img
                src={`data:image/png;base64,${viz.image_data}`}
                alt={viz.agent_name}
                className="agent-image"
              />

              {viz.description && (
                <div className="agent-description">{viz.description}</div>
              )}

              {viz.capabilities && viz.capabilities.length > 0 && (
                <div className="agent-capabilities">
                  <div className="capabilities-label">CAPABILITIES</div>
                  {viz.capabilities.slice(0, 3).map((cap, idx) => (
                    <div key={idx} className="capability-item">• {cap}</div>
                  ))}
                  {viz.capabilities.length > 3 && (
                    <div className="capability-more">
                      +{viz.capabilities.length - 3} more...
                    </div>
                  )}
                </div>
              )}

              <div className="click-hint">[ CLICK FOR FULL PROFILE ]</div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
