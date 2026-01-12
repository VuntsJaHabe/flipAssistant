import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './CategorizedFlips.css';

const CategorizedFlips = () => {
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [activeTab, setActiveTab] = useState(0);

  useEffect(() => {
    fetchCategorizedFlips();
  }, []);

  const fetchCategorizedFlips = async () => {
    try {
      setLoading(true);
      const response = await axios.get('http://localhost:8080/categorized-flips');
      setCategories(response.data.categories);
    } catch (error) {
      console.error('Error fetching categorized flips:', error);
      setError('Failed to load flip categories');
    } finally {
      setLoading(false);
    }
  };

  const formatPrice = (price) => {
    if (price >= 1000000) {
      return `${(price / 1000000).toFixed(1)}M`;
    } else if (price >= 1000) {
      return `${(price / 1000).toFixed(0)}K`;
    }
    return price.toString();
  };

  const formatMarginPercentage = (percentage) => {
    return percentage ? `${percentage.toFixed(1)}%` : '';
  };

  if (loading) {
    return <div className="loading">Loading flip categories...</div>;
  }

  if (error) {
    return <div className="error">{error}</div>;
  }

  return (
    <div className="categorized-flips">
      <div className="category-tabs">
        {categories.map((category, index) => (
          <button
            key={index}
            className={`tab ${activeTab === index ? 'active' : ''}`}
            onClick={() => setActiveTab(index)}
          >
            {category.name}
            <span className="item-count">({category.count})</span>
          </button>
        ))}
      </div>

      {categories[activeTab] && (
        <div className="category-content">
          <div className="category-header">
            <h3>{categories[activeTab].name}</h3>
            <p className="category-description">{categories[activeTab].description}</p>
          </div>

          <div className="flips-grid">
            {categories[activeTab].items.length === 0 ? (
              <div className="no-items">No items found in this category</div>
            ) : (
              categories[activeTab].items.map((flip, index) => (
                <div key={index} className="flip-card">
                  <div className="flip-card-header">
                    <img
                      src={`https://services.runescape.com/m=itemdb_oldschool/obj_sprite.gif?id=${flip.item_id}`}
                      alt={flip.item_name}
                      className="item-icon"
                      onError={(e) => { e.target.style.display = 'none' }}
                    />
                    <div className="item-name">{flip.item_name}</div>
                  </div>

                  <div className="flip-card-stats">
                    <div className="stat-row">
                      <span className="stat-label">Buy:</span>
                      <span className="stat-value buy-price">{formatPrice(Math.round(flip.sma5_buy))}</span>
                    </div>
                    <div className="stat-row">
                      <span className="stat-label">Sell:</span>
                      <span className="stat-value sell-price">{formatPrice(Math.round(flip.sma5_sell))}</span>
                    </div>
                    <div className="stat-row highlight">
                      <span className="stat-label">Profit:</span>
                      <span className="stat-value profit">{formatPrice(Math.round(flip.profit))}</span>
                    </div>
                    {categories[activeTab].name === 'High Margin Items' && (
                      <div className="stat-row">
                        <span className="stat-label">Margin:</span>
                        <span className="stat-value margin">{formatMarginPercentage(flip.margin_percentage)}</span>
                      </div>
                    )}
                  </div>

                  <div className="flip-card-actions">
                    <button className="view-btn">View History</button>
                  </div>
                </div>
              ))
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default CategorizedFlips;