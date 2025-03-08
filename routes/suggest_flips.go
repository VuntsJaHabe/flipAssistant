package routes

import (
	"flipAssistant/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuggestFlips(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT item_id, AVG(buy_price) AS avg_buy, AVG(sell_price) AS avg_sell
		FROM item_prices
		GROUP BY item_id
		HAVING avg_sell - avg_buy > 1000
		ORDER BY (avg_sell - avg_buy) DESC
		LIMIT 10;
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
		return
	}
	defer rows.Close()

	var flips []map[string]interface{}
	for rows.Next() {
		var itemID, avgBuy, avgSell int
		if err := rows.Scan(&itemID, &avgBuy, &avgSell); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data"})
			return
		}
		flips = append(flips, map[string]interface{}{
			"item_id":  itemID,
			"avg_buy":  avgBuy,
			"avg_sell": avgSell,
			"profit":   avgSell - avgBuy,
		})
	}

	c.JSON(http.StatusOK, gin.H{"suggested_flips": flips})
}
