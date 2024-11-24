package main

import (
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var dataStore [][]interface{} // Global variable to store received data

func main() {
	// Get the port from the environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create a new Gin router
	router := gin.Default()

	// Enable CORS for frontend requests
	router.Use(cors.Default())

	// Route to receive data from Postman or other clients
	router.POST("/api", func(c *gin.Context) {
		var requestBody map[string]interface{}
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON payload"})
			return
		}

		// Extract "data" as a list of arrays
		data, exists := requestBody["data"].([]interface{})
		if !exists {
			c.JSON(400, gin.H{"error": "Missing or invalid 'data' field"})
			return
		}

		// Convert data to [][]interface{} for storage
		var formattedData [][]interface{}
		for _, item := range data {
			row, ok := item.([]interface{})
			if ok {
				formattedData = append(formattedData, row)
			}
		}

		// Store the received data
		dataStore = formattedData
		fmt.Println("Data stored successfully:", dataStore)

		// Respond to the client
		c.JSON(200, gin.H{"message": "Data stored successfully"})
	})

	// Route to serve the stored data for the frontend
	router.GET("/data", func(c *gin.Context) {
		c.JSON(200, gin.H{"data": dataStore})
	})

	// Start the server
	fmt.Printf("Server is running on port %s...\n", port)
	router.Run(":" + port)
}
