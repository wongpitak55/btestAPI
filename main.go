package main

import (
	"fmt" // To print to the console
	"os"  // To get environment variables

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
		// Print the incoming request to the console
		fmt.Println("Received a POST request to /api")

		var requestBody map[string]interface{}
		if err := c.BindJSON(&requestBody); err != nil {
			// Log the error
			fmt.Println("Error: Invalid JSON payload")
			c.JSON(400, gin.H{"error": "Invalid JSON payload"})
			return
		}

		// Extract "mobile" from the request body
		mobile, exists := requestBody["mobile"].(string)
		if !exists {
			// Log missing "mobile" field
			fmt.Println("Error: Missing 'mobile' field in request payload")
			c.JSON(400, gin.H{"error": "Missing 'mobile' field"})
			return
		}

		// Print the extracted "mobile" field
		fmt.Printf("Received 'mobile': %s\n", mobile)

		// Send a JSON response
		c.JSON(200, gin.H{
			"message": "Data received successfully",
			"mobile":  mobile,
		})
	})

	// Print a message when the server starts
	fmt.Printf("Server is starting on port %s...\n", port)

	// Start the server on the specified port
	router.Run(":" + port)
}
