package routes

import (
	"flipAssistant/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuggestFlips(c *gin.Context) {
	rows, err := database.DB.Query(`
        SELECT a.item_id, a.sma5_buy, a.sma5_sell, 
               (a.sma5_sell - a.sma5_buy) AS profit_margin
        FROM item_analytics a
        WHERE a.sma5_buy > 0 AND a.sma5_sell > 0
        AND a.sma5_buy < 200000000 -- Exclude items over 200M to avoid low volume 3rd age items
        ORDER BY profit_margin DESC
        LIMIT 10;
    `)
	if err != nil {
		log.Printf("Query error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
		return
	}
	defer rows.Close()

	var flips []map[string]interface{}
	for rows.Next() {
		var itemID int
		var smaBuy, smaSell, profitMargin float64
		if err := rows.Scan(&itemID, &smaBuy, &smaSell, &profitMargin); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data"})
			return
		}
		flips = append(flips, map[string]interface{}{
			"item_id":   itemID,
			"sma5_buy":  smaBuy,
			"sma5_sell": smaSell,
			"profit":    profitMargin,
		})
	}

	c.JSON(http.StatusOK, gin.H{"suggested_flips": flips})
}
