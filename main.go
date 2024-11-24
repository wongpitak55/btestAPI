package main

import (
	"fmt" // For console output
	"net/http"
	"os" // For environment variables

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

		// Build the HTML content with embedded CSS and JavaScript
		htmlContent := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Data Table</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					margin: 20px;
					background-color: #f4f4f4;
				}
				h1 {
					text-align: center;
					color: #333;
				}
				table {
					width: 100%;
					border-collapse: collapse;
					margin: 20px 0;
					background: #fff;
					box-shadow: 0 2px 5px rgba(0,0,0,0.2);
				}
				table th, table td {
					border: 1px solid #ddd;
					padding: 10px;
					text-align: center;
				}
				table th {
					background-color: #333;
					color: #fff;
				}
				table tr:nth-child(even) {
					background-color: #f9f9f9;
				}
			</style>
			<script>
				function alertRow(rowIndex) {
					alert('You clicked on row: ' + rowIndex);
				}
			</script>
		</head>
		<body>
			<h1>Received Data Table</h1>
			<table>
				<tr>
					<th>Index</th>
					<th>Values</th>
				</tr>
		`

		// Append table rows
		for i, row := range data {
			htmlContent += fmt.Sprintf(`
				<tr onclick="alertRow(%d)">
					<td>%d</td>
					<td>%v</td>
				</tr>
			`, i, i, row)
		}

		// Close HTML tags
		htmlContent += `
			</table>
			<p style="text-align: center;">Click on any row to see the row index.</p>
		</body>
		</html>
		`

		// Render the HTML content as a response
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
	})

	// Print a message when the server starts
	fmt.Printf("Server is starting on port %s...\n", port)

	// Start the server on the specified port
	router.Run(":" + port)
}
