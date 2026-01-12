# FlipAssistant

### AI-Powered OSRS Grand Exchange Flip Analyzer

FlipAssistant is a tool designed to help Old School RuneScape (OSRS) players identify profitable Grand Exchange flipping opportunities. It automatically tracks item prices, calculates technical indicators, and suggests the best items to flip based on historical data and profit margins.

## Features

- **Real-time Price Tracking**: Automatically fetches and stores Grand Exchange prices from the OSRS Wiki API.
- **Advanced Technical Analysis**:
  - **RSI (Relative Strength Index)**: Identification of overbought/oversold conditions.
  - **MACD (Moving Average Convergence Divergence)**: Trend-following momentum indicator.
  - **SMA (Simple Moving Average)**: 5-period moving averages for calculating reliable margins.
- **Categorized Flip Suggestions**: curated lists of items based on capital requirements and strategy (High Margin, High Volume, etc.).
- **Interactive Visualizations**:
  - Price history charts with overlayed technical indicators.
  - Responsive card-grid layout for browsing items.
- **Search & History**: Robust search functionality with relevance sorting and detailed price history view.
- **Performant Backend**: Go backend with SQLite database for fast data processing and minimal resource usage.

## Architecture

### Backend (Go)
- **Framework**: Gin (HTTP router)
- **Database**: SQLite
- **Analysis**: Custom Go implementations of financial indicators (RSI, MACD)

### Frontend (React + Vite)
- **Framework**: React 18
- **Visualization**: Recharts
- **Styling**: Custom CSS with responsive grid layouts

## Quick Start

### Prerequisites
- Go 1.23+
- Node.js 22+

### Installation

1. **Clone the repository**
   ```bash
   git clone <your-repo-url>
   cd flipAssistant
   ```

2. **Run the development environment**
   ```bash
   make dev
   ```
   This command starts both the backend server (port 8080) and the frontend server (port 5173).

## API Endpoints

- `GET /suggest-flips` - Returns top flip opportunities ranked by profit margin.
- `GET /categorized-flips` - Returns items organized by category (Budget, High Value, etc.).
- `GET /item-history/:id` - Returns comprehensive price history with pre-calculated RSI and MACD values.
- `GET /item-info/:id` - Returns basic item details.
- `GET /search-item` - Search for items by name with fuzzy matching.

## Data Source Compliance

FlipAssistant is compliant with the OSRS Wiki API guidelines:
- Uses bulk `/latest` endpoints to minimize request count.
- Respects rate limits by caching data and performing batch updates.
- Identifies itself with a descriptive User-Agent.

## Disclaimer

This tool is for educational and informational purposes only. The Grand Exchange market is volatile; always perform your own due diligence before making high-value trades.

## Acknowledgments

- **OSRS Wiki** for providing the Real-time Prices API.
- RuneScape® is a trademark of Jagex Ltd. This project is not affiliated with Jagex.
