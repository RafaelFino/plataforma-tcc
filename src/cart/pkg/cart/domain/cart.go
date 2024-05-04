package domain

import (
	"crypto/rand"
	"log"
	"time"

	"products/pkg/currencies/domain"
	"products/pkg/products/domain"

	"github.com/oklog/ulid"
)

type Cart struct {
	ID            string                        `json:"id"`
	ClientID      string                        `json:"client_id"`
	ClientName    string                        `json:"client_name"`
	ClientSurname string                        `json:"client_surname"`
	ClientEmail   string                        `json:"client_email"`
	CreatedAt     time.Time                     `json:"created_at"`
	UpdatedAt     time.Time                     `json:"updated_at"`
	Products      map[string]*domain.Product    `json:"products"`
	Currencies    map[string]*domain.Currencies `json:"currencies"`
	Status        CartStatus                    `json:"status"`
}

type CartStatus int

const (
	Started CartStatus = iota
	Progress
	Done
	Canceled
)

func (c CartStatus) String() string {
	return [...]string{"STARTED", "PROGRESS", "DONE", "CANCELED"}[c]
}

func (c CartStatus) EnumIndex() int {
	return int(c)
}

func NewCart(clientId string, currencies map[string]float64) *Cart {
	c := &Cart{
		Products:   make(map[string]*domain.Product),
		Status:     Started,
		ID:         CreateID(),
		ClientID:   clientId,
		Currencies: currencies,
	}

	return c
}

func (c *Cart) GetID() string {
	return c.ID
}

func (c *Cart) AddProduct(p *domain.Product) {
	c.Products[p.ID] = p
}

func (c *Cart) RemoveProduct(id string) {
	delete(c.Products, id)
}

func (c *Cart) UpdateProduct(p *domain.Product) {
	c.Products[p.ID] = p
}

func (c *Cart) GetProduct(id string) *domain.Product {
	return c.Products[id]
}

func (c *Cart) GetProducts() map[string]*domain.Product {
	return c.Products
}

func (c *Cart) CalcTotal() float64 {
	var total float64
	for _, p := range c.Products {
		if curr, ok := c.Currencies[p.Currency]; ok {
			total += p.Price * float64(p.Quantity) * curr
		} else {
			log.Printf("[domain.Cart] Currency not found: %s, using BRL", p.Currency)
			total += p.Price * float64(p.Quantity)
		}
	}
	return total

}

func (c *Cart) AddClientName(name string) {
	c.ClientName = name
}
func CreateID() string {
	return ulid.MustNew(ulid.Now(), ulid.Monotonic(rand.Reader, 0)).String()
}
