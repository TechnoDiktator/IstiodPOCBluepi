package db

import (
	"database/sql"
	"fmt"

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
	rows, err := m.Conn.Query("SELECT id, name, price FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (m *MySQLDB) CreateProduct(p models.Product) error {
	_, err := m.Conn.Exec("INSERT INTO products (name, price) VALUES (?, ?)", p.Name, p.Price)
	return err
}
