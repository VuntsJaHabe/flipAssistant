package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type ItemPrice struct {
	High int `json:"high"` // Sell price
	Low  int `json:"low"`  // Buy price
}

type ItemData map[string]ItemPrice

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	// Include other fields as needed
}

type APIResponse struct {
	Data map[string]ItemPrice `json:"data"`
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

	// Extract query parameters for pagination and filtering
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	minMarginStr := r.URL.Query().Get("minMargin")
	maxMarginStr := r.URL.Query().Get("maxMargin")
	minBuyStr := r.URL.Query().Get("minBuy")
	maxBuyStr := r.URL.Query().Get("maxBuy")
	minSellStr := r.URL.Query().Get("minSell")
	maxSellStr := r.URL.Query().Get("maxSell")

	// Default values for pagination
	page := 1
	pageSize := 10

	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			http.Error(w, "Invalid page number", http.StatusBadRequest)
			return
		}
	}

	if pageSizeStr != "" {
		var err error
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil {
			http.Error(w, "Invalid page size", http.StatusBadRequest)
			return
		}
	}

	data, err := fetchOSRSData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Filter the items
	var filteredItems []map[string]interface{}
	for itemIDStr, price := range data {
		margin := price.High - price.Low
		if margin > 0 {
			itemID, err := strconv.Atoi(itemIDStr)
			if err != nil {
				continue
			}

			item, ok := itemNames[itemID]
			if !ok {
				item.Name = "Unknown"
			}

			// Apply additional filters
			if minMarginStr != "" {
				minMargin, err := strconv.Atoi(minMarginStr)
				if err == nil && margin < minMargin {
					continue
				}
			}
			if maxMarginStr != "" {
				maxMargin, err := strconv.Atoi(maxMarginStr)
				if err == nil && margin > maxMargin {
					continue
				}
			}
			if minBuyStr != "" {
				minBuy, err := strconv.Atoi(minBuyStr)
				if err == nil && price.Low < minBuy {
					continue
				}
			}
			if maxBuyStr != "" {
				maxBuy, err := strconv.Atoi(maxBuyStr)
				if err == nil && price.Low > maxBuy {
					continue
				}
			}
			if minSellStr != "" {
				minSell, err := strconv.Atoi(minSellStr)
				if err == nil && price.High < minSell {
					continue
				}
			}
			if maxSellStr != "" {
				maxSell, err := strconv.Atoi(maxSellStr)
				if err == nil && price.High > maxSell {
					continue
				}
			}

			// Add the item to the list
			filteredItems = append(filteredItems, map[string]interface{}{
				"id":     itemID,
				"name":   item.Name,
				"buy":    price.Low,
				"sell":   price.High,
				"margin": margin,
			})
		}
	}

	// Pagination: Slice the items array
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > len(filteredItems) {
		end = len(filteredItems)
	}
	paginatedItems := filteredItems[start:end]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paginatedItems)
}

func loadItemData() (map[int]string, error) {
	itemMap := make(map[int]string)

	// Read the JSON file
	data, err := os.ReadFile("items.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read items.json: %v", err)
	}

	// Use map[string]ItemInfo to match JSON structure
	var items map[string]Item
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

var itemNames map[int]Item

func loadItemNames() error {
	// Read the items.json file
	data, err := ioutil.ReadFile("items.json")
	if err != nil {
		return fmt.Errorf("error reading items.json: %v", err)
	}

	// Parse the JSON data into a map
	var items map[string]Item
	err = json.Unmarshal(data, &items)
	if err != nil {
		return fmt.Errorf("error parsing items.json: %v", err)
	}

	// Convert map keys to integers and store the item names
	itemNames = make(map[int]Item)
	for key, item := range items {
		id, err := strconv.Atoi(key)
		if err != nil {
			return fmt.Errorf("error converting item ID to integer: %v", err)
		}
		itemNames[id] = item
	}

	return nil
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func main() {
	err := loadItemNames()
	if err != nil {
		fmt.Println("Error loading item names:", err)
		return
	}
	// Set up the HTTP server
	http.HandleFunc("/api/items", enableCORS(itemsHandler))

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
