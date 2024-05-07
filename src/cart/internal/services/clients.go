package services

import (
	"cart/internal/config"
	"clients/internal/domain"
	"errors"
	"log"
	"net/http"
	"time"

	"encoding/json"
)

type Clients struct {
	url string
}

type ClientData struct {
	Client    *domain.Client `json:"client"`
	Elapsed   float64        `json:"elapsed"`
	Timestamp time.Time      `json:"timestamp"`
}

func NewClients(config *config.Config) *Clients {
	return &Clients{
		url: config.ClientsURL,
	}
}
func (c *Clients) Get(id string) (*domain.Client, error) {
	log.Printf("[services.Clients] Get client data for: %s", id)

	body, status, err := HttpGet(c.url)

	if err != nil {
		log.Printf("[services.Clients] Error get client details: %s", err)
		return nil, err
	}

	if status != http.StatusOK {
		log.Printf("[services.Clients] Error get client details: %d", status)
		return nil, err
	}

	data := &ClientData{}

	err = json.Unmarshal([]byte(body), data)

	if err != nil {
		log.Printf("[services.Clients] Error parsing client details: %s", err)
		return nil, err
	}

	if data.Client == nil {
		log.Printf("[services.Clients] Client is nil")
		return errors.New("client is nil")
	}

	log.Printf("[services.Clients] Client details: %+v", data.Client)

	return data.Client, nil
}
