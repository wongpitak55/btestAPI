package main

import (
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Get the port from the environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create a new Gin router
	router := gin.Default()

	// Enable CORS
	router.Use(cors.Default())

	// Define a route
	router.POST("/api", func(c *gin.Context) {
		fmt.Println("Received a POST request to /api")

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

		// Return the data as JSON
		c.JSON(200, gin.H{
			"message": "Data received successfully",
			"data":    data,
		})
	})

	fmt.Printf("Server is starting on port %s...\n", port)
	router.Run(":" + port)
}
