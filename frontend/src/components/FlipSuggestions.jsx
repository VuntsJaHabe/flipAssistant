import { useState, useEffect } from 'react'
import axios from 'axios'
import './FlipSuggestions.css'

function FlipSuggestions({ apiUrl, onItemSelect }) {
  const [flips, setFlips] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [itemNames, setItemNames] = useState({})

  const fetchItemNames = async (itemIds) => {
    try {
      const promises = itemIds.map(id => 
        axios.get(`${apiUrl}/item-info/${id}`)
          .then(response => ({ id, name: response.data.name }))
          .catch(() => ({ id, name: `Item ${id}` }))
      )
      const results = await Promise.all(promises)
      const nameMap = {}
      results.forEach(({ id, name }) => {
        nameMap[id] = name
      })
      setItemNames(prev => ({ ...prev, ...nameMap }))
    } catch (err) {
      console.error('Error fetching item names:', err)
    }
  }

  const fetchFlips = async () => {
    try {
      setLoading(true)
      const response = await axios.get(`${apiUrl}/suggest-flips`)
      const flipsData = response.data.suggested_flips || []
      setFlips(flipsData)
      
      // Fetch names for items we don't have yet
      const missingNames = flipsData
        .map(flip => flip.item_id)
        .filter(id => !itemNames[id])
      
      if (missingNames.length > 0) {
        await fetchItemNames(missingNames)
      }
      
      setError(null)
    } catch (err) {
      setError('Failed to fetch flip suggestions. Make sure the backend is running.')
      console.error('Error fetching flips:', err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchFlips()
    // Refresh every 30 seconds
    const interval = setInterval(fetchFlips, 30000)
    return () => clearInterval(interval)
  }, [apiUrl])

  const formatPrice = (price) => {
    if (price >= 1000000) {
      return `${(price / 1000000).toFixed(1)}M`
    } else if (price >= 1000) {
      return `${(price / 1000).toFixed(1)}K`
    }
    return price.toString()
  }

  const getProfitMarginClass = (profit) => {
    if (profit > 100000) return 'profit-high'
    if (profit > 50000) return 'profit-medium'
    return 'profit-low'
  }

  if (loading) {
    return (
      <div className="flip-suggestions">
        <div className="loading">
          <div className="spinner"></div>
          <p>Loading flip suggestions...</p>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flip-suggestions">
        <div className="error">
          <p>{error}</p>
          <button onClick={fetchFlips} className="retry-btn">Retry</button>
        </div>
      </div>
    )
  }

  return (
    <div className="flip-suggestions">
      <div className="suggestions-header">
        <h2>💰 Top Flip Opportunities</h2>
        <button onClick={fetchFlips} className="refresh-btn">🔄 Refresh</button>
      </div>
      
      {flips.length === 0 ? (
        <div className="no-data">
          <p>No flip suggestions available. The system needs more price data.</p>
        </div>
      ) : (
        <div className="flips-grid">
          {flips.map((flip, index) => (
            <div 
              key={flip.item_id} 
              className="flip-card"
              onClick={() => onItemSelect(flip.item_id)}
            >
              <div className="flip-rank">#{index + 1}</div>
              <div className="item-info">
                <h3>{itemNames[flip.item_id] || `Item ${flip.item_id}`}</h3>
                <p className="item-id">ID: {flip.item_id}</p>
              </div>
              
              <div className="price-info">
                <div className="price-row">
                  <span className="label">Buy (SMA5):</span>
                  <span className="price buy-price">{formatPrice(Math.round(flip.sma5_buy))}</span>
                </div>
                <div className="price-row">
                  <span className="label">Sell (SMA5):</span>
                  <span className="price sell-price">{formatPrice(Math.round(flip.sma5_sell))}</span>
                </div>
                <div className={`profit-row ${getProfitMarginClass(flip.profit)}`}>
                  <span className="label">Profit:</span>
                  <span className="profit">{formatPrice(Math.round(flip.profit))}</span>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

export default FlipSuggestions