# 🔄 FlipAssistant

AI-Powered OSRS Grand Exchange Flip Analyzer

FlipAssistant is a sophisticated tool designed to help Old School RuneScape (OSRS) players identify profitable Grand Exchange flipping opportunities. It automatically tracks item prices, calculates moving averages, and suggests the best items to flip based on historical data and profit margins.

## ✨ Features

- **Real-time Price Tracking**: Automatically fetches and stores Grand Exchange prices from the OSRS Wiki API
- **Smart Analytics**: Calculates 5-period Simple Moving Averages (SMA5) for buy and sell prices
- **Flip Suggestions**: AI-powered recommendations for the most profitable items to flip
- **Interactive Charts**: Visual price history with trend analysis using Recharts
- **Modern UI**: Clean, responsive React interface with real-time updates
- **RESTful API**: Go backend with Gin framework for fast, reliable data serving

## 🏗️ Architecture

### Backend (Go)
- **Framework**: Gin (HTTP router)
- **Database**: SQLite (lightweight, file-based)
- **Price Source**: OSRS Wiki API
- **Features**: CORS support, automatic price fetching, moving averages

### Frontend (React + Vite)
- **Framework**: React 18 with Vite
- **Styling**: Custom CSS with responsive design
- **Charts**: Recharts for interactive price visualizations
- **HTTP Client**: Axios for API communication

## 🚀 Quick Start

### Prerequisites
- Go 1.23+ 
- Node.js 22+ (use nvm for easy version management)
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone <your-repo-url>
   cd flipAssistant
   ```

2. **Install backend dependencies**
   ```bash
   go mod tidy
   ```

3. **Install frontend dependencies**
   ```bash
   make install
   # or manually: cd frontend && npm install
   ```

### Development

**Easy way using the server script:**
```bash
# Start both servers
./server.sh start

# Stop all servers  
./server.sh stop

# Check server status
./server.sh status

# Restart servers
./server.sh restart
```

**Or using Make commands:**
```bash
# Start both servers (with proper Ctrl+C handling)
make dev

# Stop all servers
make stop

# Check what's running
make status
```

This will start:
- Backend server on `http://localhost:8080`
- Frontend development server on `http://localhost:5173`

**Run servers separately:**
```bash
# Backend only
make backend

# Frontend only  
make frontend
```

## 📊 How It Works

1. **Data Collection**: The system automatically fetches price data for popular OSRS items every 5 minutes
2. **Analytics Processing**: Calculates 5-period Simple Moving Averages for both buy and sell prices
3. **Profit Analysis**: Identifies items with the highest profit margins based on SMA data
4. **Real-time Updates**: Frontend refreshes suggestions every 30 seconds
5. **Historical Tracking**: Stores all price data for trend analysis and charting

## 🎯 Currently Tracked Items

- Dragon Platebody (11840)
- Abyssal Whip (4151)
- Cannonball (2)
- Amulet of Fury (6585)
- Dragon Warhammer (13652)
- Armadyl Godsword (11802)
- Bandos Godsword (11804)
- Saradomin Godsword (11806)
- Zamorak Godsword (11808)
- Dragon Claws (13576)

## 🔧 API Endpoints

- `GET /suggest-flips` - Returns top 10 flip opportunities ranked by profit margin
- `GET /item-history/:id` - Returns price history for a specific item ID
- `GET /item-info/:id` - Returns item name and details for a specific item ID
- `GET /tracked-items` - Returns all currently tracked items with names

## 🤝 OSRS Wiki API Compliance

FlipAssistant is fully compliant with [OSRS Wiki API guidelines](https://oldschool.runescape.wiki/w/RuneScape:Real-time_Prices):

- ✅ **Bulk API Usage**: Single `/latest` call gets all items (not 141 individual requests)
- ✅ **Respectful Rate Limiting**: 1 API call every 10 minutes (6 per hour)
- ✅ **Proper User-Agent**: Descriptive identifier with contact info
- ✅ **Efficient Processing**: 141 items updated from 1 API response
- ✅ **No Individual Item Queries**: Uses bulk endpoint as recommended

## 🗄️ Database Schema

### item_prices
- `id`: Primary key
- `item_id`: OSRS item ID
- `timestamp`: Price recording time
- `buy_price`: High (sell) price from GE
- `sell_price`: Low (buy) price from GE

### item_analytics  
- `item_id`: OSRS item ID (primary key)
- `sma5_buy`: 5-period moving average of buy prices
- `sma5_sell`: 5-period moving average of sell prices
- `last_updated`: Last calculation timestamp

## 🎨 Screenshots

[Add screenshots of your application here]

## 🚧 Future Enhancements

- [ ] Machine learning price prediction models
- [ ] User accounts and portfolio tracking
- [ ] Discord/Telegram bot integration
- [ ] Mobile app (React Native)
- [ ] Advanced analytics (RSI, MACD, Bollinger Bands)
- [ ] Risk assessment scoring
- [ ] Profit calculator with GE tax consideration
- [ ] Push notifications for optimal flip opportunities

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ⚠️ Disclaimer

This tool is for educational and informational purposes only. Always do your own research before making any trading decisions. The Grand Exchange can be volatile and profits are not guaranteed.

## 🙏 Acknowledgments

- [OSRS Wiki](https://oldschool.runescape.wiki/) for providing the free price API
- RuneScape® is a trademark of Jagex Ltd
- This project is not affiliated with or endorsed by Jagex Ltd
