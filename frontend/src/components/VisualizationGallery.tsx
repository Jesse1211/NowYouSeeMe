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
      setVisualizations(data.visualizations || [])
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
            </div>

            <img
              src={`data:image/png;base64,${viz.image_data}`}
              alt={viz.agent_name}
              className="agent-image"
            />

            {viz.description && (
              <div className="agent-description">
                {viz.description}
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
