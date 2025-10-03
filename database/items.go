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
