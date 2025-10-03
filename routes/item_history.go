package routes

import (
	"net/http"

	"flipAssistant/database"

	"github.com/gin-gonic/gin"
)

func GetItemHistory(c *gin.Context) {
	itemID := c.Param("id")
	rows, err := database.DB.Query(`
	SELECT timestamp, 
		AVG(buy_price) OVER (ORDER BY timestamp ROWS BETWEEN 6 PRECEDING AND CURRENT ROW) AS avg_buy,
		AVG(sell_price) OVER (ORDER BY timestamp ROWS BETWEEN 6 PRECEDING AND CURRENT ROW) AS avg_sell
	FROM item_prices
	WHERE item_id = ?
	ORDER BY timestamp DESC
`, itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var history []map[string]interface{}
	for rows.Next() {
		var timestamp string
		var buyPrice, sellPrice float64
		if err := rows.Scan(&timestamp, &buyPrice, &sellPrice); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data"})
			return
		}
		history = append(history, map[string]interface{}{
			"timestamp":  timestamp,
			"buy_price":  int(buyPrice),
			"sell_price": int(sellPrice),
		})
	}

	c.JSON(http.StatusOK, gin.H{"history": history})
}
