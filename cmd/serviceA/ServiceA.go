package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/IstiodPOCBluepi/models"
)

// ServiceA structure
type ServiceA struct {
	ServiceCURL string
}

// Initialize ServiceA
func NewServiceA() *ServiceA {
	return &ServiceA{
		ServiceCURL: "http://service-c.default.svc.cluster.local:8083", // Service C URL inside Kubernetes
	}
}

func main() {
	// Initialize Gin router
	r := gin.Default()
	svcA := NewServiceA()

	// Forward JWT Token in requests to Service C
	r.GET("/products", func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization Header"})
			return
		}

		req, _ := http.NewRequest("GET", svcA.ServiceCURL+"/products", nil)
		req.Header.Set("Authorization", token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to contact Service C"})
			return
		}
		defer resp.Body.Close()

		var products []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response"})
			return
		}

		// Add metadata
		c.JSON(http.StatusOK, gin.H{
			"data":       products,
			"requested_at": time.Now(),
			"service":     "service-a",
		})
	})

	r.POST("/products", func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization Header"})
			return
		}

		var p models.Product
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		jsonData, _ := json.Marshal(p)
		req, err := http.NewRequest("POST", svcA.ServiceCURL+"/products", bytes.NewBuffer(jsonData))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request for Service C"})
			return
		}

		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to contact Service C"})
			return
		}
		defer resp.Body.Close()

		var createdProduct models.ProductWithMetadata
		if err := json.NewDecoder(resp.Body).Decode(&createdProduct); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response from Service C"})
			return
		}

		c.JSON(http.StatusCreated, createdProduct)
	})

	fmt.Println("Service A is running on port 8082")
	r.Run(":8082")
}



