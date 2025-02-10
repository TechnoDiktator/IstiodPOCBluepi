package db

import (
	"github.com/yourusername/IstiodPOCBluepi/models"
)

type DB interface {
	GetProducts() ([]models.Product, error)
	CreateProduct(p models.Product) error
}
