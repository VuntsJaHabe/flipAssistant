package routes

import (
	"net/http"

	"flipAssistant/database"

	"github.com/gin-gonic/gin"
)

func GetItemHistory(c *gin.Context) {
	itemID := c.Param("id")
	rows, err := database.DB.Query("SELECT timestamp, buy_price, sell_price FROM item_prices WHERE item_id = ? ORDER BY timestamp DESC", itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var history []map[string]interface{}
	for rows.Next() {
		var timestamp string
		var buyPrice, sellPrice int
		if err := rows.Scan(&timestamp, &buyPrice, &sellPrice); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data"})
			return
		}
		history = append(history, map[string]interface{}{
			"timestamp":  timestamp,
			"buy_price":  buyPrice,
			"sell_price": sellPrice,
		})
	}

	c.JSON(http.StatusOK, gin.H{"history": history})
}
