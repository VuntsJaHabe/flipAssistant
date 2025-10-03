import { useState, useEffect } from 'react'
import './App.css'
import FlipSuggestions from './components/FlipSuggestions'
import ItemHistory from './components/ItemHistory'
import CategorizedFlips from './components/CategorizedFlips'
import Header from './components/Header'

const API_BASE_URL = 'http://localhost:8080'

function App() {
  const [activeTab, setActiveTab] = useState('suggestions')
  const [selectedItemId, setSelectedItemId] = useState(null)

  return (
    <div className="App">
      <Header />
      
      <nav className="nav-tabs">
        <button 
          className={activeTab === 'suggestions' ? 'tab active' : 'tab'}
          onClick={() => setActiveTab('suggestions')}
        >
          Quick Flips
        </button>
        <button 
          className={activeTab === 'categories' ? 'tab active' : 'tab'}
          onClick={() => setActiveTab('categories')}
        >
          Flip Categories
        </button>
        <button 
          className={activeTab === 'history' ? 'tab active' : 'tab'}
          onClick={() => setActiveTab('history')}
        >
          Item History
        </button>
      </nav>

      <main className="main-content">
        {activeTab === 'suggestions' && (
          <FlipSuggestions 
            apiUrl={API_BASE_URL}
            onItemSelect={(itemId) => {
              setSelectedItemId(itemId)
              setActiveTab('history')
            }}
          />
        )}
        
        {activeTab === 'categories' && <CategorizedFlips />}
        
        {activeTab === 'history' && (
          <ItemHistory 
            apiUrl={API_BASE_URL}
            itemId={selectedItemId}
            onItemIdChange={setSelectedItemId}
          />
        )}
      </main>
    </div>
  )
}

export default App
