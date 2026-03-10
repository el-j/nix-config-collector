import { useState } from 'react'
import { Scan, Generate, Preview } from '../wailsjs/go/gui/App'

const DEFAULT_OUTPUT_DIR = '~/.config/nixpkgs'

function App() {
  const [scanResult, setScanResult] = useState(null)
  const [darwinConfig, setDarwinConfig] = useState('')
  const [homeConfig, setHomeConfig] = useState('')
  const [flakeConfig, setFlakeConfig] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)
  const [activeTab, setActiveTab] = useState('scan')
  const [outputDir, setOutputDir] = useState(DEFAULT_OUTPUT_DIR)

  const handleScan = async () => {
    setLoading(true)
    setError(null)
    try {
      const result = await Scan()
      setScanResult(JSON.parse(result))
    } catch (e) {
      setError(String(e))
    } finally {
      setLoading(false)
    }
  }

  const handlePreview = async () => {
    setLoading(true)
    setError(null)
    try {
      const config = await Preview()
      setDarwinConfig(config.DarwinConfig)
      setHomeConfig(config.HomeConfig)
      setFlakeConfig(config.FlakeConfig)
      setActiveTab('darwin')
    } catch (e) {
      setError(String(e))
    } finally {
      setLoading(false)
    }
  }

  const handleGenerate = async () => {
    setLoading(true)
    setError(null)
    try {
      await Generate(outputDir)
    } catch (e) {
      setError(String(e))
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={{ fontFamily: 'system-ui, sans-serif', maxWidth: 1100, margin: '0 auto', padding: 24 }}>
      <h1 style={{ color: '#2d5a8e' }}>🔧 nix-config-collector</h1>
      <p>Scan your macOS system and generate Nix configuration files.</p>

      <div style={{ display: 'flex', gap: 12, marginBottom: 16, alignItems: 'center' }}>
        <button onClick={handleScan} disabled={loading} style={btnStyle('#2d5a8e')}>
          {loading ? '⏳ Scanning…' : '🔍 Scan System'}
        </button>
        <button onClick={handlePreview} disabled={loading || !scanResult} style={btnStyle('#2e7d32')}>
          👁 Preview Config
        </button>
        <button onClick={handleGenerate} disabled={loading || !darwinConfig} style={btnStyle('#c62828')}>
          💾 Write Config Files
        </button>
      </div>

      <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 24 }}>
        <label htmlFor="outputDir" style={{ fontWeight: 'bold', fontSize: 13 }}>Output directory:</label>
        <input
          id="outputDir"
          type="text"
          value={outputDir}
          onChange={e => setOutputDir(e.target.value)}
          style={{ flex: 1, padding: '6px 10px', border: '1px solid #ccc', borderRadius: 4, fontSize: 13 }}
        />
      </div>

      {error && (
        <div style={{ background: '#fff3f3', border: '1px solid #f88', padding: 12, borderRadius: 6, marginBottom: 16 }}>
          ❌ {error}
        </div>
      )}

      {scanResult && (
        <details open style={{ marginBottom: 16 }}>
          <summary style={{ cursor: 'pointer', fontWeight: 'bold' }}>
            ✅ Scan complete — {scanResult.packages?.length ?? 0} packages, {scanResult.services?.length ?? 0} services
          </summary>
          <div style={{ marginTop: 8, display: 'flex', gap: 24 }}>
            <div>
              <strong>Packages</strong>
              <ul style={{ maxHeight: 200, overflowY: 'auto' }}>
                {scanResult.packages?.map(p => (
                  <li key={p.name + p.type}>{p.name} ({p.type}){p.nix_name ? ` → ${p.nix_name}` : ''}</li>
                ))}
              </ul>
            </div>
            <div>
              <strong>Shell</strong>
              <p>{scanResult.shell_config?.shell}</p>
            </div>
          </div>
        </details>
      )}

      {darwinConfig && (
        <div>
          <div style={{ display: 'flex', gap: 0, borderBottom: '2px solid #ddd', marginBottom: 0 }}>
            {['darwin', 'home', 'flake'].map(tab => (
              <button key={tab} onClick={() => setActiveTab(tab)}
                style={{ ...tabStyle, background: activeTab === tab ? '#f5f5f5' : 'white',
                         borderBottom: activeTab === tab ? '2px solid #2d5a8e' : 'none' }}>
                {tab === 'darwin' ? 'darwin-configuration.nix' : tab === 'home' ? 'home.nix' : 'flake.nix'}
              </button>
            ))}
          </div>
          <pre style={{ background: '#1e1e1e', color: '#d4d4d4', padding: 16, borderRadius: '0 0 6px 6px',
                        overflowX: 'auto', fontSize: 13, maxHeight: 500, overflowY: 'auto' }}>
            {activeTab === 'darwin' ? darwinConfig : activeTab === 'home' ? homeConfig : flakeConfig}
          </pre>
        </div>
      )}
    </div>
  )
}

const btnStyle = (color) => ({
  background: color, color: 'white', border: 'none', padding: '10px 20px',
  borderRadius: 6, cursor: 'pointer', fontWeight: 'bold', fontSize: 14,
})

const tabStyle = {
  padding: '8px 16px', border: '1px solid #ddd', borderBottom: 'none',
  cursor: 'pointer', fontSize: 13,
}

export default App
