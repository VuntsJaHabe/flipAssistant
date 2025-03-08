package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
)

// Struct matching API response
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

	// Debug: Print raw response
	log.Println("API Response:", string(resp.Body()))

	// Define a struct to parse the response
	var response struct {
		Data map[string]map[string]interface{} `json:"data"`
	}

	// Parse JSON response
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		log.Println("Error parsing API response:", err)
		return
	}

	// Check for the specific item in the response
	itemKey := fmt.Sprintf("%d", itemID)
	if itemData, exists := response.Data[itemKey]; exists {
		// Extract price data if the item exists
		highPrice := itemData["high"]
		lowPrice := itemData["low"]
		log.Printf("Item %d - High Price: %v, Low Price: %v", itemID, highPrice, lowPrice)
	} else {
		log.Printf("Item %d not found in response.", itemID)
	}
}

func main() {
	for {
		FetchAndStorePrices(11840)
		time.Sleep(10 * time.Minute)
	}
}
