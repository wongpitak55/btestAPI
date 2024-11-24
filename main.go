package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// A map to store data for multiple clients
var clientData = make(map[string][][]interface{})

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

	// API endpoint for each client
	clients := []string{"worldair", "client2", "client3"} // Add more clients here

	// Create POST and GET endpoints for each client
	for _, client := range clients {
		client := client // Capture range variable

		// POST API for receiving data for a specific client
		router.POST(fmt.Sprintf("/api/%s", client), func(c *gin.Context) {
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

			// Store the data for the specific client
			clientData[client] = formattedData
			fmt.Printf("Data stored successfully for %s: %v\n", client, clientData[client])

			// Respond to the client
			c.JSON(200, gin.H{"message": fmt.Sprintf("Data stored successfully for %s", client)})
		})

		// GET API to serve data for a specific client
		router.GET(fmt.Sprintf("/data/%s", client), func(c *gin.Context) {
			data := clientData[client]
			c.JSON(200, gin.H{"data": data})
		})
	}

	// GET API to serve data for all clients
	router.GET("/data/all", func(c *gin.Context) {
		allData := []map[string]interface{}{}
		for client, data := range clientData {
			allData = append(allData, map[string]interface{}{
				"client": client,
				"data":   data,
			})
		}
		c.JSON(200, gin.H{"clients": allData})
	})

	// Start the self-ping mechanism
	go selfPing()

	// Start the server
	fmt.Printf("Server is running on port %s...\n", port)
	router.Run(":" + port)
}

// selfPing periodically sends a GET request to the server's own endpoint to keep it awake
func selfPing() {
	serverURL := "https://btestapi-am67.onrender.com/"
	for {
		resp, err := http.Get(serverURL)
		if err != nil {
			fmt.Println("Error self-pinging the server:", err)
		} else {
			fmt.Println("Self-ping response:", resp.Status)
			resp.Body.Close()
		}
		time.Sleep(5 * time.Minute) // Adjust interval as needed
	}
}
