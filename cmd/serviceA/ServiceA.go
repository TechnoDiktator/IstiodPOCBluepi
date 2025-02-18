package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"net/http"
	"io"
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

		log.Println("Response from Service C:", resp.Status)

		defer resp.Body.Close()

		// ✅ Read the body before decoding
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
			return
		}
		// ✅ Print raw response before unmarshalling
		bodyString := string(bodyBytes)
		log.Println("🔹 Raw Response from Service C:", bodyString)

		// ✅ Unmarshal into map[string]interface{}
		var serviceCResp map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &serviceCResp); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response FROM Service C"})
			return
		}

		// ✅ Extract "data" field
		products, ok := serviceCResp["data"].([]interface{})
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response format from Service C"})
			return
		}

		// ✅ Return only the "data" field from Service C
		c.JSON(http.StatusOK, gin.H{
			"data":        products,  // ✅ Extracted array
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



