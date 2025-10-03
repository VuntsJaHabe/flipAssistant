package routes

import (
	"database/sql"
	"flipAssistant/database"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// FlipCategory represents different types of flip opportunities
type FlipCategory struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Items       []map[string]interface{} `json:"items"`
	Count       int                      `json:"count"`
}

// GetCategorizedFlips returns flip suggestions organized by categories
func GetCategorizedFlips(c *gin.Context) {
	categories := []FlipCategory{
		{
			Name:        "High Value Items",
			Description: "Items worth 1M+ GP - High profit potential but requires significant capital",
			Items:       getFlipsByValueRange(1000000, 999999999),
			Count:       0,
		},
		{
			Name:        "Mid Value Items",
			Description: "Items worth 100K-1M GP - Good balance of profit and accessibility",
			Items:       getFlipsByValueRange(100000, 999999),
			Count:       0,
		},
		{
			Name:        "Budget Items",
			Description: "Items worth less than 100K GP - Low capital required, great for beginners",
			Items:       getFlipsByValueRange(0, 99999),
			Count:       0,
		},
		{
			Name:        "High Margin Items",
			Description: "Items with profit margin >5% of item value - Percentage-based profits",
			Items:       getFlipsByMarginPercentage(5.0),
			Count:       0,
		},
		{
			Name:        "High Volume Potential",
			Description: "Items with high GE buy limits - Suitable for bulk trading",
			Items:       getFlipsByBuyLimit(),
			Count:       0,
		},
		{
			Name:        "Quick Flips",
			Description: "Items with consistent small margins - Fast turnover opportunities",
			Items:       getFlipsByConsistency(),
			Count:       0,
		},
	}

	// Set count for each category
	for i := range categories {
		categories[i].Count = len(categories[i].Items)
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
		"timestamp":  time.Now().Unix(),
	})
}

// getFlipsByValueRange returns flips within a specific price range
func getFlipsByValueRange(minPrice, maxPrice int) []map[string]interface{} {
	rows, err := database.DB.Query(`
		SELECT ia.item_id, ia.sma5_buy, ia.sma5_sell, (ia.sma5_sell - ia.sma5_buy) as profit_margin
		FROM item_analytics ia
		WHERE ia.sma5_buy BETWEEN ? AND ?
		AND ia.sma5_buy > 0 AND ia.sma5_sell > 0
		AND (ia.sma5_sell - ia.sma5_buy) > 0
		ORDER BY profit_margin DESC
		LIMIT 10
	`, minPrice, maxPrice)

	if err != nil {
		return []map[string]interface{}{}
	}
	defer rows.Close()

	return processFlipRows(rows)
}

// getFlipsByMarginPercentage returns flips with high percentage margins
func getFlipsByMarginPercentage(minPercentage float64) []map[string]interface{} {
	rows, err := database.DB.Query(`
		SELECT ia.item_id, ia.sma5_buy, ia.sma5_sell, (ia.sma5_sell - ia.sma5_buy) as profit_margin,
		       ((ia.sma5_sell - ia.sma5_buy) / ia.sma5_buy * 100) as margin_percentage
		FROM item_analytics ia
		WHERE ia.sma5_buy > 0 AND ia.sma5_sell > 0
		AND ((ia.sma5_sell - ia.sma5_buy) / ia.sma5_buy * 100) >= ?
		AND (ia.sma5_sell - ia.sma5_buy) > 0
		ORDER BY margin_percentage DESC
		LIMIT 10
	`, minPercentage)

	if err != nil {
		return []map[string]interface{}{}
	}
	defer rows.Close()

	return processFlipRowsWithPercentage(rows)
}

// getFlipsByBuyLimit returns items with high buy limits (approximated by item category)
func getFlipsByBuyLimit() []map[string]interface{} {
	// Items that typically have high buy limits in OSRS
	highLimitItems := []int{
		// Consumables (high limits)
		2, 560, 561, 562, 563, 564, 565, // Cannonballs and runes
		384, 386, 373, 361, // Food items
		199, 201, 203, 205, 207, 209, 211, 213, 215, 217, 219, // Herbs
		440, 442, 444, 447, 449, 451, 453, // Ores
		2349, 2351, 2353, 2355, 2357, 2359, 2361, 2363, // Bars
	}

	if len(highLimitItems) == 0 {
		return []map[string]interface{}{}
	}

	// Create placeholders for the IN clause
	placeholders := "?"
	for i := 1; i < len(highLimitItems); i++ {
		placeholders += ",?"
	}

	query := fmt.Sprintf(`
		SELECT ia.item_id, ia.sma5_buy, ia.sma5_sell, (ia.sma5_sell - ia.sma5_buy) as profit_margin
		FROM item_analytics ia
		WHERE ia.item_id IN (%s)
		AND ia.sma5_buy > 0 AND ia.sma5_sell > 0
		AND (ia.sma5_sell - ia.sma5_buy) > 0
		ORDER BY profit_margin DESC
		LIMIT 10
	`, placeholders)

	// Convert int slice to interface slice for query args
	args := make([]interface{}, len(highLimitItems))
	for i, v := range highLimitItems {
		args[i] = v
	}

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return []map[string]interface{}{}
	}
	defer rows.Close()

	return processFlipRows(rows)
}

// getFlipsByConsistency returns items with consistent but smaller margins
func getFlipsByConsistency() []map[string]interface{} {
	rows, err := database.DB.Query(`
		SELECT ia.item_id, ia.sma5_buy, ia.sma5_sell, (ia.sma5_sell - ia.sma5_buy) as profit_margin
		FROM item_analytics ia
		WHERE ia.sma5_buy > 0 AND ia.sma5_sell > 0
		AND (ia.sma5_sell - ia.sma5_buy) BETWEEN 100 AND 50000
		AND ia.sma5_buy < 500000
		ORDER BY (ia.sma5_sell - ia.sma5_buy) DESC
		LIMIT 10
	`)

	if err != nil {
		return []map[string]interface{}{}
	}
	defer rows.Close()

	return processFlipRows(rows)
}

// processFlipRows processes SQL rows into flip data
func processFlipRows(rows *sql.Rows) []map[string]interface{} {
	var flips []map[string]interface{}

	for rows.Next() {
		var itemID int
		var smaBuy, smaSell, profitMargin float64
		if err := rows.Scan(&itemID, &smaBuy, &smaSell, &profitMargin); err != nil {
			continue
		}
		flips = append(flips, map[string]interface{}{
			"item_id":   itemID,
			"sma5_buy":  smaBuy,
			"sma5_sell": smaSell,
			"profit":    profitMargin,
		})
	}
	return flips
}

// processFlipRowsWithPercentage processes SQL rows including margin percentage
func processFlipRowsWithPercentage(rows *sql.Rows) []map[string]interface{} {
	var flips []map[string]interface{}

	for rows.Next() {
		var itemID int
		var smaBuy, smaSell, profitMargin, marginPercentage float64
		if err := rows.Scan(&itemID, &smaBuy, &smaSell, &profitMargin, &marginPercentage); err != nil {
			continue
		}
		flips = append(flips, map[string]interface{}{
			"item_id":           itemID,
			"sma5_buy":          smaBuy,
			"sma5_sell":         smaSell,
			"profit":            profitMargin,
			"margin_percentage": marginPercentage,
		})
	}
	return flips
}
