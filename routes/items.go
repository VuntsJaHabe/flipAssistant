package routes

import (
	"flipAssistant/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetItemInfo(c *gin.Context) {
	itemIDStr := c.Param("id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	itemName := database.GetItemName(itemID)

	c.JSON(http.StatusOK, gin.H{
		"id":   itemID,
		"name": itemName,
	})
}

func GetAllTrackedItems(c *gin.Context) {
	popularItems := database.GetPopularItems()
	var items []map[string]interface{}

	for _, itemID := range popularItems {
		items = append(items, map[string]interface{}{
			"id":   itemID,
			"name": database.GetItemName(itemID),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
		"count": len(items),
	})
}
