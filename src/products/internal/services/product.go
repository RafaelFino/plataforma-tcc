package services

import (
	"crypto/rand"
	"products/internal/storage"

	"github.com/oklog/ulid"
)

type Product struct {
	storage *storage.Product
}

func NewProduct(storage *storage.Product) *Product {
	return &Product{
		storage: storage,
	}
}

func (c *Product) CreateID() string {
	return ulid.MustNew(ulid.Now(), ulid.Monotonic(rand.Reader, 0)).String()
}
