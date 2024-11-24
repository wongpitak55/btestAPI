package main

import (
	"os" // To get environment variables

	"github.com/gin-gonic/gin"
)

func main() {
	// Get the port from the environment variable
	port := os.Getenv("PORT")
	if port == "" {
		// Default to port 8080 if PORT is not set
		port = "8080"
	}

	// Create a new Gin router
	router := gin.Default()

	// Define a route
	router.POST("/api", func(c *gin.Context) {
		var requestBody map[string]interface{}
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON payload"})
			return
		}
		mobile := requestBody["mobile"].(string)
		c.JSON(200, gin.H{
			"message": "Data received successfully",
			"mobile":  mobile,
		})
	})

	// Start the server on the specified port
	router.Run(":" + port)
}
