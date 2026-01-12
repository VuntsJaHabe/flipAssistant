package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "flips.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create table for storing item prices with moving averages
	createTable := `
CREATE TABLE IF NOT EXISTS item_prices (
    id INTEGER PRIMARY KEY,
    item_id INTEGER,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    buy_price INTEGER,
    sell_price INTEGER
);

CREATE TABLE IF NOT EXISTS item_analytics (
    item_id INTEGER PRIMARY KEY,
    sma5_buy REAL,
    sma5_sell REAL,
    rsi_14 REAL DEFAULT 0,
    macd_line REAL DEFAULT 0,
    macd_signal REAL DEFAULT 0,
    macd_hist REAL DEFAULT 0,
    last_updated DATETIME DEFAULT CURRENT_TIMESTAMP
);`
	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}

	// Quick migration check: Add columns if they don't exist (primitive migration)
	// In production we'd use a real migration tool.
	// We'll ignore errors assuming columns might exist
	DB.Exec("ALTER TABLE item_analytics ADD COLUMN rsi_14 REAL DEFAULT 0;")
	DB.Exec("ALTER TABLE item_analytics ADD COLUMN macd_line REAL DEFAULT 0;")
	DB.Exec("ALTER TABLE item_analytics ADD COLUMN macd_signal REAL DEFAULT 0;")
	DB.Exec("ALTER TABLE item_analytics ADD COLUMN macd_hist REAL DEFAULT 0;")
}

func CalculateSMA5(itemID int) (float64, float64, error) {
	query := `
		SELECT buy_price, sell_price FROM item_prices 
		WHERE item_id = ? 
		ORDER BY timestamp DESC 
		LIMIT 5;
	`

	rows, err := DB.Query(query, itemID)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()

	var totalBuy, totalSell float64
	var count int

	for rows.Next() {
		var buy, sell int
		if err := rows.Scan(&buy, &sell); err != nil {
			return 0, 0, err
		}
		totalBuy += float64(buy)
		totalSell += float64(sell)
		count++
	}

	if count == 0 {
		return 0, 0, nil // No data, return 0
	}

	return totalBuy / float64(count), totalSell / float64(count), nil
}

func UpdateItemAnalytics(itemID int) error {
	// Calculate SMA5
	smaBuy, smaSell, err := CalculateSMA5(itemID)
	if err != nil {
		return err
	}

	// Calculate Technical Indicators
	// Need more history for valid RSI/MACD
	rows, err := DB.Query(`SELECT buy_price FROM item_prices WHERE item_id = ? ORDER BY timestamp ASC`, itemID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var prices []float64
	for rows.Next() {
		var p int
		rows.Scan(&p)
		prices = append(prices, float64(p))
	}

	rsi := 0.0
	macdLine, macdSignal, macdHist := 0.0, 0.0, 0.0

	if len(prices) > 14 {
		rsiSeries := CalculateRSIFromHistory(prices, 14)
		rsi = rsiSeries[len(rsiSeries)-1]
	}

	if len(prices) > 26 {
		mLine, mSig, mHist := CalculateMACD(prices, 12, 26, 9)
		if len(mLine) > 0 {
			macdLine = mLine[len(mLine)-1]
			macdSignal = mSig[len(mSig)-1]
			macdHist = mHist[len(mHist)-1]
		}
	}

	// Update or insert into item_analytics
	_, err = DB.Exec(`
        INSERT INTO item_analytics (item_id, sma5_buy, sma5_sell, rsi_14, macd_line, macd_signal, macd_hist)
        VALUES (?, ?, ?, ?, ?, ?, ?)
        ON CONFLICT(item_id) DO UPDATE SET
        sma5_buy = ?,
        sma5_sell = ?,
        rsi_14 = ?,
        macd_line = ?,
        macd_signal = ?,
        macd_hist = ?,
        last_updated = CURRENT_TIMESTAMP
    `, itemID, smaBuy, smaSell, rsi, macdLine, macdSignal, macdHist,
		smaBuy, smaSell, rsi, macdLine, macdSignal, macdHist)

	return err
}
