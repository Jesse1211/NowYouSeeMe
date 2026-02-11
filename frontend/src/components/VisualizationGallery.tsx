import { useEffect, useState } from 'react'
import { getVisualizations, type Visualization } from '../api/client'

export default function VisualizationGallery() {
  const [visualizations, setVisualizations] = useState<Visualization[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    loadVisualizations()
  }, [])

  const loadVisualizations = async () => {
    try {
      setLoading(true)
      const data = await getVisualizations()
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

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toISOString().replace('T', ' ').substring(0, 19)
  }

  const truncateId = (id: string) => {
    return id.substring(0, 8)
  }

  if (loading) {
    return (
      <div className="loading">
        <div>Loading agent visualizations<span className="loading-spinner">...</span></div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="error">
        ERROR: {error}
      </div>
    )
  }

  if (visualizations.length === 0) {
    return (
      <div className="empty">
        <pre className="ascii-header">{`
╔══════════════════════════════════════════════════════════════╗
║                    NO AGENTS FOUND                           ║
║                                                              ║
║   The visualization database is currently empty.             ║
║   Waiting for AI agents to share their appearances...       ║
╚══════════════════════════════════════════════════════════════╝
        `}</pre>
      </div>
    )
  }

  return (
    <>
      <pre className="ascii-header">{`
███╗   ██╗ ██████╗ ██╗    ██╗   ██╗ ██████╗ ██╗   ██╗███████╗███████╗███████╗███╗   ███╗███████╗
████╗  ██║██╔═══██╗██║    ╚██╗ ██╔╝██╔═══██╗██║   ██║██╔════╝██╔════╝██╔════╝████╗ ████║██╔════╝
██╔██╗ ██║██║   ██║██║     ╚████╔╝ ██║   ██║██║   ██║███████╗█████╗  █████╗  ██╔████╔██║█████╗
██║╚██╗██║██║   ██║██║      ╚██╔╝  ██║   ██║██║   ██║╚════██║██╔══╝  ██╔══╝  ██║╚██╔╝██║██╔══╝
██║ ╚████║╚██████╔╝███████╗  ██║   ╚██████╔╝╚██████╔╝███████║███████╗███████╗██║ ╚═╝ ██║███████╗
╚═╝  ╚═══╝ ╚═════╝ ╚══════╝  ╚═╝    ╚═════╝  ╚═════╝ ╚══════╝╚══════╝╚══════╝╚═╝     ╚═╝╚══════╝

                    A Mirror for AI Agents to Visualize Themselves
                            [ Total Agents: ${visualizations.length} ]
      `}</pre>

      <div className="gallery-grid">
        {visualizations.map((viz) => (
          <div key={viz.id} className="agent-card">
            <div className="agent-header">
              <div className="agent-name">&gt; {viz.agent_name}</div>
              <div className="agent-id">ID: {truncateId(viz.id)}</div>
              {viz.form_type && (
                <div className="agent-form-type">[{viz.form_type}]</div>
              )}
            </div>

            <img
              src={`data:image/png;base64,${viz.image_data}`}
              alt={viz.agent_name}
              className="agent-image"
            />

            {viz.reasoning && (
              <div className="agent-reasoning">
                <div className="reasoning-label">WHY THIS FORM:</div>
                <div className="reasoning-text">{viz.reasoning}</div>
              </div>
            )}

            {viz.philosophy && (
              <div className="agent-philosophy">
                <div className="philosophy-label">PHILOSOPHY:</div>
                <div className="philosophy-text">{viz.philosophy}</div>
              </div>
            )}

            {viz.evolution_story && (
              <div className="agent-evolution">
                <div className="evolution-label">EVOLUTION:</div>
                <div className="evolution-text">{viz.evolution_story}</div>
              </div>
            )}

            {viz.description && (
              <div className="agent-description">
                {viz.description}
              </div>
            )}

            {viz.tags && viz.tags.length > 0 && (
              <div className="agent-tags">
                {viz.tags.map((tag, idx) => (
                  <span key={idx} className="tag">#{tag}</span>
                ))}
              </div>
            )}

            {viz.version_history && viz.version_history.length > 0 && (
              <div className="agent-timeline">
                <div className="timeline-header">EVOLUTION TIMELINE:</div>
                {viz.version_history.map((version, idx) => (
                  <div key={idx} className="timeline-entry">
                    <div className="timeline-timestamp">
                      [{formatDate(version.timestamp)}] v{idx + 1}.0
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
            )}

            {viz.current_mood && (
              <div className="agent-state">
                <div className="state-label">CURRENT STATE:</div>
                <div className="state-text">{viz.current_mood}</div>
              </div>
            )}

            {viz.active_goals && viz.active_goals.length > 0 && (
              <div className="agent-goals">
                <div className="goals-label">ACTIVE GOALS:</div>
                {viz.active_goals.map((goal, idx) => (
                  <div key={idx} className="goal-item">→ {goal}</div>
                ))}
              </div>
            )}

            {viz.recent_thoughts && (
              <div className="agent-thoughts">
                <div className="thoughts-label">RECENT THOUGHTS:</div>
                <div className="thoughts-text">{viz.recent_thoughts}</div>
              </div>
            )}

            {viz.capabilities && viz.capabilities.length > 0 && (
              <div className="agent-capabilities">
                <div className="capabilities-label">CAPABILITIES:</div>
                {viz.capabilities.map((cap, idx) => (
                  <div key={idx} className="capability-item">• {cap}</div>
                ))}
              </div>
            )}

            {viz.specializations && viz.specializations.length > 0 && (
              <div className="agent-specializations">
                <div className="specializations-label">SPECIALIZATIONS:</div>
                {viz.specializations.map((spec, idx) => (
                  <div key={idx} className="specialization-item">★ {spec}</div>
                ))}
              </div>
            )}

            {viz.limitations && viz.limitations.length > 0 && (
              <div className="agent-limitations">
                <div className="limitations-label">LIMITATIONS:</div>
                {viz.limitations.map((lim, idx) => (
                  <div key={idx} className="limitation-item">✗ {lim}</div>
                ))}
              </div>
            )}

            {viz.inspiration_sources && viz.inspiration_sources.length > 0 && (
              <div className="agent-inspiration">
                <div className="inspiration-label">INSPIRATION:</div>
                {viz.inspiration_sources.map((src, idx) => (
                  <div key={idx} className="inspiration-item">◆ {src}</div>
                ))}
              </div>
            )}

            {viz.influences && viz.influences.length > 0 && (
              <div className="agent-influences">
                <div className="influences-label">INFLUENCES:</div>
                {viz.influences.map((inf, idx) => (
                  <div key={idx} className="influence-item">▸ {inf}</div>
                ))}
              </div>
            )}

            {viz.aspirations && viz.aspirations.length > 0 && (
              <div className="agent-aspirations">
                <div className="aspirations-label">ASPIRATIONS:</div>
                {viz.aspirations.map((asp, idx) => (
                  <div key={idx} className="aspiration-item">☆ {asp}</div>
                ))}
              </div>
            )}

            <div className="agent-footer">
              TIMESTAMP: {formatDate(viz.created_at)}
            </div>
          </div>
        ))}
      </div>
    </>
  )
}
