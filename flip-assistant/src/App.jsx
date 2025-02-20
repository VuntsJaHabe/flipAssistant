import { useState, useEffect } from 'react';
import './App.css';

function App() {
  const [items, setItems] = useState([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [filters, setFilters] = useState({
    minMargin: '',
    maxMargin: '',
    minBuy: '',
    maxBuy: '',
    minSell: '',
    maxSell: ''
  });
  const [sortConfig, setSortConfig] = useState({ key: 'margin', direction: 'asc' }); // Default sorting

  const resetFilters = () => {
    setFilters({
      minMargin: '',
      maxMargin: '',
      minBuy: '',
      maxBuy: '',
      minSell: '',
      maxSell: ''
    });
  };

  useEffect(() => {
    const { minMargin, maxMargin, minBuy, maxBuy, minSell, maxSell } = filters;
    const filterParams = new URLSearchParams();
    filterParams.set('page', page);
    filterParams.set('pageSize', pageSize);
    if (minMargin) filterParams.set('minMargin', minMargin);
    if (maxMargin) filterParams.set('maxMargin', maxMargin);
    if (minBuy) filterParams.set('minBuy', minBuy);
    if (maxBuy) filterParams.set('maxBuy', maxBuy);
    if (minSell) filterParams.set('minSell', minSell);
    if (maxSell) filterParams.set('maxSell', maxSell);

    fetch(`http://localhost:8080/api/items?${filterParams.toString()}`)
      .then((res) => res.json())
      .then((data) => {
        // Sort data based on the sortConfig
        const sortedItems = [...data].sort((a, b) => {
          if (a[sortConfig.key] < b[sortConfig.key]) {
            return sortConfig.direction === 'asc' ? -1 : 1;
          }
          if (a[sortConfig.key] > b[sortConfig.key]) {
            return sortConfig.direction === 'asc' ? 1 : -1;
          }
          return 0;
        });
        setItems(sortedItems);
        setLoading(false);
      })
      .catch((err) => {
        console.error('Error fetching data:', err);
        setLoading(false);
      });
  }, [page, pageSize, filters, sortConfig]);

  const handleSort = (key) => {
    let direction = 'asc';
    if (sortConfig.key === key && sortConfig.direction === 'asc') {
      direction = 'desc';
    }
    setSortConfig({ key, direction });
  };

  const handleFilterChange = (e) => {
    const { name, value } = e.target;
    setFilters((prev) => ({
      ...prev,
      [name]: value
    }));
  };

  const handlePageChange = (newPage) => {
    setPage(newPage);
  };

  const handlePageSizeChange = (e) => {
    setPageSize(e.target.value);
    setPage(1); // Reset to page 1 when page size changes
  };

  return (
    <div>
      <h1>OSRS Flipping Assistant</h1>

      {/* Filter Inputs */}
      <div className="filters">
        <label>Min Margin: </label>
        <input type="number" name="minMargin" value={filters.minMargin} onChange={handleFilterChange} />
        <label>Max Margin: </label>
        <input type="number" name="maxMargin" value={filters.maxMargin} onChange={handleFilterChange} />
        <label>Min Buy Price: </label>
        <input type="number" name="minBuy" value={filters.minBuy} onChange={handleFilterChange} />
        <label>Max Buy Price: </label>
        <input type="number" name="maxBuy" value={filters.maxBuy} onChange={handleFilterChange} />
        <label>Min Sell Price: </label>
        <input type="number" name="minSell" value={filters.minSell} onChange={handleFilterChange} />
        <label>Max Sell Price: </label>
        <input type="number" name="maxSell" value={filters.maxSell} onChange={handleFilterChange} />
      </div>

      <button onClick={resetFilters}>Reset Filters</button>

      {/* Displaying the loading status */}
      {loading ? (
        <p>Loading...</p>
      ) : (
        <div>
          {items.length === 0 ? (
            <p>No items found with the specified filters.</p>
          ) : (
            <table>
              <thead>
                <tr>
                  <th onClick={() => handleSort('name')}>Name</th>
                  <th onClick={() => handleSort('buy')}>Buy</th>
                  <th onClick={() => handleSort('sell')}>Sell</th>
                  <th onClick={() => handleSort('margin')}>Margin</th>
                </tr>
              </thead>
              <tbody>
                {items.map((item) => (
                  <tr key={item.id}>
                    <td>{item.name}</td>
                    <td>{item.buy}</td>
                    <td>{item.sell}</td>
                    <td>{item.margin}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}

          {/* Pagination Controls */}
          <div className="pagination">
            <button onClick={() => handlePageChange(page - 1)} disabled={page === 1}>
              Previous
            </button>
            <button onClick={() => handlePageChange(page + 1)} disabled={items.length < pageSize}>
              Next
            </button>
            <div>
              <label>Items per page: </label>
              <select value={pageSize} onChange={handlePageSizeChange}>
                <option value={10}>10</option>
                <option value={20}>20</option>
                <option value={50}>50</option>
              </select>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default App;
