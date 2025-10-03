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
    last_updated DATETIME DEFAULT CURRENT_TIMESTAMP
);`
	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
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

	// Update or insert into item_analytics
	_, err = DB.Exec(`
        INSERT INTO item_analytics (item_id, sma5_buy, sma5_sell)
        VALUES (?, ?, ?)
        ON CONFLICT(item_id) DO UPDATE SET
        sma5_buy = ?,
        sma5_sell = ?,
        last_updated = CURRENT_TIMESTAMP
    `, itemID, smaBuy, smaSell, smaBuy, smaSell)

	return err
}
