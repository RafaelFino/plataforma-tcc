package services

import (
	"clients/internal/config"
	"clients/internal/storage"
	"clients/pkg/clients/domain"
	"crypto/rand"
	"errors"
	"log"

	"github.com/oklog/ulid"
)

type Client struct {
	storage *storage.Client
}

func NewClient(config *config.Config) *Client {
	return &Client{
		storage: storage.NewClient(storage.NewDbConnection(config.DBPath)),
	}
}

func (c *Client) GetById(id string) (*domain.Client, error) {
	log.Printf("[services.Client] Getting client with ID: %s", id)
	ret, err := c.storage.Get(id)

	if err != nil {
		log.Printf("[services.Client] Error getting client: %s", err)
		return nil, err
	}

	return ret, nil
}

func (c *Client) Get() ([]*domain.Client, error) {
	log.Printf("[services.Client] Getting all clients")
	return c.storage.GetAll()
}

func (c *Client) Update(client *domain.Client) error {
	if client == nil {
		log.Printf("[services.Client] client is nil")
		return errors.New("client is nil")
	}

	if client.ID == "" {
		log.Printf("[services.Client] client ID is empty")
		return errors.New("client ID is empty")
	}

	log.Printf("[services.Client] Updating client with ID: %s to %+v", client.ID, client)
	return c.storage.Update(client)
}

func (c *Client) Insert(client *domain.Client) (string, error) {
	if client == nil {
		log.Printf("[services.Client] client is nil")
		return "", errors.New("client is nil")
	}

	client.ID = c.CreateID()

	log.Printf("[services.Client] Inserting client with ID: %s -> %+v", client.ID, client)
	return client.ID, c.storage.Insert(client)
}

func (c *Client) Delete(id string) error {
	log.Printf("[services.Client] Deleting client with ID: %s", id)
	return c.storage.Delete(id)
}

func (c *Client) Close() error {
	log.Printf("[services.Client] Closing storage")
	return c.storage.Close()
}

func (c *Client) CreateID() string {
	return ulid.MustNew(ulid.Now(), ulid.Monotonic(rand.Reader, 0)).String()
}
