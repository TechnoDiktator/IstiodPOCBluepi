package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type ServiceC struct {
	ServiceBURL string
	Auth0Domain string
}

func NewServiceC() *ServiceC {
	return &ServiceC{
		ServiceBURL: os.Getenv("SERVICE_B_URL"), // Fetch Service B URL from env
		Auth0Domain: os.Getenv("AUTH0_DOMAIN"),  // Fetch Auth0 domain from env
	}
}

// Validate JWT Token using Auth0 JWKS
func (s *ServiceC) ValidateJWT(tokenString string) (*jwt.Token, error) {
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
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

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

// Call Service B
func (s *ServiceC) ForwardRequestToServiceB(token string) ([]byte, error) {
	req, err := http.NewRequest("GET", s.ServiceBURL+"/products", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := json.Marshal(resp.Body)
	return body, nil
}

func main() {
	svcC := NewServiceC()
	r := gin.Default()

	authMiddleware := func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization Header"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := svcC.ValidateJWT(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or Expired Token"})
			c.Abort()
			return
		}

		c.Next()
	}

	r.GET("/products", authMiddleware, func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		response, err := svcC.ForwardRequestToServiceB(authHeader)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch data from Service B"})
			return
		}
		c.Data(200, "application/json", response)
	})

	fmt.Println("Service C is running on port 8083")
	r.Run(":8083")
}
