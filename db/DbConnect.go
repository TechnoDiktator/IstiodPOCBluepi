package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/yourusername/IstiodPOCBluepi/models"
)

type MySQLDB struct {
	Conn *sql.DB
}

func NewMySQLDB(dsn string) (*MySQLDB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("database is unreachable: %v", err)
	}
	fmt.Println("Connected to MySQL successfully!")
	return &MySQLDB{Conn: db}, nil
}

func (m *MySQLDB) GetProducts() ([]models.Product, error) {
	log.Println("================ GetProducts Start =================")

	// Check if DB connection is nil
	if m.Conn == nil {
		log.Panicln(" Database connection is nil!")
		return nil, fmt.Errorf("database connection is nil")
	}

	rows, err := m.Conn.Query("SELECT id, name, price FROM products")
	if err != nil {
		log.Panicln(" Query failed:", err)  // ✅ Print the error
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			log.Panicln(" Error scanning row:", err)  // ✅ Print scan errors
			return nil, err
		}
		products = append(products, p)
	}
	if len(products) == 0 {
		log.Println("⚠️ No products found, returning empty array")
	}
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
