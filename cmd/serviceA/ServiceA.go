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

// {
// 	"sub": "google-oauth2|104494965273019075736",
// 	"given_name": "Tarang ",
// 	"family_name": "Rastogi",
// 	"nickname": "tarang",
// 	"name": "Tarang Rastogi",
// 	"picture": "https://lh3.googleusercontent.com/a-/ALV-UjVQ_lqYpdFVvlkchRr6AzwOjtnopQ2EMzIjUY1f_wV54pqD7zEpqSZMyLfeG97iZyiEJvvKnX1b7mXGvYlqfQ59JvwNOo7k_CLZ8xoMk0OlRFrvE4-oa6ypz05rHpoFcbxHN_37YOt128pMMfcCz8BMmwp45lUq5WkiDGKJseO0TANCd0hKWUytNTKm_JGsu4T-FiAe8AvfLdSeuACwe4oD-SANVBZipHQOcput8DUaLalnX1zgc9Xd7xIeps2f2ER8aBaxxNmS8rbc-ieivKHupjh_yWJ0FbcT567ox_0Z1udvp2VT5eeWrCs8d9IRJ8W_PN_CoUn3oH8Z8Ak6d1Efy2gQsDeLC2BME5yOOUa2ns8rERQ9SjphWKTmsO02D3qoK8hBPh3_ZMeFNU0mwUmMJmeR-iDfS4VqxfQnpOiVmQTCoKvOtw5FxjSCdrWjHPZxSThXVi4xMndF9je3x05wBGTaLUdSyma-I_V_iPPmgOkmvpuyOVEtuwRKMJzQZTu3sUVMJRqQKM3uwxu6fjhhVOnCl-JeCLlEBcTlLwgfb4-QlNSkwKDCdduR3zqY-6qnwABqhRDwXYgV8E3JVQ3d4p8vLZ8zfS-O7gmph_QZGE4vR48KxaMrgFx5Gx3z-uvzuorGwu0ZkafcWGAQ0TSnqFxwRj8e-wwmxFRc34G2jFHv6jXHXKxwfhrge15YNepWvplK9pq7ChWESFO1WIHsAuqn6qq-KF52OtmB_PdVqn4gvbE3P5QqcJxtDmCvePOEeuy290SgKdDWJ_pBnTDQgraJfJ_nlCNoo2NqOc6xUA2ToL2j6iQC-TlXi795ptFbo58vV7yVMXkrqmauSqUkyWrzP9M73EQXMhFITgbnUhBcb3ixfr4_f_c5NFIa-bfeF50MrV-frZG0m7-WZY_lvn_CBEPkfE8u7p0WknInV84Unhy_CCuXasRsbB-mk5cdbHyT21JrxpcRUzpWoNPhJ6w=s96-c",
// 	"updated_at": "2025-02-11T10:52:14.720Z"
//   }



type ServiceA struct {
	ServiceBURL string
}
//C:\Users\tarang\Desktop\IstiodPOC\IstiodPOCBluepi\cmd\serviceA\ServiceA.go

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
	svcA := &ServiceA{ServiceBURL: "http://service-b.default.svc.cluster.local:8081"}


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
