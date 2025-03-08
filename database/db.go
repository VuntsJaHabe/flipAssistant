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

	// Create table for storing item prices
	createTable := `
	CREATE TABLE IF NOT EXISTS item_prices (
		id INTEGER PRIMARY KEY,
		item_id INTEGER,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		buy_price INTEGER,
		sell_price INTEGER
	);`
	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
}
