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

func GetPopularItems() []int {
	// Expanded list of popular high-volume items across different categories
	return []int{
		// High-value weapons
		4151, 11802, 11804, 11806, 11808, 13576, 13652, 12006,
		// Armour pieces
		11840, 6570, 4087, 4585, 2572, 2577, 2581,
		// Consumables & Food
		2, 384, 386, 390, 392, 361, 373, 3144, 385, 7946,
		// Jewelry
		6585, 1712, 6731, 1704, 1725, 1731, 2572, 6737,
		// Resources & Materials
		1513, 1515, 1517, 1519, 1521, 1623, 1625, 1627, 1629, 1631,
		// Runes (high volume)
		554, 555, 556, 557, 558, 559, 560, 561, 562, 563, 564, 565,
		// Herbs (popular ones)
		199, 201, 203, 205, 207, 209, 211, 213, 215, 217, 219, 2485,
		// Potions (4 dose - most traded)
		2436, 113, 2440, 2442, 2444, 139, 3024, 3026, 3028, 3030,
		// Ores & Bars
		440, 442, 444, 447, 449, 451, 453, 2349, 2351, 2353, 2355, 2357, 2359, 2361, 2363,
		// Seeds (profitable)
		5295, 5296, 5297, 5298, 5299, 5300, 5301, 5302, 5303, 5304,
		// Construction materials
		8778, 8780, 960, 8782, 8784, 8786, 8788,
		// Dragon items
		1187, 1213, 1215, 1231, 1249, 1305, 4087, 11335,
		// Barrows items
		4708, 4710, 4712, 4714, 4716, 4718, 4720, 4722, 4724, 4726, 4728, 4730, 4732, 4734, 4736, 4738, 4745, 4747, 4749, 4751, 4753, 4755, 4757, 4759,
	}
}
