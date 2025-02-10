package models

import (
	"time"
)

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type ProductWithMetadata struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Price       int       `json:"price"`
	RequestedAt time.Time `json:"requested_at"`
	ServiceAID  string    `json:"service_a_id"`
}
