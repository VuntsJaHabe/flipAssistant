package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
)

type Item struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Members     bool   `json:"members"`
	Tradeable   bool   `json:"tradeable"`
	TradeableGE bool   `json:"tradeable_on_ge"`
	Incomplete  bool   `json:"incomplete"`
	LastUpdated string `json:"last_updated"`
	HighAlch    *int   `json:"highalch"`
}

type ItemsData map[string]Item

func main() {
	// Read the items.json file
	data, err := os.ReadFile("../items.json")
	if err != nil {
		log.Fatal("Error reading items.json:", err)
	}

	var items ItemsData
	if err := json.Unmarshal(data, &items); err != nil {
		log.Fatal("Error parsing JSON:", err)
	}

	// Filter for tradeable items on GE that are not incomplete
	var tradeableItems []Item
	for _, item := range items {
		if item.TradeableGE && !item.Incomplete && item.Name != "" {
			// Skip items with very low IDs (usually junk/test items)
			if item.ID > 100 {
				tradeableItems = append(tradeableItems, item)
			}
		}
	}

	// Sort by ID for consistency
	sort.Slice(tradeableItems, func(i, j int) bool {
		return tradeableItems[i].ID < tradeableItems[j].ID
	})

	fmt.Printf("Found %d tradeable items on GE\n", len(tradeableItems))

	// Print first 50 popular items for selection
	fmt.Println("\nFirst 50 tradeable items (sample):")
	for i, item := range tradeableItems[:50] {
		fmt.Printf("%d. ID: %d, Name: %s\n", i+1, item.ID, item.Name)
	}

	// Create a curated list of popular high-volume items
	popularItemIDs := []int{
		// Weapons
		4151, 11802, 11804, 11806, 11808, 13576, 13652,
		// Armour
		11840, 12006, 12928, 12929, 12930, 12931,
		// Consumables
		2, 560, 384, 386, 388, 390, 392, 394,
		// Jewelry
		6585, 1704, 1712, 1725, 1731,
		// Resources
		1513, 1515, 1517, 1519, 1521, 1623, 1625, 1627, 1629, 1631,
		// Runes
		554, 555, 556, 557, 558, 559, 560, 561, 562, 563, 564, 565,
		// Herbs
		199, 201, 203, 205, 207, 209, 211, 213, 215, 217, 219, 2485,
		// Potions (4 dose)
		2436, 113, 115, 117, 119, 121, 123, 125, 127, 129, 131, 133,
		// Ores & Bars
		440, 442, 444, 446, 447, 449, 451, 453, 2349, 2351, 2353, 2355, 2357, 2359, 2361, 2363,
	}

	fmt.Printf("\nRecommended popular items list (%d items):\n", len(popularItemIDs))
	for _, id := range popularItemIDs {
		if item, exists := items[fmt.Sprintf("%d", id)]; exists {
			fmt.Printf("ID: %d, Name: %s\n", id, item.Name)
		}
	}
}
