package storage

import (
	"errors"
	"log"
)

type Cart struct {
	conn *DbConnection
}

func NewCart(conn *DbConnection) *Cart {
	ret := &Cart{
		conn: conn,
	}

	err := ret.Init()

	if err != nil {
		log.Printf("[storage.Cart] Error initializing storage: %s", err)
		panic(err)
	}

	return ret
}

func (c *Cart) Init() error {
	create := `
CREATE TABLE IF NOT EXISTS carts (
	id TEXT PRIMARY KEY NOT NULL,
	client_id TEXT NOT NULL,
	status_id INTEGER NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

CREATE TABLE IF NOT EXISTS cart_items (
	cart_id TEXT NOT NULL,
	product_id TEXT NOT NULL,
	quantity INTEGER NOT NULL,
	price REAL NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (cart_id, product_id)
	);

CREATE TABLE IF NOT EXISTS cart_status (
	id TEXT PRIMARY KEY NOT NULL,	
	STATUS TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
`

	if c.conn == nil {
		log.Printf("[storage.Cart] Error creating tables: db is null")
		return errors.New("db is null")
	}

	err := c.conn.Exec(create)

	if err != nil {
		log.Printf("[storage.Cart] Error creating tables: %s", err)
	}

	return err
}

func (c *Cart) Close() error {
	if c.conn == nil {
		log.Printf("[storage.Cart] Database is already closed")
		return nil
	}

	return c.conn.Close()
}
