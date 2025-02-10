package db

import (
	"database/sql"
	"fmt"
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
