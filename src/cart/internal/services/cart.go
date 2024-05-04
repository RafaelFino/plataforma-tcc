package services

import (
	"cart/internal/config"
	"cart/internal/domain"
	"cart/internal/storage"
	"io"
	"json"
	"log"
	"net/http"
	"time"

	"github.com/goccy/go-json"
)

type Cart struct {
	storage    *storage.Cart
	currencies map[string]float64
	Config     *config.Config
	last       time.Time
}

func NewCart(config *config.Config) *Cart {
	ret := &Cart{
		storage:    storage.NewCart(storage.NewDbConnection(config.DBPath)),
		Config:     config,
		currencies: make(map[string]float64),
	}

	err := ret.UpdateCurrencies()

	if err != nil {
		log.Printf("[services.Cart] Error updating currencies: %s", err)
	}

	return ret
}

func (c *Cart) httpGet(url string) (string, int, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Printf("[services.Cart] Error getting url: %s", err)
		return "", http.StatusInternalServerError, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("[services.Cart] Error reading body: %s", err)
		return "", res.StatusCode, err
	}

	if s.Config.Debug {
		log.Printf("[services.Cart] HTTP-GET %d Response: %s", res.StatusCode, body)
	}

	return string(body), res.StatusCode, nil
}

func (c *Cart) UpdateCurrencies() error {
	log.Printf("[services.Cart] Getting currencies")

	body, status, err := c.httpGet(c.Config.CurrenciesURL)

	if err != nil {
		log.Printf("[services.Cart] Error getting currencies: %s", err)
		return err
	}

	if status != http.StatusOK {
		log.Printf("[services.Cart] Error getting currencies: %d", status)
		return err
	}

	jsonData := make(map[string]interface{})

	err = json.Unmarshal([]byte(body), &jsonData)

	if err != nil {
		log.Printf("[services.Cart] Error parsing currencies: %s", err)
		return err
	}

	if data, ok := jsonData["currencies"].(map[string]interface{}); ok {
		log.Printf("[services.Cart] Parsing currencies: %+v", data)

		for _, item := range data {
			log.Printf("[services.Cart] Parsing item: %+v", item)
			curr := item.(map[string]interface{})

			if val, ok := curr["code"].(string); ok {
				c.currencies[val] = curr["value"].(float64)
				log.Printf("[services.Cart] Add currency: %s -> %+v", val, c.currencies[val])
			} else {
				log.Printf("[services.Cart] Error parsing currency code -> %+v", curr)
			}
		}
	}

	c.last = time.Now()

	return nil
}

func (c *Cart) SetClientDetails(cart *domain.Cart) error {
	log.Printf("[services.Cart] Getting client details for cart: %s", cart.ClientID)

	body, status, err := c.httpGet(c.Config.ClientsURL)

	if err != nil {
		log.Printf("[services.Cart] Error getting client details: %s", err)
		return err
	}

	if status != http.StatusOK {
		log.Printf("[services.Cart] Error getting client details: %d", status)
		return err
	}

	jsonData := make(map[string]interface{})

	err = json.Unmarshal([]byte(body), &jsonData)

	if err != nil {
		log.Printf("[services.Cart] Error parsing client details: %s", err)
		return err
	}

	if data, ok := jsonData["client"].(map[string]interface{}); ok {
		log.Printf("[services.Cart] Parsing client details: %+v", data)

		if val, ok := data["name"].(string); ok {
			cart.ClientName = val
		}

		if val, ok := data["email"].(string); ok {
			cart.ClientEmail = val
		}

		if val, ok := data["surname"].(string); ok {
			cart.ClientSurname = val
		}

		log.Printf("[services.Cart] Client details: %+v", cart)
	}

	return nil
}

func (c *Cart) Close() error {
	log.Printf("[services.Cart] Closing storage")
	return c.storage.Close()
}

func (c *Cart) CreateCart(clientID string) (string, error) {
	log.Printf("[services.Cart] Creating cart for client: %s", clientID)

	cart := &domain.NewCart(clientID, currencies)

	ret, err := c.storage.CreateCart(clientID)

	if err != nil {
		log.Printf("[services.Cart] Error creating cart: %s", err)
		return "", err
	}

	log.Printf("[services.Cart] Cart created: %+v", ret)

	return ret, nil
}

func (c *Cart) GetById(id string) (*domain.Cart, error) {
	log.Printf("[services.Cart] Getting cart with ID: %s", id)

	ret, err := c.storage.GetCart(id)

	if err != nil {
		log.Printf("[services.Cart] Error getting cart: %s", err)
		return nil, err
	}

	log.Printf("[services.Cart] Cart: %+v", ret)

	return ret, nil
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
