package domain

import (
	"crypto/rand"
	"log"
	"time"

	"github.com/oklog/ulid"
	"golang.org/x/vuln/client"
)

type Client struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Email     string    `json:"email"`
	BirthDate string    `json:"birth_date"`
	Enable    bool      `json:"enable"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"updated_at"`
}

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Currency struct {
	Code       string  `json:"code"`
	Codein     string  `json:"codein"`
	Name       string  `json:"name"`
	High       float64 `json:"high,omitempty"`
	Low        float64 `json:"low,omitempty"`
	VarBid     float64 `json:"varBid,omitempty"`
	PctChange  float64 `json:"pctChange,omitempty"`
	Bid        float64 `json:"bid,omitempty"`
	Ask        float64 `json:"ask,omitempty"`
	Timestamp  string  `json:"timestamp,omitempty"`
	CreateDate string  `json:"create_date,omitempty"`
}

type Cart struct {
	ID        string               `json:"id"`
	Client    *Client              `json:"client"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
	Items     map[string]*CartItem `json:"items"`
	Total     float64              `json:"total"`
	Status    CartStatus           `json:"status"`
}

type CartItem struct {
	Product  *Product `json:"product"`
	Quantity float64  `json:"quantity"`
	Currency string   `json:"currency"`
	Quote    float64  `json:"quote"`
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
		Items:  make(map[string]*CartItem),
		Status: Started,
		ID:     CreateID(),
		Client: &client.Client{ID: clientId},
		Total:  0,
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
