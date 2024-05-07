package domain

import (
	"crypto/rand"
	"log"
	"time"

	"github.com/oklog/ulid"
)

type Cart struct {
	ID        string               `json:"id" bson:"_id"`
	Client    *Client              `json:"client" bson:"client"`
	CreatedAt time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time            `json:"updated_at" bson:"updated_at"`
	Items     map[string]*CartItem `json:"items" bson:"items"`
	Total     map[string]float64   `json:"currency" bson:"currency"`
	Status    CartStatus           `json:"status" bson:"status"`
}

type ClientCarts struct {
	ClientID string   `json:"client_id" bson:"_id"`
	Carts    []string `json:"carts" bson:"carts"`
}

type CartItem struct {
	Product  *Product `json:"product" bson:"product"`
	Quantity float64  `json:"quantity" bson:"quantity"`
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

func NewCart(clientId string) *Cart {
	c := &Cart{
		Items:     make(map[string]*CartItem),
		Status:    Started,
		ID:        CreateID(),
		Client:    &Client{ID: clientId},
		CreatedAt: time.Now().Local(),
		UpdatedAt: time.Now().Local(),
		Total: map[string]float64{
			"BRL": 0.0,
		},
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
}

func (c *Cart) RemoveItem(id string) {
	delete(c.Items, id)
}

func (c *Cart) UpdateItem(productId string, qty float64) {
	if _, ok := c.Items[productId]; ok {
		c.Items[productId].Quantity = qty
	}
}

func (c *Cart) AddItemQty(productId string, qty float64) {
	if _, ok := c.Items[productId]; ok {
		c.Items[productId].Quantity += qty
	}
}

func (c *Cart) GetItem(id string) *CartItem {
	if _, ok := c.Items[id]; !ok {
		return nil
	}

	return c.Items[id]
}

func (c *Cart) CalcTotal(currencies map[string]float64) {
	Total := 0.0
	for _, p := range c.Items {
		Total += p.Quantity * p.Product.Price
	}

	for k, v := range currencies {
		c.Total[k] = Total * v
		log.Printf("[domain.Cart] Total[%s]: %f -> %f", k, Total, c.Total[k])
	}
}
func CreateID() string {
	return ulid.MustNew(ulid.Now(), ulid.Monotonic(rand.Reader, 0)).String()
}
