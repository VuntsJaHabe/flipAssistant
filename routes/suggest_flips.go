package routes

import (
	"flipAssistant/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuggestFlips(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT item_id, sma5_buy, sma5_sell
		FROM item_prices
		WHERE sma5_buy IS NOT NULL AND sma5_sell IS NOT NULL
		GROUP BY item_id
		HAVING sma5_sell - sma5_buy > 1000
		ORDER BY (sma5_sell - sma5_buy) DESC
		LIMIT 10;
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
		return
	}
	defer rows.Close()

	var flips []map[string]interface{}
	for rows.Next() {
		var itemID int
		var smaBuy, smaSell float64
		if err := rows.Scan(&itemID, &smaBuy, &smaSell); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data"})
			return
		}
		flips = append(flips, map[string]interface{}{
			"item_id":   itemID,
			"sma5_buy":  smaBuy,
			"sma5_sell": smaSell,
			"profit":    smaSell - smaBuy,
		})
	}

	c.JSON(http.StatusOK, gin.H{"suggested_flips": flips})
}
