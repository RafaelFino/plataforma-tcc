package services

import (
	"crypto/rand"
	"log"
	"products/internal/config"
	"products/internal/domain"
	"products/internal/storage"

	"github.com/oklog/ulid"
)

type Product struct {
	storage *storage.Product
}

func NewProduct(config *config.Config) *Product {
	return &Product{
		storage: storage.NewProduct(storage.NewDbConnection(config.DBPath)),
	}
}

func (c *Product) GetById(id string) (*domain.Product, error) {
	log.Printf("[services.Product] Getting product with ID: %s", id)
	ret, err := c.storage.Get(id)

	if err != nil {
		log.Printf("[services.Product] Error getting product: %s", err)
		return nil, err
	}

	return ret, nil
}

func (c *Product) CreateID() string {
	return ulid.MustNew(ulid.Now(), ulid.Monotonic(rand.Reader, 0)).String()
}
