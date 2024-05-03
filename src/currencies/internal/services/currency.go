package services

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"currencies/internal/domain"
)

type Currency struct {
	url      string
	data     map[string]*domain.Currency
	last     time.Time
	interval int64
}

func NewCurrency(url string) *Currency {
	return &Currency{
		url:      url,
		interval: 60,
	}
}

func (c *Currency) Update() error {
	if time.Now().After(c.last.Add(time.Duration(c.interval) * time.Minute)) {
		log.Printf("[services.Currency] minimal interval not reached")
		return nil
	}

	ret, err := http.Get(c.url)
	if err != nil {
		log.Printf("[services.Currency] Error getting url: %s", err)
		return err
	}

	defer ret.Body.Close()
	body, err := io.ReadAll(ret.Body)
	if err != nil {
		log.Printf("[services.Currency] Error reading body: %s", err)
		return err
	}

	data := string(body)

	log.Printf("[services.Currency] body: %s", data)

	c.last = time.Now()

	err = c.parserData(data)

	if err != nil {
		log.Printf("[services.Currency] Error parsing data: %s", err)
	}

	go func() {
		log.Printf("[services.Currency] waiting %d minutes to update", c.interval)
		time.Sleep(time.Duration(c.interval) * time.Minute)
		c.Update()
	}()

	return err
}

func (c *Currency) GetCurrency() map[string]*domain.Currency {
	return c.data
}

func (c *Currency) LastUpdate() time.Time {
	return c.last
}

func (c *Currency) GetCurrencyByCode(code string) (*domain.Currency, error) {
	if currency, ok := c.data[code]; ok {
		return currency, nil
	}

	log.Printf("[services.Currency] currency not found: %s", code)

	return nil, fmt.Errorf("%s code currency not found", code)
}

func (c *Currency) parserData(data string) error {
	if c.data == nil {
		c.data = make(map[string]*domain.Currency)
	}

	currencies, err := domain.FromJSONList(data)

	if err != nil {
		log.Printf("[services.Currency] Error parsing data: %s", err)
		return err
	}

	for code, currency := range currencies {
		log.Printf("[services.Currency] currency: %s -> %s", code, currency.ToJSON())
		c.data[code] = currency
	}

	return nil
}
