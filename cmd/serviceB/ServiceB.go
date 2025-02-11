package main

import (

	"fmt"
	"log"

	"github.com/yourusername/IstiodPOCBluepi/models"
	//"github.com/yourusername/IstiodPOCBluepi/db/productcrud"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/yourusername/IstiodPOCBluepi/serviceinit"
)



func main() {
	dsn := "root:PlusOne98@17@tcp(127.0.0.1:3306)/testdb"
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
