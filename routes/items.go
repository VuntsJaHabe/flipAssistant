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

// SearchItemByName searches for an item by name and returns its ID
func SearchItemByName(c *gin.Context) {
	itemName := c.Query("name")
	if itemName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Item name is required"})
		return
	}

	itemID := database.GetItemIDByName(itemName)
	if itemID == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":   itemID,
		"name": database.GetItemName(itemID),
	})
}
