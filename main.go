package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
		// Example third-party API endpoint
		thirdPartyAPI := "https://jsonplaceholder.typicode.com/users" // Replace with your actual third-party API

		// Make a GET request to the third-party API
		resp, err := http.Get(thirdPartyAPI)
		if err != nil {
			fmt.Println("Error fetching data from third-party API:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data from third-party API"})
			return
		}
		defer resp.Body.Close()

		// Parse the third-party API response
		var data []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			fmt.Println("Error decoding API response:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse third-party API response"})
			return
		}

		// Log the data and send it to the frontend
		fmt.Println("Data received from third-party API:", data)
		c.JSON(http.StatusOK, gin.H{
			"message": "Data fetched successfully",
			"data":    data,
		})
	})

	// Print a message when the server starts
	fmt.Printf("Server is starting on port %s...\n", port)
	router.Run(":" + port)
}
