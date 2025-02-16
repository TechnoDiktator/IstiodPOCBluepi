package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/IstiodPOCBluepi/models"
)

type ServiceC struct {
	ServiceBURL string
}

// Initialize ServiceC
func NewServiceC() *ServiceC {
	return &ServiceC{
		ServiceBURL: os.Getenv("SERVICE_B_URL"), // Fetch Service B URL from env
	}
}

func main() {
	svcC := NewServiceC()
	r := gin.Default()

	// GET /products - Forward request to Service B
	r.GET("/products", func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization Header"})
			log.Println("Missing Authorization Header")
			return
		}

		req, err := http.NewRequest("GET", svcC.ServiceBURL+"/products", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request for Service B"})
			log.Println("Failed to create request for Service B")
			return
		}

		req.Header.Set("Authorization", token) // Forward JWT token to Service B

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to reach Service B"})
			log.Println("Failed to reach Service B")
			return
		}
		defer resp.Body.Close()

		var products []models.Product
		if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response from Service B"})
			log.Panicln("Failed to decode response from Service B")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":        products,
			"requested_at": time.Now(),
			"service":     "service-c",
		})
	})

	// POST /products - Forward request to Service B
	r.POST("/products", func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization Header"})
			log.Println("Missing Authorization Header")
			return
		}

		var p models.Product
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			log.Println("Invalid input")
			return
		}

		jsonData, _ := json.Marshal(p)
		req, err := http.NewRequest("POST", svcC.ServiceBURL+"/products", bytes.NewBuffer(jsonData))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request for Service B"})
			log.Println("Failed to create request for Service B")
			return
		}

		req.Header.Set("Authorization", token) // Forward JWT token to Service B
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to reach Service B"})
			log.Println("Failed to reach Service B")
			return
		}
		defer resp.Body.Close()

		var createdProduct models.Product
		if err := json.NewDecoder(resp.Body).Decode(&createdProduct); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response from Service B"})
			log.Println("Failed to decode response from Service B")
			return
		}

		c.JSON(http.StatusCreated, createdProduct)
	})

	fmt.Println("Service C is running on port 8083")
	r.Run(":8083")
}
