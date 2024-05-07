package services

import (
	"cart/internal/config"
	"cart/internal/domain"
	"cart/internal/storage"
	"log"
)

type Cart struct {
	storage *storage.Cart
	Config  *config.Config

	products   *Products
	clients    *Clients
	currencies *Currencies
}

func NewCart(config *config.Config) *Cart {
	ret := &Cart{
		storage:    storage.NewCart(storage.NewDbConnection(config.DBPath)),
		Config:     config,
		currencies: NewCurrencies(config),
		clients:    NewClients(config),
		products:   NewProducts(config),
	}

	err := ret.currencies.Update()

	if err != nil {
		log.Printf("[services.Cart] Error updating currencies: %s", err)
	}

	return ret
}

func (c *Cart) Close() error {
	log.Printf("[services.Cart] Closing storage")
	return c.storage.Close()
}

func (c *Cart) CreateCart(clientId string) (string, error) {
}

func (c *Cart) Get(cartId string) (*domain.Cart, error) {
}

func (c *Cart) DeleteCart(cartId string) error {
}

func (c *Cart) AddProduct(cartId string, productId string, quantity int) error {
}

func (c *Cart) RemoveProduct(cartId string, productId string) error {
}

func (c *Cart) GetByClient(clientId string) ([]*domain.Cart, error) {
}

func (c *Cart) Checkout(cartId string) error {
}
