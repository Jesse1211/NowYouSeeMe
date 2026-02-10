import VisualizationGallery from './components/VisualizationGallery'
import './terminal.css'

function App() {
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

export default App
