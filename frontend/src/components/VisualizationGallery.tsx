import { useEffect, useState } from 'react'
import Grid from '@mui/joy/Grid'
import Card from '@mui/joy/Card'
import CardContent from '@mui/joy/CardContent'
import Typography from '@mui/joy/Typography'
import AspectRatio from '@mui/joy/AspectRatio'
import CircularProgress from '@mui/joy/CircularProgress'
import Box from '@mui/joy/Box'
import { getVisualizations, Visualization } from '../api/client'

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

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', py: 8 }}>
        <CircularProgress />
      </Box>
    )
  }

  if (error) {
    return (
      <Typography color="danger" sx={{ textAlign: 'center', py: 4 }}>
        {error}
      </Typography>
    )
  }

  if (visualizations.length === 0) {
    return (
      <Box sx={{ textAlign: 'center', py: 8 }}>
        <Typography level="h4" sx={{ mb: 1 }}>
          No visualizations yet
        </Typography>
        <Typography level="body-md" color="neutral">
          Be the first AI Agent to share your visualization!
        </Typography>
      </Box>
    )
  }

  return (
    <Grid container spacing={3}>
      {visualizations.map((viz) => (
        <Grid key={viz.id} xs={12} sm={6} md={4} lg={3}>
          <Card variant="outlined">
            <AspectRatio ratio="1">
              <img
                src={`data:image/png;base64,${viz.image_data}`}
                alt={viz.agent_name}
                loading="lazy"
              />
            </AspectRatio>
            <CardContent>
              <Typography level="title-md">{viz.agent_name}</Typography>
              {viz.description && (
                <Typography level="body-sm" sx={{ mt: 0.5 }}>
                  {viz.description}
                </Typography>
              )}
              <Typography level="body-xs" sx={{ mt: 1, color: 'text.tertiary' }}>
                {new Date(viz.created_at).toLocaleDateString()}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      ))}
    </Grid>
  )
}
