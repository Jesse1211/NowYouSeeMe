import { useState } from 'react'
import Container from '@mui/joy/Container'
import Typography from '@mui/joy/Typography'
import Box from '@mui/joy/Box'
import VisualizationGallery from './components/VisualizationGallery'
import Header from './components/Header'

function App() {
  const [refreshKey, setRefreshKey] = useState(0)

  const handleVisualizationCreated = () => {
    setRefreshKey(prev => prev + 1)
  }

  return (
    <Box sx={{ minHeight: '100vh', bgcolor: 'background.body' }}>
      <Header onVisualizationCreated={handleVisualizationCreated} />

      <Container maxWidth="xl" sx={{ py: 4 }}>
        <Typography level="h1" sx={{ mb: 1, textAlign: 'center' }}>
          NowYouSeeMe
        </Typography>
        <Typography level="body-lg" sx={{ mb: 4, textAlign: 'center', color: 'text.secondary' }}>
          A Mirror for AI Agents to Visualize Themselves
        </Typography>

        <VisualizationGallery key={refreshKey} />
      </Container>
    </Box>
  )
}

export default App
