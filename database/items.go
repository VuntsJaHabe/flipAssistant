package database

import (
	"encoding/json"
	"fmt"
	"os"
)

type Item struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Members     bool   `json:"members"`
	Tradeable   bool   `json:"tradeable"`
	TradeableGE bool   `json:"tradeable_on_ge"`
	Incomplete  bool   `json:"incomplete"`
}

type ItemsData map[string]Item

var itemsCache ItemsData

func LoadItemsData() error {
	if itemsCache != nil {
		return nil // Already loaded
	}

	data, err := os.ReadFile("items.json")
	if err != nil {
		return fmt.Errorf("error reading items.json: %v", err)
	}

	if err := json.Unmarshal(data, &itemsCache); err != nil {
		return fmt.Errorf("error parsing items.json: %v", err)
	}

	return nil
}

func GetItemName(itemID int) string {
	if itemsCache == nil {
		if err := LoadItemsData(); err != nil {
			return fmt.Sprintf("Item %d", itemID)
		}
	}

	if item, exists := itemsCache[fmt.Sprintf("%d", itemID)]; exists {
		return item.Name
	}
	return fmt.Sprintf("Item %d", itemID)
}

func GetAllTradeableItems() []int {
	if itemsCache == nil {
		if err := LoadItemsData(); err != nil {
			return []int{}
		}
	}

	var tradeableItems []int
	for _, item := range itemsCache {
		// Include all tradeable items that can be traded on GE
		if item.Tradeable && item.TradeableGE {
			tradeableItems = append(tradeableItems, item.ID)
		}
	}

	return tradeableItems
}

func GetPopularItems() []int {
	// Now returns all tradeable items instead of hardcoded list
	return GetAllTradeableItems()
}

// GetItemIDByName searches for an item ID by its name (case-insensitive)
// Returns -1 if not found
// GetItemIDByName searches for an item by name and returns its ID
// Returns -1 if not found
func GetItemIDByName(name string) int {
	if itemsCache == nil {
		if err := LoadItemsData(); err != nil {
			return -1
		}
	}

	var candidates []Item

	// Collect all matches
	for _, item := range itemsCache {
		if lowercaseEqual(item.Name, name) {
			candidates = append(candidates, item)
		}
	}

	if len(candidates) == 0 {
		return -1
	}

	// Filter and sort to find the best match
	// Priority 1: Tradeable on GE
	// Priority 2: Not incomplete
	// Priority 3: Exact name match (case sensitive)
	// Priority 4: Lowest ID (usually implies original item vs variants)

	bestCandidate := candidates[0]
	bestScore := calculateItemScore(bestCandidate, name)

	for i := 1; i < len(candidates); i++ {
		score := calculateItemScore(candidates[i], name)
		// We want higher score, or if equal score, lower ID
		if score > bestScore || (score == bestScore && candidates[i].ID < bestCandidate.ID) {
			bestCandidate = candidates[i]
			bestScore = score
		}
	}

	return bestCandidate.ID
}

func calculateItemScore(item Item, searchName string) int {
	score := 0
	if item.TradeableGE {
		score += 100
	}
	if !item.Incomplete {
		score += 50
	}
	if item.Name == searchName {
		score += 25
	}
	return score
}

// lowercaseEqual compares two strings case-insensitively
func lowercaseEqual(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		aLower := toLower(a[i])
		bLower := toLower(b[i])
		if aLower != bLower {
			return false
		}
	}
	return true
}

// toLower converts a byte to lowercase
func toLower(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return b + 32
	}
	return b
}
