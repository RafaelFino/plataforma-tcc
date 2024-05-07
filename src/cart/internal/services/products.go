package services

import (
	"cart/internal/config"
	"cart/internal/domain"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/goccy/go-json"
)

type Products struct {
	url string
}

type ProductData struct {
	Product   *domain.Product `json:"product"`
	Elapsed   float64         `json:"elapsed"`
	Timestamp time.Time       `json:"timestamp"`
}

func NewProducts(config *config.Config) *Products {
	return &Products{
		url: config.ProductsURL,
	}
}

func (p *Products) Get(id string) (*domain.Product, error) {
	log.Printf("[services.Products] Get product info for ID: %s", id)

	body, status, err := HttpGet(p.url)

	if err != nil {
		log.Printf("[services.Products] Error getting product info: %s", err)
		return nil, err
	}

	if status != http.StatusOK {
		log.Printf("[services.Products] Error getting product info: %d", status)
		return nil, err
	}

	data := &ProductData{}

	err = json.Unmarshal([]byte(body), data)

	if err != nil {
		log.Printf("[services.Products] Error parsing product info: %s", err)
		return nil, err
	}

	if data.Product == nil {
		log.Printf("[services.Products] Product is nil")
		return nil, errors.New("product is nil")
	}

	log.Printf("[services.Products] Product info: %+v", data.Product)

	return data.Product, nil
}
