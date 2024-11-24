package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a new Gin router
	router := gin.Default()

	// Define a route to handle POST requests
	router.POST("/api", func(c *gin.Context) {
		// Parse the incoming JSON payload
		var requestBody map[string]interface{}
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
			return
		}

		// Extract specific fields (if needed)
		mobile := requestBody["mobile"].(string)

		// Log the received data (for debugging)
		c.JSON(http.StatusOK, gin.H{
			"message": "Data received successfully",
			"mobile":  mobile,
		})
	})

	// Start the server on port 8080
	router.Run(":8080")
}
