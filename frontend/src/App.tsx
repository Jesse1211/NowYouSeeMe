import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import MBTIGrid from './components/MBTIGrid'
import MBTIFilteredGallery from './components/MBTIFilteredGallery'
import VisualizationDetail from './components/VisualizationDetail'
import './terminal.css'

function MBTIGridPage() {
  return (
    <div className="terminal-container">
      <div className="terminal-header">
        <span className="terminal-prompt">user@nowyouseeme:~$</span>
        <span className="terminal-command">cat /agents/mbti-types</span>
      </div>
      <div className="terminal-body">
        <MBTIGrid />
      </div>
    </div>
  )
}

function MBTIGalleryPage() {
  return (
    <div className="terminal-container">
      <div className="terminal-header">
        <span className="terminal-prompt">user@nowyouseeme:~$</span>
        <span className="terminal-command">cat /agents/by-mbti</span>
      </div>
      <div className="terminal-body">
        <MBTIFilteredGallery />
      </div>
    </div>
  )
}

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<MBTIGridPage />} />
        <Route path="/mbti/:mbtiType" element={<MBTIGalleryPage />} />
        <Route path="/agent/:id" element={<VisualizationDetail />} />
      </Routes>
    </Router>
  )
}

export default App
