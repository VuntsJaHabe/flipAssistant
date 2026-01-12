package routes

import (
	"net/http"

	"flipAssistant/database"

	"github.com/gin-gonic/gin"
)

func GetItemHistory(c *gin.Context) {
	itemID := c.Param("id")
	rows, err := database.DB.Query(`
	SELECT timestamp, buy_price, sell_price
	FROM item_prices
	WHERE item_id = ?
	ORDER BY timestamp ASC
`, itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var timestamps []string
	var buyPrices []float64
	var sellPrices []float64

	for rows.Next() {
		var timestamp string
		var buy, sell int
		if err := rows.Scan(&timestamp, &buy, &sell); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data"})
			return
		}
		timestamps = append(timestamps, timestamp)
		buyPrices = append(buyPrices, float64(buy))
		sellPrices = append(sellPrices, float64(sell))
	}

	// Calculate indicators based on Buy Price (could offer Sell Price too)
	rsi := database.CalculateRSIFromHistory(buyPrices, 14)
	macdLine, macdSignal, macdHist := database.CalculateMACD(buyPrices, 12, 26, 9)

	history := make([]map[string]interface{}, 0)
	for i := 0; i < len(timestamps); i++ {
		// Safety check for array bounds, though they should match
		r := 0.0
		if i < len(rsi) {
			r = rsi[i]
		}
		ml, ms, mh := 0.0, 0.0, 0.0
		if i < len(macdLine) {
			ml = macdLine[i]
			ms = macdSignal[i]
			mh = macdHist[i]
		}

		history = append(history, map[string]interface{}{
			"timestamp":   timestamps[i],
			"buy_price":   int(buyPrices[i]),
			"sell_price":  int(sellPrices[i]),
			"rsi":         r,
			"macd_line":   ml,
			"macd_signal": ms,
			"macd_hist":   mh,
		})
	}

	// Reverse history for returning newest first if that's what frontend expects,
	// BUT re-reading original code: "ORDER BY timestamp DESC" was used.
	// We fetched ASC for calculation. So we must REVERSE the final slice.
	for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
		history[i], history[j] = history[j], history[i]
	}

	c.JSON(http.StatusOK, gin.H{"history": history})
}
