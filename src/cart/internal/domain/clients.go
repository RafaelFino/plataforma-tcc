package domain

import "time"

type Client struct {
	ID        string    `json:"id" bson:"_id"`
	Name      string    `json:"name" bson:"name"`
	Surname   string    `json:"surname" bson:"surname"`
	Email     string    `json:"email" bson:"email"`
	BirthDate string    `json:"birth_date" bson:"birth_date"`
	Enable    bool      `json:"enable" bson:"enable"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdateAt  time.Time `json:"updated_at" bson:"updated_at"`
}
