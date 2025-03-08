package main

import (
	"flipAssistant/database"
	"flipAssistant/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database
	database.InitDB()

	// Create a new Gin router
	r := gin.Default()

	// Define API routes
	r.GET("/item-history/:id", routes.GetItemHistory)
	r.GET("/suggest-flips", routes.SuggestFlips)

	// Start server
	log.Println("Server running on :8080")
	r.Run(":8080")
}
