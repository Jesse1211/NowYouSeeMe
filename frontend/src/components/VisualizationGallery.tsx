import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { getVisualizations, type Visualization } from '../api/client'

export default function VisualizationGallery() {
  const [visualizations, setVisualizations] = useState<Visualization[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const navigate = useNavigate()

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

  const truncateId = (id: string) => {
    return id.substring(0, 8)
  }

  const handleCardClick = (id: string) => {
    navigate(`/agent/${id}`)
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
                            [ Click any agent to view details ]
      `}</pre>

      <div className="gallery-grid">
        {visualizations.map((viz) => (
          <div
            key={viz.id}
            className="agent-card agent-card-clickable"
            onClick={() => handleCardClick(viz.id)}
          >
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

            {viz.capabilities && viz.capabilities.length > 0 && (
              <div className="agent-capabilities">
                <div className="capabilities-label">CAPABILITIES:</div>
                {viz.capabilities.slice(0, 3).map((cap, idx) => (
                  <div key={idx} className="capability-item">• {cap}</div>
                ))}
                {viz.capabilities.length > 3 && (
                  <div className="capability-more">+ {viz.capabilities.length - 3} more...</div>
                )}
              </div>
            )}

            <div className="agent-footer">
              <div className="click-hint">[ CLICK FOR FULL DETAILS ]</div>
            </div>
          </div>
        ))}
      </div>
    </>
  )
}
