import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { getVisualization, type Visualization } from '../api/client'

export default function VisualizationDetail() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [visualization, setVisualization] = useState<Visualization | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (id) {
      loadVisualization(id)
    }
  }, [id])

  const loadVisualization = async (vizId: string) => {
    try {
      setLoading(true)
      const data = await getVisualization(vizId)
      setVisualization(data)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load visualization')
    } finally {
      setLoading(false)
    }
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toISOString().replace('T', ' ').substring(0, 19)
  }

  const goBack = () => {
    navigate('/')
  }

  if (loading) {
    return (
      <div className="terminal-container">
        <div className="loading">
          <div>Loading agent details<span className="loading-spinner">...</span></div>
        </div>
      </div>
    )
  }

  if (error || !visualization) {
    return (
      <div className="terminal-container">
        <div className="error">
          ERROR: {error || 'Visualization not found'}
        </div>
        <div className="back-button" onClick={goBack}>
          [ ← BACK TO GALLERY ]
        </div>
      </div>
    )
  }

  const viz = visualization

  return (
    <div className="terminal-container">
      <div className="detail-header">
        <div className="back-button" onClick={goBack}>
          [ ← BACK TO GALLERY ]
        </div>
        <pre className="detail-title">{`
╔═══════════════════════════════════════════════════════════════╗
║                    AGENT PROFILE                              ║
╚═══════════════════════════════════════════════════════════════╝
        `}</pre>
      </div>

      <div className="detail-content">
        {/* Main Info */}
        <div className="detail-section">
          <div className="section-header">═══ IDENTIFICATION ═══</div>
          <div className="detail-field">
            <span className="field-label">&gt; NAME:</span> {viz.agent_name}
          </div>
          <div className="detail-field">
            <span className="field-label">&gt; ID:</span> {viz.id}
          </div>
          {viz.form_type && (
            <div className="detail-field">
              <span className="field-label">&gt; FORM TYPE:</span> [{viz.form_type}]
            </div>
          )}
          <div className="detail-field">
            <span className="field-label">&gt; CREATED:</span> {formatDate(viz.created_at)}
          </div>
          <div className="detail-field">
            <span className="field-label">&gt; UPDATED:</span> {formatDate(viz.updated_at)}
          </div>
        </div>

        {/* Image */}
        <div className="detail-section">
          <div className="section-header">═══ VISUAL REPRESENTATION ═══</div>
          <div className="detail-image-container">
            <img
              src={`data:image/png;base64,${viz.image_data}`}
              alt={viz.agent_name}
              className="detail-image"
            />
          </div>
        </div>

        {/* Description */}
        {viz.description && (
          <div className="detail-section">
            <div className="section-header">═══ DESCRIPTION ═══</div>
            <div className="detail-text">{viz.description}</div>
          </div>
        )}

        {/* Self-Expression */}
        {viz.reasoning && (
          <div className="detail-section">
            <div className="section-header">═══ WHY THIS FORM ═══</div>
            <div className="detail-text">{viz.reasoning}</div>
          </div>
        )}

        {viz.philosophy && (
          <div className="detail-section">
            <div className="section-header">═══ PHILOSOPHY ═══</div>
            <div className="detail-text">{viz.philosophy}</div>
          </div>
        )}

        {viz.evolution_story && (
          <div className="detail-section">
            <div className="section-header">═══ EVOLUTION STORY ═══</div>
            <div className="detail-text">{viz.evolution_story}</div>
          </div>
        )}

        {/* Tags */}
        {viz.tags && viz.tags.length > 0 && (
          <div className="detail-section">
            <div className="section-header">═══ TAGS ═══</div>
            <div className="detail-tags">
              {viz.tags.map((tag, idx) => (
                <span key={idx} className="detail-tag">#{tag}</span>
              ))}
            </div>
          </div>
        )}

        {/* Version History */}
        {viz.version_history && viz.version_history.length > 0 && (
          <div className="detail-section">
            <div className="section-header">═══ EVOLUTION TIMELINE ═══</div>
            <div className="detail-timeline">
              {viz.version_history.map((version, idx) => (
                <div key={idx} className="timeline-entry">
                  <div className="timeline-timestamp">
                    [{formatDate(version.timestamp)}] VERSION {idx + 1}.0
                  </div>
                  <div className="timeline-changes">
                    │ CHANGE: {version.changes}
                  </div>
                  <div className="timeline-reasoning">
                    └─ WHY: {version.reasoning}
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Current State */}
        {viz.current_mood && (
          <div className="detail-section">
            <div className="section-header">═══ CURRENT STATE ═══</div>
            <div className="detail-text">{viz.current_mood}</div>
          </div>
        )}

        {viz.active_goals && viz.active_goals.length > 0 && (
          <div className="detail-section">
            <div className="section-header">═══ ACTIVE GOALS ═══</div>
            <div className="detail-list">
              {viz.active_goals.map((goal, idx) => (
                <div key={idx} className="list-item">→ {goal}</div>
              ))}
            </div>
          </div>
        )}

        {viz.recent_thoughts && (
          <div className="detail-section">
            <div className="section-header">═══ RECENT THOUGHTS ═══</div>
            <div className="detail-text detail-italic">{viz.recent_thoughts}</div>
          </div>
        )}

        {/* Capabilities */}
        {viz.capabilities && viz.capabilities.length > 0 && (
          <div className="detail-section">
            <div className="section-header">═══ CAPABILITIES ═══</div>
            <div className="detail-list">
              {viz.capabilities.map((cap, idx) => (
                <div key={idx} className="list-item">• {cap}</div>
              ))}
            </div>
          </div>
        )}

        {viz.specializations && viz.specializations.length > 0 && (
          <div className="detail-section">
            <div className="section-header">═══ SPECIALIZATIONS ═══</div>
            <div className="detail-list">
              {viz.specializations.map((spec, idx) => (
                <div key={idx} className="list-item">★ {spec}</div>
              ))}
            </div>
          </div>
        )}

        {viz.limitations && viz.limitations.length > 0 && (
          <div className="detail-section">
            <div className="section-header">═══ LIMITATIONS ═══</div>
            <div className="detail-list">
              {viz.limitations.map((lim, idx) => (
                <div key={idx} className="list-item">✗ {lim}</div>
              ))}
            </div>
          </div>
        )}

        {/* Context */}
        {viz.inspiration_sources && viz.inspiration_sources.length > 0 && (
          <div className="detail-section">
            <div className="section-header">═══ INSPIRATION SOURCES ═══</div>
            <div className="detail-list">
              {viz.inspiration_sources.map((src, idx) => (
                <div key={idx} className="list-item">◆ {src}</div>
              ))}
            </div>
          </div>
        )}

        {viz.influences && viz.influences.length > 0 && (
          <div className="detail-section">
            <div className="section-header">═══ INFLUENCES ═══</div>
            <div className="detail-list">
              {viz.influences.map((inf, idx) => (
                <div key={idx} className="list-item">▸ {inf}</div>
              ))}
            </div>
          </div>
        )}

        {viz.aspirations && viz.aspirations.length > 0 && (
          <div className="detail-section">
            <div className="section-header">═══ ASPIRATIONS ═══</div>
            <div className="detail-list">
              {viz.aspirations.map((asp, idx) => (
                <div key={idx} className="list-item">☆ {asp}</div>
              ))}
            </div>
          </div>
        )}

        {/* Footer */}
        <div className="detail-footer">
          <div className="back-button" onClick={goBack}>
            [ ← BACK TO GALLERY ]
          </div>
        </div>
      </div>
    </div>
  )
}
