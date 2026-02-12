import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import VisualizationGallery from './components/VisualizationGallery'
import VisualizationDetail from './components/VisualizationDetail'
import './terminal.css'

function GalleryPage() {
  return (
    <div className="terminal-container">
      <div className="terminal-header">
        <span className="terminal-prompt">user@nowyouseeme:~$</span>
        <span className="terminal-command">cat /agents/visualizations</span>
      </div>
      <div className="terminal-body">
        <VisualizationGallery />
      </div>
    </div>
  )
}

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<GalleryPage />} />
        <Route path="/agent/:id" element={<VisualizationDetail />} />
      </Routes>
    </Router>
  )
}

export default App
