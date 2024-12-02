package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Maps to store data for multiple clients and bot logs
var clientData = make(map[string][][]interface{})
var botLogData = make(map[string][][]interface{})   // Separate storage for bot logs
var hardDiskData = make(map[string][][]interface{}) // Separate storage for hardDiskData logs

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

	// Bot Log APIs
	botLogEndpoints := []string{"worldair", "log2", "log3"} // Add more bot log categories as needed

	for _, botLog := range botLogEndpoints {
		botLog := botLog // Capture range variable

		// POST API for receiving bot log data
		router.POST(fmt.Sprintf("/api/bot/%s", botLog), func(c *gin.Context) {
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

			// Store the data for the specific bot log category
			botLogData[botLog] = formattedData
			fmt.Printf("Bot log data stored successfully for %s: %v\n", botLog, botLogData[botLog])

			// Respond to the client
			c.JSON(200, gin.H{"message": fmt.Sprintf("Bot log data stored successfully for %s", botLog)})
		})

		// GET API to serve bot log data for a specific category
		router.GET(fmt.Sprintf("/data/bot/%s", botLog), func(c *gin.Context) {
			data := botLogData[botLog]
			c.JSON(200, gin.H{"data": data})
		})
	}

	// GET API to serve all bot log data
	router.GET("/data/bot/all", func(c *gin.Context) {
		allBotLogData := []map[string]interface{}{}
		for botLog, data := range botLogData {
			allBotLogData = append(allBotLogData, map[string]interface{}{
				"botLog": botLog,
				"data":   data,
			})
		}
		c.JSON(200, gin.H{"botLogs": allBotLogData})
	})

	// Hard Disk Used APIs
	hardDiskUsedEndpoints := []string{"worldair", "log2", "log3"} // Add more hard disk categories as needed

	for _, hardDisk := range hardDiskUsedEndpoints {
		hardDisk := hardDisk // Capture range variable

		// POST API for receiving hard disk used data
		router.POST(fmt.Sprintf("/api/harddisk/%s", hardDisk), func(c *gin.Context) {
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

			// Store the data for the specific hard disk category
			hardDiskData[hardDisk] = formattedData
			fmt.Printf("Hard disk data stored successfully for %s: %v\n", hardDisk, hardDiskData[hardDisk])

			// Respond to the client
			c.JSON(200, gin.H{"message": fmt.Sprintf("Hard disk data stored successfully for %s", hardDisk)})
		})

		// GET API to serve hard disk used data for a specific category
		router.GET(fmt.Sprintf("/data/harddisk/%s", hardDisk), func(c *gin.Context) {
			data := hardDiskData[hardDisk]
			c.JSON(200, gin.H{"data": data})
		})
	}

	// GET API to serve all hard disk used data
	router.GET("/data/harddisk/all", func(c *gin.Context) {
		allHardDiskData := []map[string]interface{}{}
		for hardDisk, data := range hardDiskData {
			allHardDiskData = append(allHardDiskData, map[string]interface{}{
				"hardDisk": hardDisk,
				"data":     data,
			})
		}
		c.JSON(200, gin.H{"hardDisks": allHardDiskData})
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
		time.Sleep(10 * time.Minute) // Adjust interval as needed
	}
}
