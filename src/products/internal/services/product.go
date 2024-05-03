package services

import (
	"crypto/rand"
	"errors"
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

func (c *Product) Get() ([]*domain.Product, error) {
	log.Printf("[services.Product] Getting all products")
	return c.storage.GetAll()
}

func (c *Product) Update(product *domain.Product) error {
	if product == nil {
		log.Printf("[services.Product] product is nil")
		return errors.New("product is nil")
	}

	if product.ID == "" {
		log.Printf("[services.Product] product ID is empty")
		return errors.New("product ID is empty")
	}

	log.Printf("[services.Product] Updating product with ID: %s to %+v", product.ID, product)
	return c.storage.Update(product)
}

func (c *Product) Insert(product *domain.Product) (string, error) {
	if product == nil {
		log.Printf("[services.Product] product is nil")
		return "", errors.New("product is nil")
	}

	product.ID = c.CreateID()

	log.Printf("[services.Product] Inserting product with ID: %s", product.ID)
	return product.ID, c.storage.Insert(product)
}

func (c *Product) Delete(id string) error {
	log.Printf("[services.Product] Deleting product with ID: %s", id)
	return c.storage.Delete(id)
}

func (c *Product) CreateID() string {
	return ulid.MustNew(ulid.Now(), ulid.Monotonic(rand.Reader, 0)).String()
}
