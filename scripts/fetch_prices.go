package scripts

import (
	"encoding/json"
	"flipAssistant/database"
	"fmt"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
)

// OSRSItem represents an item price data structure
type OSRSItem struct {
	High     int `json:"high"`     // Sell price
	HighTime int `json:"highTime"` // Timestamp for high price
	Low      int `json:"low"`      // Buy price
	LowTime  int `json:"lowTime"`  // Timestamp for low price
}

func FetchAndStorePrices(itemID int) {
	client := resty.New()
	resp, err := client.R().Get("https://prices.runescape.wiki/api/v1/osrs/latest")
	if err != nil {
		log.Println("Failed to fetch OSRS GE prices:", err)
		return
	}

	var response struct {
		Data map[string]map[string]interface{} `json:"data"`
	}

	// Parse JSON response
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		log.Println("Error parsing API response:", err)
		return
	}

	// Check if the item exists in the response
	itemKey := fmt.Sprintf("%d", itemID)
	if itemData, exists := response.Data[itemKey]; exists {
		highPrice := int(itemData["high"].(float64))
		lowPrice := int(itemData["low"].(float64))

		// Insert new price data
		_, err := database.DB.Exec(`
			INSERT INTO item_prices (item_id, buy_price, sell_price) 
			VALUES (?, ?, ?)
		`, itemID, lowPrice, highPrice)
		if err != nil {
			log.Println("Error inserting data:", err)
			return
		}

		// Calculate SMA5 and update database
		smaBuy, smaSell, err := database.CalculateSMA5(itemID)
		if err != nil {
			log.Println("Error calculating SMA5:", err)
			return
		}

		_, err = database.DB.Exec(`
			UPDATE item_prices 
			SET sma5_buy = ?, sma5_sell = ? 
			WHERE item_id = ? 
			ORDER BY timestamp DESC 
			LIMIT 1
		`, smaBuy, smaSell, itemID)
		if err != nil {
			log.Println("Error updating SMA5:", err)
		}
	} else {
		log.Printf("Item %d not found in API response.\n", itemID)
	}
}

func main() {
	for {
		FetchAndStorePrices(11840) // Example item ID
		time.Sleep(10 * time.Minute)
	}
}
