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
	ID        string               `json:"id"`
	Client    *domain.Client       `json:"client"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
	Items     map[string]*CartItem `json:"items"`
	Total     float64              `json:"total"`
	Status    CartStatus           `json:"status"`
}

type CartItem struct {
	Product  *domain.Product `json:"product"`
	Quantity float64         `json:"quantity"`
	Currency string          `json:"currency"`
	Quote    float64         `json:"quote"`
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

func (c *CartItem) GetTotal() float64 {
	return c.Product.Price * float64(c.Quantity) * c.Quote
}

func NewCart(clientId string) *Cart {
	c := &Cart{
		Items:    make(map[string]*CartItem),
		Status:   Started,
		ID:       CreateID(),
		ClientID: clientId,
		Total:    0,
	}

	return c
}

func (c *Cart) GetID() string {
	return c.ID
}

func (c *Cart) AddItem(p *CartItem) {
	if _, ok := c.Items[p.Product.ID]; ok {
		c.Items[p.Product.ID].Quantity += p.Quantity
	} else {
		c.Items[p.Product.ID] = p
	}

	c.CalcTotal()
}

func (c *Cart) RemoveItem(id string) {
	delete(c.Items, id)
	c.CalcTotal()
}

func (c *Cart) UpdateItem(productId string, qty float64) {
	if _, ok := c.Items[productId]; ok {
		c.Items[productId].Quantity = qty
		c.CalcTotal()
	}

}

func (c *Cart) AddItemQty(productId string, qty float64) {
	if _, ok := c.Items[productId]; ok {
		c.Items[productId].Quantity += qty
		c.CalcTotal()
	}
}

func (c *Cart) GetItem(id string) *CartItem {
	if _, ok := c.Items[id]; !ok {
		return nil
	}

	return c.Items[id]
}

func (c *Cart) CalcTotal() {
	c.Total = 0
	for _, p := range c.Items {
		c.Total += p.GetTotal()
	}

	log.Printf("[domain.Cart] Total: %f", c.Total)
}
func CreateID() string {
	return ulid.MustNew(ulid.Now(), ulid.Monotonic(rand.Reader, 0)).String()
}
