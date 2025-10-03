import { useState, useEffect } from 'react'
import axios from 'axios'
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts'
import './ItemHistory.css'

function ItemHistory({ apiUrl, itemId, onItemIdChange }) {
  const [history, setHistory] = useState([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)
  const [inputItemId, setInputItemId] = useState(itemId || '')
  const [itemName, setItemName] = useState('')

  const fetchItemName = async (id) => {
    try {
      const response = await axios.get(`${apiUrl}/item-info/${id}`)
      setItemName(response.data.name)
    } catch (err) {
      setItemName(`Item ${id}`)
    }
  }

  const fetchHistory = async (id) => {
    if (!id) return
    
    try {
      setLoading(true)
      
      // Fetch both history and item name
      const [historyResponse] = await Promise.all([
        axios.get(`${apiUrl}/item-history/${id}`),
        fetchItemName(id)
      ])
      
      // Process data for the chart
      const processedData = historyResponse.data.history.map((item, index) => ({
        ...item,
        index: index + 1,
        timestamp: new Date(item.timestamp).toLocaleDateString()
      })).reverse() // Reverse to show oldest to newest
      
      setHistory(processedData)
      setError(null)
    } catch (err) {
      setError('Failed to fetch item history. Make sure the item ID exists.')
      console.error('Error fetching history:', err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (itemId) {
      setInputItemId(itemId)
      fetchHistory(itemId)
    }
  }, [itemId, apiUrl])

  const handleSubmit = (e) => {
    e.preventDefault()
    const id = parseInt(inputItemId)
    if (id && id > 0) {
      onItemIdChange(id)
      fetchHistory(id)
    }
  }

  const formatPrice = (price) => {
    if (price >= 1000000) {
      return `${(price / 1000000).toFixed(1)}M`
    } else if (price >= 1000) {
      return `${(price / 1000).toFixed(1)}K`
    }
    return price.toString()
  }

  const currentItemName = itemName || `Item ${itemId}`

  return (
    <div className="item-history">
      <div className="history-header">
        <h2>📈 Item Price History</h2>
        
        <form onSubmit={handleSubmit} className="item-input-form">
          <input
            type="number"
            placeholder="Enter Item ID"
            value={inputItemId}
            onChange={(e) => setInputItemId(e.target.value)}
            className="item-input"
          />
          <button type="submit" className="fetch-btn">Fetch History</button>
        </form>
      </div>

      {itemId && (
        <div className="current-item">
          <h3>Current Item: {currentItemName}</h3>
          <p>ID: {itemId}</p>
        </div>
      )}

      {loading && (
        <div className="loading">
          <div className="spinner"></div>
          <p>Loading item history...</p>
        </div>
      )}

      {error && (
        <div className="error">
          <p>{error}</p>
        </div>
      )}

      {!loading && !error && history.length === 0 && itemId && (
        <div className="no-data">
          <p>No price history available for this item.</p>
        </div>
      )}

      {!loading && !error && history.length > 0 && (
        <div className="chart-container">
          <ResponsiveContainer width="100%" height={500}>
            <LineChart data={history}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis 
                dataKey="timestamp" 
                angle={-45}
                textAnchor="end"
                height={80}
              />
              <YAxis tickFormatter={formatPrice} />
              <Tooltip 
                formatter={(value) => [formatPrice(value), '']}
                labelFormatter={(label) => `Date: ${label}`}
              />
              <Legend />
              <Line 
                type="monotone" 
                dataKey="buy_price" 
                stroke="#e53e3e" 
                name="Buy Price"
                strokeWidth={2}
              />
              <Line 
                type="monotone" 
                dataKey="sell_price" 
                stroke="#38a169" 
                name="Sell Price"
                strokeWidth={2}
              />
            </LineChart>
          </ResponsiveContainer>

          <div className="history-stats">
            <h4>Recent Stats</h4>
            <div className="stats-grid">
              {history.slice(-5).map((item, index) => (
                <div key={index} className="stat-item">
                  <div className="stat-date">{item.timestamp}</div>
                  <div className="stat-prices">
                    <span className="buy-stat">Buy: {formatPrice(item.buy_price)}</span>
                    <span className="sell-stat">Sell: {formatPrice(item.sell_price)}</span>
                    <span className="profit-stat">
                      Profit: {formatPrice(item.sell_price - item.buy_price)}
                    </span>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {!itemId && !loading && (
        <div className="no-item">
          <p>👆 Enter an item ID above to view its price history, or select an item from the flip suggestions.</p>
        </div>
      )}
    </div>
  )
}

export default ItemHistory