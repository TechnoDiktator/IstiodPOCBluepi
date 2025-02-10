package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/IstiodPOCBluepi/models"
)

type ServiceA struct {
	ServiceBURL string
}

func (s *ServiceA) FetchProducts() ([]models.ProductWithMetadata, error) {
	resp, err := http.Get(s.ServiceBURL + "/products")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %v", err)
	}
	defer resp.Body.Close()

	var products []models.ProductWithMetadata
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	for i := range products {
		products[i].RequestedAt = time.Now()
		products[i].ServiceAID = "service-a-123"
	}

	return products, nil
}

func (s *ServiceA) CreateProduct(p models.Product) (*models.ProductWithMetadata, error) {
	jsonData, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal product data: %v", err)
	}

	resp, err := http.Post(s.ServiceBURL+"/products", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to send request to Service B: %v", err)
	}
	defer resp.Body.Close()

	var createdProduct models.ProductWithMetadata
	if err := json.NewDecoder(resp.Body).Decode(&createdProduct); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	createdProduct.RequestedAt = time.Now()
	createdProduct.ServiceAID = "service-a-123"

	return &createdProduct, nil
}

func main() {
	r := gin.Default()
	svcA := &ServiceA{ServiceBURL: "http://localhost:8081"}

	r.GET("/products", func(c *gin.Context) {
		products, err := svcA.FetchProducts()
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch products"})
			return
		}
		c.JSON(200, products)
	})

	r.POST("/products", func(c *gin.Context) {
		var p models.Product
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		createdProduct, err := svcA.CreateProduct(p)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create product in Service B"})
			return
		}

		c.JSON(201, createdProduct)
	})

	fmt.Println("Service A is running on port 8082")
	r.Run(":8082")
}
