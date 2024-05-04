package storage

import (
	"errors"
	"log"

	"cart/internal/domain"
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

func (c *Cart) CreateCart(cart *domain.Cart) error {
	if cart == nil {
		log.Printf("[storage.Cart] cart is nil")
		return errors.New("cart is nil")
	}

	if cart.ID == "" {
		log.Printf("[storage.Cart] cart ID is empty")
		return errors.New("cart ID is empty")
	}

	log.Printf("[storage.Cart] Inserting cart with ID: %s -> %+v", cart.ID, cart)

	create := `
	INSERT INTO carts (id, client_id, status_id) VALUES (?, ?, ?);
	INSERT INTO cart_status (id, status) VALUES (?, ?);
`

	err := c.conn.Exec(create, cart.ID, cart.ClientID, cart.Status, cart.ID, cart.Status)

	if err != nil {
		log.Printf("[storage.Cart] Error inserting cart: %s", err)
	}

	return err
}

func (c *Cart) SetCartStatus(cartId string, newStatus domain.CartStatus) error {
	if id == "" {
		log.Printf("[storage.Cart] cart ID is empty")
		return errors.New("cart ID is empty")
	}

	log.Printf("[storage.Cart] Updating cart status with ID: %s to %+v", cartId, newStatus)

	update := `	
	UPDATE carts SET status_id = ? WHERE id = ?;
	INSERT INTO cart_status (id, status) VALUES (?, ?);

`
	err := c.conn.Exec(update, newStatus, cartId, cartId, newStatus)

	if err != nil {
		log.Printf("[storage.Cart] Error updating cart status: %s", err)
	}

	return err
}

func (c *Cart) AddCartItem(cartId string, item *domain.Product) error {
	insert := `
	INSERT INTO cart_items (cart_id, product_id, quantity, price) VALUES (?, ?, ?, ?);
`

	err := c.conn.Exec(insert, cartId, item.ID, item.Quantity, item.Price)

	if err != nil {
		log.Printf("[storage.Cart] Error inserting cart item: %s", err)
	}

	return err
}

func (c *Cart) RemoveCartItem(cartId string, productId string) error {
	del := `
	DELETE FROM cart_items WHERE cart_id = ? AND product_id = ?;
`

	err := c.conn.Exec(del, cartId, productId)

	if err != nil {
		log.Printf("[storage.Cart] Error deleting cart item: %s", err)
	}

	return err
}

func (c *Cart) GetCart(cartId string) (*domain.Cart, error) {
	query := `
	SELECT 
		id, 
		client_id, 
		status, 
		created_at,
		updated_at
	FROM carts
	WHERE id = ?;
	`

	conn, err := c.conn.GetConn()

	if err != nil {
		log.Printf("[storage.Cart] Error getting connection: %s", err)
		return nil, err
	}

	rows, err := conn.Query(query, cartId)

	if err != nil {
		log.Printf("[storage.Cart] Error executing query: %s -> error: %s", query, err)
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		log.Printf("[storage.Cart] No rows returned")
		return nil, errors.New("no rows returned")
	}

	ret := &domain.Cart{}

	err = rows.Scan(&ret.ID, &ret.ClientID, &ret.Status, &ret.CreatedAt, &ret.UpdatedAt)

	if err != nil {
		log.Printf("[storage.Cart] Error scanning row: %s", err)
		return nil, err
	}

	err = c.GetCartItems(ret)

	if err != nil {
		log.Printf("[storage.Cart] Error getting cart items: %s", err)
	}

	return ret, err
}

func (c *Cart) GetCartItems(cart *domain.Cart) error {
	if cart == nil {
		log.Printf("[storage.Cart] cart is nil")
		return errors.New("cart is nil")
	}

	query := `
	SELECT 
		product_id, 
		quantity, 
		price
	FROM cart_items
	WHERE cart_id = ?;
	`

	conn, err := c.conn.GetConn()

	if err != nil {
		log.Printf("[storage.Cart] Error getting connection: %s", err)
		return err
	}

	rows, err := conn.Query(query, cart.ID)

	if err != nil {
		log.Printf("[storage.Cart] Error executing query: %s -> error: %s", query, err)
		return err
	}

	defer rows.Close()

	if cart.Products == nil {
		cart.Products = make(map[string]*domain.Product, 0)
	}

	for rows.Next() {
		item := &domain.Product{}
		err = rows.Scan(&item.ID, &item.Quantity, &item.Price)

		if err != nil {
			log.Printf("[storage.Cart] Error scanning row: %s", err)
			return err
		}

		cart.Products[item.ID] = item
	}

	return nil
}

func (c *Cart) GetClientCarts(clientId string) ([]*domain.Cart, error) {
	query := `
	SELECT 
		id
	FROM carts
	WHERE client_id = ?;
	`

	conn, err := c.conn.GetConn()

	if err != nil {
		log.Printf("[storage.Cart] Error getting connection: %s", err)
		return nil, err
	}

	rows, err := conn.Query(query, clientId)

	if err != nil {
		log.Printf("[storage.Cart] Error executing query: %s -> error: %s", query, err)
		return nil, err
	}

	defer rows.Close()

	ret := make([]*domain.Cart, 0)

	for rows.Next() {
		var id string

		err = rows.Scan(&id)
		if err != nil {
			log.Printf("[storage.Cart] Error scanning row: %s", err)
			return nil, err
		}

		cart, err := c.GetCart(id)

		if err != nil {
			log.Printf("[storage.Cart] Error getting cart items: %s", err)
			return nil, err
		}

		ret = append(ret, cart)
	}

	return ret, nil
}
