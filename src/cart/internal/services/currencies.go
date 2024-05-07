package services

import (
	"cart/internal/config"
	"log"
	"net/http"
	"products/internal/domain"
	"time"

	"github.com/goccy/go-json"
)

type Currencies struct {
	data map[string]*domain.Currency
	url  string
	last time.Time
	ttl  int
}

type CurrencieData struct {
	Currencies []*domain.Currency `json:"currencies"`
	Elapsed    float64            `json:"elapsed"`
	Timestamp  time.Time          `json:"timestamp"`
}

func NewCurrencies(config *config.Config) *Currencies {
	return &Currencies{
		data: make(map[string]*domain.Currency),
		url:  config.CurrenciesURL,
		ttl:  config.CurrenciesInterval,
		last: time.Now().Local().Add(-time.Hour * 24),
	}
}

func (c *Currencies) Update() error {
	if time.Since(c.last).Minutes() > float64(c.ttl) {
		return nil
	}

	log.Printf("[services.Currencies] Updating currencies")

	body, status, err := HttpGet(c.url)

	if err != nil {
		log.Printf("[services.Currencies] Error getting currencies: %s", err)
		return err
	}

	if status != http.StatusOK {
		log.Printf("[services.Currencies] Error getting currencies: %d", status)
		return err
	}

	data := &CurrencieData{}

	err = json.Unmarshal([]byte(body), data)

	if err != nil {
		log.Printf("[services.Currencies] Error parsing currencies: %s", err)
		return err
	}

	c.data = data.Currencies

	for _, currency := range data.Currencies {
		c.data[currency.Code] = currency
		log.Printf("[services.Cart] currency: %s -> %+v", currency.Code, currency)
	}

	c.last = time.Now()

	return nil
}

func (c *Currencies) GetAll() map[string]*domain.Currency {
	err := c.Update()

	if err != nil {
		log.Printf("[services.Cart] Error updating currencies: %s", err)
	}

	return c.data
}

func (c *Currencies) Get(code string) *domain.Currency {
	err := c.Update()

	if err != nil {
		log.Printf("[services.Cart] Error updating currencies: %s", err)
	}

	if ret, ok := c.data[code]; ok {
		return ret
	}

	log.Printf("[services.Currencies] Currency not found: %s", code)

	return &domain.Currency{
		Code:      code,
		Name:      "Unknown",
		High:      1,
		Low:       1,
		VarBid:    1,
		PctChange: 1,
		Bid:       1,
		Ask:       1,
	}
}
