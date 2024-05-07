package domain

import "time"

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

func NewClient() *Client {
	return &Client{}
}

func (c *Client) GetID() string {
	return c.ID
}
