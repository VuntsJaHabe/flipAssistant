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

          <div className="flips-table">
            <div className={`table-header ${categories[activeTab].name === 'High Margin Items' ? 'five-columns' : ''}`}>
              <div>Item Name</div>
              <div>Buy Price</div>
              <div>Sell Price</div>
              <div>Profit</div>
              {categories[activeTab].name === 'High Margin Items' && <div>Margin %</div>}
            </div>
            
            {categories[activeTab].items.length === 0 ? (
              <div className="no-items">No items found in this category</div>
            ) : (
              categories[activeTab].items.map((flip, index) => (
                <div key={index} className={`table-row ${categories[activeTab].name === 'High Margin Items' ? 'five-columns' : ''}`}>
                  <div className="item-id">{flip.item_name}</div>
                  <div className="buy-price">{formatPrice(Math.round(flip.sma5_buy))}</div>
                  <div className="sell-price">{formatPrice(Math.round(flip.sma5_sell))}</div>
                  <div className="profit">{formatPrice(Math.round(flip.profit))}</div>
                  {categories[activeTab].name === 'High Margin Items' && (
                    <div className="margin-percentage">{formatMarginPercentage(flip.margin_percentage)}</div>
                  )}
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