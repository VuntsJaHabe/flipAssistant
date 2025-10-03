package scripts

import (
	"encoding/json"
	"flipAssistant/database"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

// OSRSItem represents an item price data structure
type OSRSItem struct {
	High     int `json:"high"`     // Sell price
	HighTime int `json:"highTime"` // Timestamp for high price
	Low      int `json:"low"`      // Buy price
	LowTime  int `json:"lowTime"`  // Timestamp for low price
}

// FetchAndStorePricesForAllItems fetches prices for all tracked items in a single API call
func FetchAndStorePricesForAllItems(itemIDs []int) {
	client := resty.New()

	// Set proper User-Agent to be respectful to the API (as required by OSRS Wiki)
	client.SetHeader("User-Agent", "FlipAssistant/1.0 - OSRS GE Flip Analysis Tool - Contact: github.com/VuntsJaHabe/flipAssistant")

	resp, err := client.R().Get("https://prices.runescape.wiki/api/v1/osrs/latest")
	if err != nil {
		log.Println("Failed to fetch OSRS GE prices:", err)
		return
	}

	var response struct {
		Data map[string]OSRSItem `json:"data"`
	}

	// Parse JSON response
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		log.Println("Error parsing API response:", err)
		return
	}

	successCount := 0
	notFoundCount := 0

	// Process all our tracked items from the single API response
	for _, itemID := range itemIDs {
		itemKey := fmt.Sprintf("%d", itemID)
		if itemData, exists := response.Data[itemKey]; exists {
			// Insert new price data
			_, err := database.DB.Exec(`
				INSERT INTO item_prices (item_id, buy_price, sell_price) 
				VALUES (?, ?, ?)
			`, itemID, itemData.Low, itemData.High)
			if err != nil {
				log.Printf("Error inserting data for item %d: %v", itemID, err)
				continue
			}

			// Update analytics
			if err := database.UpdateItemAnalytics(itemID); err != nil {
				log.Printf("Error updating analytics for item %d: %v", itemID, err)
			}
			successCount++
		} else {
			notFoundCount++
		}
	}

	log.Printf("Price update complete: %d items updated, %d items not found in API", successCount, notFoundCount)
}

// Fetch5MinuteAverages fetches 5-minute price averages for better trend analysis
func Fetch5MinuteAverages(itemIDs []int) {
	client := resty.New()
	client.SetHeader("User-Agent", "FlipAssistant/1.0 - OSRS GE Flip Analysis Tool - Contact: github.com/user/flipAssistant")

	resp, err := client.R().Get("https://prices.runescape.wiki/api/v1/osrs/5m")
	if err != nil {
		log.Println("Failed to fetch 5-minute averages:", err)
		return
	}

	var response struct {
		Data map[string]struct {
			AvgHighPrice    int `json:"avgHighPrice"`
			AvgLowPrice     int `json:"avgLowPrice"`
			HighPriceVolume int `json:"highPriceVolume"`
			LowPriceVolume  int `json:"lowPriceVolume"`
		} `json:"data"`
		Timestamp int64 `json:"timestamp"`
	}

	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		log.Println("Error parsing 5m API response:", err)
		return
	}

	// This could be used for volume-based flip suggestions in the future
	log.Printf("Fetched 5-minute averages with %d items at timestamp %d", len(response.Data), response.Timestamp)
}

// Legacy function for single item (now deprecated, but kept for compatibility)
func FetchAndStorePrices(itemID int) {
	FetchAndStorePricesForAllItems([]int{itemID})
}
