package services

import (
	"cart/internal/config"
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

/*
	s.engine.POST("/cart/", s.handler.CreateCart)
	s.engine.GET("/cart/:cart_id", s.handler.Get)
	s.engine.DELETE("/cart/:cart_id", s.handler.DeleteCart)

	s.engine.POST("/cart/:cart_id", s.handler.AddProduct)
	s.engine.DELETE("/cart/:cart_id/:product_id", s.handler.RemoveProduct)

	s.engine.GET("/client/:client_id", s.handler.GetByClient)
	s.engine.PUT("/cart/:cart_id", s.handler.Checkout)
*/
