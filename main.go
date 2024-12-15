package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

// Structure to store the last heartbeat time for each computer
type ComputerStatus struct {
	LastSeen time.Time
	Status   string
}

var (
	statusMap = make(map[string]*ComputerStatus) // Map to store computer statuses
	mu        sync.Mutex                         // Mutex to protect access to statusMap
	timeout   = 8 * time.Minute                  // Timeout duration
)

// Maps to store data for multiple clients and bot logs
var clientData = make(map[string][][]interface{})
var botLogData = make(map[string][][]interface{})   // Separate storage for bot logs
var hardDiskData = make(map[string][][]interface{}) // Separate storage for hardDiskData logs

// Function to send email
func sendEmail(subject, body string) error {

	// Email configuration
	smtpHost := "smtp.gmail.com" // Replace with your SMTP server address
	smtpPort := 587              // SMTP server port (use 465 for SSL or 587 for STARTTLS)
	emailFrom := "maintenancegreenmoons@gmail.com"
	emailPassword := "aoqtlaepsucvdksf" // Replace with your email password or app-specific password
	emailTo := "maintenancegreenmoons@gmail.com"

	// Create a new message
	m := gomail.NewMessage()
	m.SetHeader("From", emailFrom)
	m.SetHeader("To", emailTo)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	// Create a new dialer
	d := gomail.NewDialer(smtpHost, smtpPort, emailFrom, emailPassword)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Error sending email:", err)
		return err
	}

	fmt.Println("Email sent successfully!")

	return nil
}

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

	//Part   check remote computer status online or not
	// Define a route to check if remote computer is online
	// Define a route to handle "check-online" requests
	router.POST("/check-online", func(c *gin.Context) {
		var requestData struct {
			ComputerName string `json:"computer_name"` // Name of the remote computer
			Status       string `json:"status"`        // "online" or "offline"
		}

		// Parse JSON from the client
		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request format",
			})
			return
		}

		// Update the statusMap with the current timestamp
		mu.Lock()
		if _, exists := statusMap[requestData.ComputerName]; !exists {
			statusMap[requestData.ComputerName] = &ComputerStatus{
				LastSeen: time.Now(),
				Status:   "offline", // Default to offline
			}
		}
		statusMap[requestData.ComputerName].LastSeen = time.Now()
		statusMap[requestData.ComputerName].Status = requestData.Status
		mu.Unlock()

		c.JSON(http.StatusOK, gin.H{
			"message":         "Status updated to online",
			"computer_name":   requestData.ComputerName,
			"received_status": requestData.Status,
		})
	})

	// Goroutine to check for inactive computers
	go func() {
		for {
			time.Sleep(1 * time.Minute) // Check every minute
			mu.Lock()
			for computerName, status := range statusMap {
				if time.Since(status.LastSeen) > timeout {
					status.Status = "offline"
					go sendOfflineStatus(computerName) // Send API asynchronously
				}
			}
			mu.Unlock()
		}
	}()

	// Define a route to get the current status of all computers
	router.GET("/statuses", func(c *gin.Context) {
		mu.Lock()
		defer mu.Unlock()

		response := make(map[string]string)
		emailBody := "The following computers are offline:\n"
		offlineFound := false

		for computerName, status := range statusMap {
			response[computerName] = status.Status

			if status.Status == "offline" {
				offlineFound = true
				emailBody += fmt.Sprintf("- %s\n", computerName)
			}
		}

		if offlineFound {
			if err := sendEmail("Offline Computers Alert", emailBody); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email", "details": err.Error()})
				return
			}
		}

		c.JSON(http.StatusOK, response)
	})

	// Define a route to get the current status of all computers
	/*
		router.GET("/statuses", func(c *gin.Context) {
			mu.Lock()
			defer mu.Unlock()

			response := make(map[string]string)
			for computerName, status := range statusMap {
				response[computerName] = status.Status
			}
			c.JSON(http.StatusOK, response)
		})
	*/

	//Part  errorlog data and botprocesslog data
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
	hardDiskUsedEndpoints := []string{"hcs", "worldair", "log2", "log3"} // Add more hard disk categories as needed

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

// Function to send an API request to mark a computer as offline
func sendOfflineStatus(computerName string) {
	apiURL := "https://btestapi-am67.onrender.com/check-online" // Self API endpoint

	// Prepare payload
	payload := map[string]string{
		"computer_name": computerName,
		"status":        "offline",
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return
	}

	// Send POST request
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}
	defer resp.Body.Close()
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
