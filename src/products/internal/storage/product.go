package storage

import (
	"errors"
	"log"

	"products/internal/domain"
)

type Product struct {
	conn *DbConnection
}

func NewProduct(conn *DbConnection) *Product {
	ret := &Product{
		conn: conn,
	}

	err := ret.Init()

	if err != nil {
		log.Printf("[storage.Product] Error initializing storage: %s", err)
		panic(err)
	}

	return ret
}

func (p *Product) Init() error {
	create := `
CREATE TABLE IF NOT EXISTS products (
	id TEXT PRIMARY KEY NOT NULL,
	name TEXT NOT NULL,
	description TEXT NOT NULL,
	price REAL NOT NULL,
	stock INTEGER NOT NULL,
	enable BOOLEAN DEFAULT TRUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	if p.conn == nil {
		log.Printf("[storage.Product] Error creating tables: db is null")
		return errors.New("db is null")
	}

	err := p.conn.Exec(create)

	if err != nil {
		log.Printf("[storage.Product] Error creating tables: %s", err)
	}

	return err
}

func (p *Product) Close() error {
	if p.conn == nil {
		log.Printf("[storage.Product] Database is already closed")
		return nil
	}

	return p.conn.Close()
}

func (p *Product) Insert(product *domain.Product) error {
	insert := `
INSERT INTO products (id, name, description, price, stock)
VALUES (?, ?, ?, ?, ?);
`

	if p.conn == nil {
		log.Printf("[storage.Product] Error inserting product: db is null")
		return errors.New("db is null")
	}

	return p.conn.Exec(insert, product.ID, product.Name, product.Description, product.Price, product.Stock)
}

func (p *Product) Get(id string) (*domain.Product, error) {
	query := `
SELECT id, name, description, price, stock, created_at, updated_at
FROM products
WHERE id = ?;
`
	if p.conn == nil {
		log.Printf("[storage.Product] Error getting product: db is null")
		return nil, errors.New("db is null")
	}

	conn, err := p.conn.GetConn()

	if err != nil {
		log.Printf("[storage.Product] Error getting connection: %s", err)
		return nil, err
	}

	if conn == nil {
		log.Printf("[storage.Product] Error getting connection")
		return nil, errors.New("connection is null")
	}

	rows, err := conn.Query(query, id)

	if err != nil {
		log.Printf("[storage.Product] Error executing query: %s -> error: %s", query, err)
		return nil, err
	}

	defer rows.Close()

	var product domain.Product

	if rows.Next() {
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CreatedAt, &product.UpdatedAt)

		if err != nil {
			log.Printf("[storage.Product] Error scanning row: %s", err)
			return nil, err
		}
	}

	return &product, nil
}

func (p *Product) GetAll() ([]*domain.Product, error) {
	query := `
SELECT id, name, description, price, stock, created_at, updated_at
FROM products
ORDER BY created_at DESC;
`

	if p.conn == nil {
		log.Printf("[storage.Product] Error getting all products: db is null")
		return nil, errors.New("db is null")
	}

	conn, err := p.conn.GetConn()

	if err != nil {
		log.Printf("[storage.Product] Error getting connection: %s", err)
		return nil, err
	}

	if conn == nil {
		log.Printf("[storage.Product] Error getting connection")
		return nil, errors.New("connection is null")
	}

	rows, err := conn.Query(query)

	if err != nil {
		log.Printf("[storage.Product] Error executing query: %s -> error: %s", query, err)
		return nil, err
	}

	defer rows.Close()

	products := make([]*domain.Product, 0)

	for rows.Next() {
		var product domain.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CreatedAt, &product.UpdatedAt)

		if err != nil {
			log.Printf("[storage.Product] Error scanning row: %s", err)
			return nil, err
		}

		products = append(products, &product)
	}

	return products, nil
}

func (p *Product) Update(product *domain.Product) error {
	update := `
UPDATE products
SET name = ?, description = ?, price = ?, stock = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;
`

	if p.conn == nil {
		log.Printf("[storage.Product] Error updating product: db is null")
		return errors.New("db is null")
	}

	return p.conn.Exec(update, product.Name, product.Description, product.Price, product.Stock, product.ID)
}

func (p *Product) Delete(id string) error {
	del := `
UPDATE products
SET enable = FALSE, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;
`

	if p.conn == nil {
		log.Printf("[storage.Product] Error deleting product: db is null")
		return errors.New("db is null")
	}

	return p.conn.Exec(del, id)
}
