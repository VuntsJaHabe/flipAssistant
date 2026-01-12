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

  const searchItemByName = async (name) => {
    try {
      const response = await axios.get(`${apiUrl}/search-item`, {
        params: { name }
      })
      return response.data.id
    } catch (err) {
      return null
    }
  }

  const resolveItemInput = async (input) => {
    // Try to parse as number first
    const numInput = parseInt(input)
    if (!isNaN(numInput) && numInput > 0) {
      return numInput
    }

    // Try to search by name
    const itemId = await searchItemByName(input)
    return itemId
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

    if (!inputItemId.trim()) {
      setError('Please enter an item ID or name')
      return
    }

    const parseAndFetch = async () => {
      const resolvedId = await resolveItemInput(inputItemId)

      if (!resolvedId) {
        setError('Item not found. Please check the ID or name and try again.')
        return
      }

      onItemIdChange(resolvedId)
      fetchHistory(resolvedId)
    }

    parseAndFetch()
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
            type="text"
            placeholder="Enter Item ID or Name (e.g., '560' or 'Cannonball')"
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
          <div className="chart-section main-chart">
            <h3>Price Action</h3>
            <ResponsiveContainer width="100%" height={400}>
              <LineChart data={history}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis
                  dataKey="timestamp"
                  angle={-45}
                  textAnchor="end"
                  height={80}
                  tick={{ fontSize: 12 }}
                />
                <YAxis tickFormatter={formatPrice} domain={['auto', 'auto']} />
                <Tooltip
                  formatter={(value) => [formatPrice(value), '']}
                  labelFormatter={(label) => `Date: ${label}`}
                  contentStyle={{ backgroundColor: '#2d3748', border: 'none', color: '#fff' }}
                />
                <Legend verticalAlign="top" />
                <Line
                  type="monotone"
                  dataKey="buy_price"
                  stroke="#ef4444"
                  name="Buy Price"
                  strokeWidth={2}
                  dot={false}
                />
                <Line
                  type="monotone"
                  dataKey="sell_price"
                  stroke="#48bb78"
                  name="Sell Price"
                  strokeWidth={2}
                  dot={false}
                />
              </LineChart>
            </ResponsiveContainer>
          </div>

          <div className="chart-section technical-chart">
            <h3>RSI (14)</h3>
            <ResponsiveContainer width="100%" height={200}>
              <LineChart data={history}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="timestamp" hide />
                <YAxis domain={[0, 100]} ticks={[30, 50, 70]} />
                <Tooltip labelFormatter={(label) => `Date: ${label}`} />
                {/* Reference Lines for Overbought/Oversold */}
                <Line type="monotone" dataKey="rsi" stroke="#8884d8" dot={false} strokeWidth={2} />
              </LineChart>
            </ResponsiveContainer>
          </div>

          <div className="chart-section technical-chart">
            <h3>MACD (12, 26, 9)</h3>
            <ResponsiveContainer width="100%" height={200}>
              <LineChart data={history}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="timestamp" hide />
                <YAxis />
                <Tooltip labelFormatter={(label) => `Date: ${label}`} />
                <Legend verticalAlign="top" height={36} />
                <Line type="monotone" dataKey="macd_line" stroke="#3182ce" name="MACD" dot={false} strokeWidth={2} />
                <Line type="monotone" dataKey="macd_signal" stroke="#ed8936" name="Signal" dot={false} strokeWidth={2} />
              </LineChart>
            </ResponsiveContainer>
          </div>

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
                  <div className="stat-technical">
                    <span className="rsi-stat">RSI: {item.rsi?.toFixed(1)}</span>
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