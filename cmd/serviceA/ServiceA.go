package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/go-resty/resty/v2"

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
	ServiceBURL    string
	Auth0Domain    string
	Audience       string
	ClientID       string
	ClientSecret   string
}

func NewServiceA() *ServiceA {
	return &ServiceA{
		ServiceBURL:    "http://service-b.default.svc.cluster.local:8081",
		Auth0Domain:    os.Getenv("AUTH0_DOMAIN"),
		Audience:       os.Getenv("AUTH0_AUDIENCE"),
		ClientID:       os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret:   os.Getenv("AUTH0_CLIENT_SECRET"),
	}
}
//C:\Users\tarang\Desktop\IstiodPOC\IstiodPOCBluepi\cmd\serviceA\ServiceA.go



// Fetch JWT Token from Auth0
func (s *ServiceA) getJWTToken() (string, error) {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{
			"grant_type":    "client_credentials",
			"client_id":     s.ClientID,
			"client_secret": s.ClientSecret,
			"audience":      s.Audience,
		}).
		Post(fmt.Sprintf("https://%s/oauth/token", s.Auth0Domain))

	if err != nil {
		return "", fmt.Errorf("failed to request JWT token: %v", err)
	}

	var result map[string]string
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return "", fmt.Errorf("failed to decode JWT response: %v", err)
	}

	token, ok := result["access_token"]
	if !ok {
		return "", fmt.Errorf("JWT token missing in response")
	}

	return token, nil
}


// Middleware to Validate JWT Token
func (s *ServiceA) ValidateJWT(tokenString string) (*jwt.Token, error) {
	jwksURL := fmt.Sprintf("https://%s/.well-known/jwks.json", s.Auth0Domain)

	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %v", err)
	}
	defer resp.Body.Close()

	var jwks struct {
		Keys []map[string]interface{} `json:"keys"`
	}
	json.NewDecoder(resp.Body).Decode(&jwks)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		kid := token.Header["kid"]
		for _, key := range jwks.Keys {
			if key["kid"] == kid {
				return []byte(key["x5c"].(string)), nil
			}
		}
		return nil, fmt.Errorf("no matching key found")
	})

	if err != nil {
		return nil, err
	}
	return token, nil
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
	svcA := NewServiceA()

	authMiddleware := func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization Header"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := svcA.ValidateJWT(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or Expired Token"})
			c.Abort()
			return
		}

		c.Next()
	}

	r.GET("/products",authMiddleware ,  func(c *gin.Context) {

		tokenString, _ := svcA.getJWTToken()
		req, _ := http.NewRequest("GET", svcA.ServiceBURL+"/products", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch products"})
			return
		}
		defer resp.Body.Close()

		var products []map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&products)

		c.JSON(200, products)
	})

	r.POST("/products", authMiddleware, func(c *gin.Context) {
		var p models.Product
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}
	
		// Extract the JWT token from the request
		tokenString := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	
		// Send request to Service B with JWT token
		jsonData, _ := json.Marshal(p)
		req, err := http.NewRequest("POST", svcA.ServiceBURL+"/products", bytes.NewBuffer(jsonData))
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create request for Service B"})
			return
		}
		req.Header.Set("Authorization", "Bearer "+tokenString)
		req.Header.Set("Content-Type", "application/json")
	
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to send request to Service B"})
			return
		}
		defer resp.Body.Close()
	
		var createdProduct models.ProductWithMetadata
		if err := json.NewDecoder(resp.Body).Decode(&createdProduct); err != nil {
			c.JSON(500, gin.H{"error": "Failed to decode response from Service B"})
			return
		}
	
		c.JSON(201, createdProduct)
	})

	fmt.Println("Service A is running on port 8082")
	r.Run(":8082")
}
