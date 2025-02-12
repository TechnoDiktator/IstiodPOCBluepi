package main

import (
	//"encoding/json"
	"fmt"
	"log"
	"os"
	//"strings"

	"net/http"

	"github.com/yourusername/IstiodPOCBluepi/models"

	//"github.com/yourusername/IstiodPOCBluepi/db/productcrud"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	//"github.com/golang-jwt/jwt/v5"
	"github.com/yourusername/IstiodPOCBluepi/serviceinit"
)

type ServiceB struct {
	Auth0Domain string
}

func NewServiceB() *ServiceB {
	
	return &ServiceB{
		Auth0Domain: os.Getenv("AUTH0_DOMAIN"), // Set this in environment variables
	}
}

// Validate JWT Token using Auth0 JWKS
// func (s *ServiceB) ValidateJWT(tokenString string) (*jwt.Token, error) {
// 	jwksURL := fmt.Sprintf("https://%s/.well-known/jwks.json", s.Auth0Domain)

// 	// Fetch Auth0 JWKS
// 	resp, err := http.Get(jwksURL)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to fetch JWKS: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	var jwks struct {
// 		Keys []map[string]interface{} `json:"keys"`
// 	}
// 	json.NewDecoder(resp.Body).Decode(&jwks)

// 	// Parse and validate the JWT
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		// Ensure the signing method is RS256
// 		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}

// 		// Find the public key matching the JWT `kid`
// 		kid := token.Header["kid"]
// 		for _, key := range jwks.Keys {
// 			if key["kid"] == kid {
// 				return []byte(key["x5c"].(string)), nil
// 			}
// 		}
// 		return nil, fmt.Errorf("no matching key found")
// 	})

// 	if err != nil {
// 		return nil, err
// 	}
// 	return token, nil
// }

func main() {
	fmt.Println("Connecting To DB ================== Starting Service B")

	// Read database credentials from environment variables
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	// Initialize service with DB connection
	service, err := serviceinit.NewService(dsn)
	if err != nil {
		log.Fatalf("Failed to initialize services: %v", err)
	}

	r := gin.Default()

	// GET /products - Fetch all products (Istio handles authentication)
	r.GET("/products", func(c *gin.Context) {
		products, err := service.DBService.GetProducts()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
			return
		}
		c.JSON(http.StatusOK, products)
	})

	// POST /products - Create a new product (Istio handles authentication)
	r.POST("/products", func(c *gin.Context) {
		var p models.Product
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if err := service.DBService.CreateProduct(p); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully"})
	})

	fmt.Println("Service B is running on port 8081")
	r.Run(":8081")
}


