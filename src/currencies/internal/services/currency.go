package services

import (
	"encoding/json"
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
	c := &Currency{
		url:      url,
		interval: 60,
		last:     time.Now().Add(-time.Duration(24) * time.Hour),
	}

	err := c.Update()

	if err != nil {
		log.Printf("[services.Currency] Error updating currency: %s", err)
	}

	return c
}

func (c *Currency) Update() error {
	log.Printf("[services.Currency] Updating currency")

	if time.Now().Before(c.last.Add(time.Duration(c.interval) * time.Minute)) {
		log.Printf("[services.Currency] Minimal interval not reached, wait for %s", c.last.Add(time.Duration(c.interval)*time.Minute).Format(time.RFC3339))
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

func (c *Currency) Get() map[string]*domain.Currency {
	return c.data
}

func (c *Currency) LastUpdate() time.Time {
	return c.last
}

func (c *Currency) GetByCode(code string) (*domain.Currency, error) {
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

	jsonData := make(map[string]interface{})

	err := json.Unmarshal([]byte(data), &jsonData)

	if err != nil {
		log.Printf("[services.Currency] Error parsing json data to map: %s", err)
		return err
	}

	for code, currency := range jsonData {
		log.Printf("[services.Currency] currency: %s -> %+v", code, currency)
		c.data[code] = domain.CurrencyFromMap(currency.(map[string]interface{}))
	}

	return nil
}
