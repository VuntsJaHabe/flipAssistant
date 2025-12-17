package main

import (
	"flipAssistant/database"
	"flipAssistant/routes"
	"flipAssistant/scripts"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database
	database.InitDB()

	// Load items data from JSON
	if err := database.LoadItemsData(); err != nil {
		log.Printf("Warning: Could not load items data: %v", err)
	}

	// Get all tradeable items to track (comprehensive coverage)
	tradeableItems := database.GetAllTradeableItems()
	log.Printf("Tracking %d tradeable items for flip opportunities", len(tradeableItems))

	// Fetch prices for all items in a single API call (API-friendly)
	go func() {
		for {
			log.Println("Fetching prices for all tracked items...")
			scripts.FetchAndStorePricesForAllItems(tradeableItems)

			// Wait 10 minutes between batch updates to be respectful to the API
			// This means we update all 141 items with just 1 API call every 10 minutes
			time.Sleep(10 * time.Minute)
		}
	}()

	// Create a new Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(cors.Default())

	// Define API routes
	r.GET("/item-history/:id", routes.GetItemHistory)
	r.GET("/suggest-flips", routes.SuggestFlips)
	r.GET("/categorized-flips", routes.GetCategorizedFlips)
	r.GET("/item-info/:id", routes.GetItemInfo)
	r.GET("/tracked-items", routes.GetAllTrackedItems)
	r.GET("/search-item", routes.SearchItemByName)

	// Start server
	log.Println("Server running on :8080")
	r.Run(":8080")
}
