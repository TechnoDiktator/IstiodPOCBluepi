package main

import (
	"database/sql"

	"fmt"
	"log"

	"github.com/yourusername/IstiodPOCBluepi/models"
	//"github.com/yourusername/IstiodPOCBluepi/db/productcrud"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/yourusername/IstiodPOCBluepi/serviceinit"
)

var db *sql.DB

func connectDB() {
	var err error
	dsn := "root:password@tcp(localhost:3306)/testdb"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Database is unreachable: %v", err)
	}
	fmt.Println("Connected to MySQL successfully!")
}

func main() {
	dsn := "root:password@tcp(localhost:3306)/testdb"
	service, err := serviceinit.NewService(dsn)
	if err != nil {
		log.Fatalf("Failed to initialize services: %v", err)
	}

	r := gin.Default()
	r.GET("/products", func(c *gin.Context) {
		products, err := service.DBService.GetProducts()
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

		if err := service.DBService.CreateProduct(p); err != nil {
			c.JSON(500, gin.H{"error": "Failed to create product"})
			return
		}

		c.JSON(201, gin.H{"message": "Product created successfully"})
	})

	fmt.Println("Service B is running on port 8081")
	r.Run(":8081")
}
