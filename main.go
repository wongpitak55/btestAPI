package main

import (
	"fmt" // To print to the console
	"net/http"
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
		// Print the incoming request to the console
		fmt.Println("Received a POST request to /api")

		var requestBody map[string]interface{}
		if err := c.BindJSON(&requestBody); err != nil {
			// Log the error
			fmt.Println("Error: Invalid JSON payload")
			c.JSON(400, gin.H{"error": "Invalid JSON payload"})
			return
		}

		// Extract "data" as a list of arrays
		data, exists := requestBody["data"].([]interface{})
		if !exists {
			// Log missing "data" field
			fmt.Println("Error: Missing or invalid 'data' field in request payload")
			c.JSON(400, gin.H{"error": "Missing or invalid 'data' field"})
			return
		}

		// Print the extracted data
		fmt.Printf("Received 'data': %v\n", data)

		// Create an HTML table
		htmlContent := "<html><body><table border='1'>"
		htmlContent += "<tr><th>Index</th><th>Values</th></tr>"

		for i, row := range data {
			htmlContent += fmt.Sprintf("<tr><td>%d</td><td>%v</td></tr>", i, row)
		}

		htmlContent += "</table></body></html>"

		// Render the HTML table as a response
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
	})

	// Print a message when the server starts
	fmt.Printf("Server is starting on port %s...\n", port)

	// Start the server on the specified port
	router.Run(":" + port)
}
