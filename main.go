package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type ItemPrice struct {
	High int `json:"high"` // Sell price
	Low  int `json:"low"`  // Buy price
}

type ItemData map[string]ItemPrice

type ItemInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	// Include other fields as needed
}

type APIResponse struct {
	Data map[string]ItemPrice `json:"data"` // <-- Extract the "data" field
}

func fetchOSRSData() (ItemData, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://prices.runescape.wiki/api/v1/osrs/latest", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("User-Agent", "FlipAssistant/1.0 (karmo.kupits@gmail.com)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	fmt.Println("Raw API Response:", string(body)) // Debugging output

	var apiResponse APIResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return apiResponse.Data, nil // Extract the correct "data" field
}

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	data, err := fetchOSRSData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var filteredItems []map[string]interface{}
	for itemID, price := range data {
		margin := price.High - price.Low
		if margin > 1000 {
			// Get the item name from the itemNames map
			name, ok := itemNames[itemID]
			if !ok {
				name = "Unknown" // Fallback name if itemID is not found
			}

			filteredItems = append(filteredItems, map[string]interface{}{
				"id":     itemID,
				"name":   name,
				"buy":    price.Low,
				"sell":   price.High,
				"margin": margin,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredItems)
}

func loadItemData() (map[int]string, error) {
	itemMap := make(map[int]string)

	// Read the JSON file
	data, err := os.ReadFile("items.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read items.json: %v", err)
	}

	// Use map[string]ItemInfo to match JSON structure
	var items map[string]ItemInfo
	err = json.Unmarshal(data, &items)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Convert keys (which are strings) to int
	for key, item := range items {
		var id int
		_, err := fmt.Sscanf(key, "%d", &id) // Convert string key to int
		if err != nil {
			fmt.Printf("Warning: skipping item with key %s due to conversion error: %v\n", key, err)
			continue
		}
		itemMap[id] = item.Name
	}

	return itemMap, nil
}

var itemNames map[string]string

func loadItemNames() error {
	// Read the items.json file
	data, err := ioutil.ReadFile("items.json")
	if err != nil {
		return fmt.Errorf("error reading items.json: %v", err)
	}

	// Parse JSON into the itemNames map
	err = json.Unmarshal(data, &itemNames)
	if err != nil {
		return fmt.Errorf("error parsing items.json: %v", err)
	}

	return nil
}

func main() {

	err := loadItemNames()
	if err != nil {
		log.Fatalf("Error loading item names: %v", err)
	}
	// Set up the HTTP server
	http.HandleFunc("/api/items", itemsHandler)

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
