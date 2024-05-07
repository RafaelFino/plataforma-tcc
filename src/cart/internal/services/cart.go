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
		storage:    storage.NewCart(config),
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
	cart := domain.NewCart(clientId)

	log.Printf("[services.Cart] Creating cart for client: %s -> CartID: %s", clientId, cart.ID)
	err := c.storage.CreateCart(cart)

	if err != nil {
		log.Printf("[services.Cart] Error creating cart: %s", err)
		return "", err
	}

	return cart.ID, nil
}

func (c *Cart) Get(cartId string) (*domain.Cart, error) {
	log.Printf("[services.Cart] Getting cart with ID: %s", cartId)
	return c.storage.Get(cartId)
}

func (c *Cart) DeleteCart(cartId string) error {
	log.Printf("[services.Cart] Deleting cart with ID: %s", cartId)
	cart.Status = CartStatus.Canceled
	return c.storage.UpdateCart(cartId)
}

func (c *Cart) AddProduct(cartId string, productId string, quantity float64) error {
	log.Printf("[services.Cart] Adding product to cart: %s -> %s", cartId, productId)
	cart, err := c.storage.Get(cartId)

	if err != nil {
		log.Printf("[services.Cart] Error getting cart: %s", err)
		return err
	}

	product, err := c.products.Get(productId)

	if err != nil {
		log.Printf("[services.Cart] Error getting product: %s", err)
		return err
	}

	if product == nil {
		log.Printf("[services.Cart] Product not found: %s", productId)
		return nil
	}

	item := &domain.CartItem{
		Product:  product,
		Quantity: quantity,
	}

	cart.AddItem(item)
	cart.Status = CartStatus.Progress

	return c.storage.UpdateCart(cart)
}

func (c *Cart) RemoveProduct(cartId string, productId string) error {
	log.Printf("[services.Cart] Removing product from cart: %s -> %s", cartId, productId)
	cart, err := c.storage.Get(cartId)

	if err != nil {
		log.Printf("[services.Cart] Error getting cart: %s", err)
		return err
	}

	cart.RemoveItem(productId)

	return c.storage.UpdateCart(cart)
}

func (c *Cart) GetByClient(clientId string) ([]*domain.Cart, error) {
	log.Printf("[services.Cart] Getting carts for client: %s", clientId)
	return c.storage.GetByClient(clientId)
}

func (c *Cart) Checkout(cartId string, currency string) error {
	cart, err := c.storage.Get(cartId)

	if err != nil {
		log.Printf("[services.Cart] Error getting cart: %s", err)
		return err
	}

	cart.CalcTotal(c.currencies.GetCurrencies())
	cart.Status = CartStatus.Done

	return c.storage.UpdateCart(cart)
}
