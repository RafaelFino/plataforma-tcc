package storage

import (
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

}

func (c *Cart) Close() error {
	if c.conn == nil {
		log.Printf("[storage.Cart] Database is already closed")
		return nil
	}

	return c.conn.Close()
}
