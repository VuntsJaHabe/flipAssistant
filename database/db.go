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
		sell_price INTEGER,
		sma5_buy REAL DEFAULT NULL,
		sma5_sell REAL DEFAULT NULL
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
