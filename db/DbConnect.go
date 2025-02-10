package db

import (
	"database/sql"
	"fmt"
)

func connectDB() (*sql.DB, error) {

	dsn := "root:password@tcp(localhost:3306)/testdb"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("database is unreachable: %v", err)
	}
	fmt.Println("Connected to MySQL successfully!")
	return db, nil
}
