package db

import (
	"database/sql"
	"fmt"
	"log"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"

	"github.com/yourusername/IstiodPOCBluepi/models"
)

type MySQLDB struct {
	Conn *sql.DB
}

func NewMySQLDB(dsn string) (*MySQLDB, error) {
	log.Println("🚀 Connecting to MySQL with DSN:", dsn)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Panicln("❌ Error opening database connection:", err)  // ✅ Log error
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	// Test connection
	if err = db.Ping(); err != nil {
		log.Panicln("❌ Database is unreachable:", err)  // ✅ Log error
		return nil, fmt.Errorf("database is unreachable: %v", err)
	}

	log.Println("✅ Connected to MySQL successfully!")
	return &MySQLDB{Conn: db}, nil
}

func (m *MySQLDB) GetProducts() ([]models.Product, error) {
	log.Println("================ GetProducts Start =================")

	// Check if DB connection is nil


	log.Println("=================1=================")

	if m.Conn == nil {
		log.Println(" Database connection is nil!")
		return nil, fmt.Errorf("database connection is nil")
	}

	log.Println("=================2=================")
	rows, err := m.Conn.Query("SELECT id, name, price FROM products")
	if err != nil {
		log.Println(" Query failed:", err)  // Print the error
		return nil, err
	}
	defer rows.Close()

	log.Println("=================3=================")
	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			log.Println(" Error scanning row:", err)  //  Print scan errors
			return nil, err
		}
		products = append(products, p)
	}
	log.Println("=================4=================")
	if len(products) == 0 {
		log.Println("⚠️ No products found, returning empty array")
	}else{
		productJSON, _ := json.MarshalIndent(products, "", "  ")
		log.Println("Products found:", string(productJSON))
	}

	log.Println("=================5=================")
	log.Println("================ GetProducts End =================")
	return products, nil
}


func (m *MySQLDB) CreateProduct(p models.Product) error {
	log.Println("================ CreateProduct Start =================")
	_, err := m.Conn.Exec("INSERT INTO products (name, price) VALUES (?, ?)", p.Name, p.Price)
	log.Println("================ CreateProduct End =================")
	if err != nil {
		log.Println("Error in CreateProduct")
		log.Println(err)
	}
	return err
}
