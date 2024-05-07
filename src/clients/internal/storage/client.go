package storage

import (
	"errors"
	"log"

	domain "github.com/rafaelfino/plataforma-tcc/src/clients/pkg/domain"
)

type Client struct {
	conn *DbConnection
}

func NewClient(conn *DbConnection) *Client {
	ret := &Client{
		conn: conn,
	}

	err := ret.Init()

	if err != nil {
		log.Printf("[storage.Client] Error initializing storage: %s", err)
		panic(err)
	}

	return ret
}

func (c *Client) Init() error {
	create := `
CREATE TABLE IF NOT EXISTS clients (
	id TEXT PRIMARY KEY NOT NULL,
	name TEXT NOT NULL,
	surname TEXT NOT NULL,
	email TEXT DEFAULT NULL,
	birth_date DATE DEFAULT NULL,
	enable BOOLEAN DEFAULT TRUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
`

	if c.conn == nil {
		log.Printf("[storage.Client] Error creating tables: db is null")
		return errors.New("db is null")
	}

	err := c.conn.Exec(create)

	if err != nil {
		log.Printf("[storage.Client] Error creating tables: %s", err)
	}

	return err
}

func (c *Client) Close() error {
	if c.conn == nil {
		log.Printf("[storage.Client] Database is already closed")
		return nil
	}

	return c.conn.Close()
}

func (c *Client) Insert(client *domain.Client) error {
	insert := `
INSERT INTO clients (id, name, surname, email, birth_date)
VALUES (?, ?, ?, ?, ?);
`

	err := c.conn.Exec(insert, client.ID, client.Name, client.Surname, client.Email, client.BirthDate)

	if err != nil {
		log.Printf("[storage.Client] Error executing query: %s -> error: %s", insert, err)
	}

	return err
}

func (c *Client) Update(client *domain.Client) error {
	update := `
UPDATE clients
SET name = ?, surname = ?, email = ?, birth_date = ?, enable = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;
`

	err := c.conn.Exec(update, client.Name, client.Surname, client.Email, client.BirthDate, client.Enable, client.ID)

	if err != nil {
		log.Printf("[storage.Client] Error executing query: %s -> error: %s", update, err)
	}

	return err
}

func (c *Client) Delete(id string) error {
	del := `
	UPDATE clients
	SET enable = FALSE, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?;
	`

	err := c.conn.Exec(del, id)

	if err != nil {
		log.Printf("[storage.Client] Error executing query: %s -> error: %s", del, err)
	}

	return err
}

func (c *Client) Get(id string) (*domain.Client, error) {
	get := `
SELECT id, name, surname, email, birth_date, enable, created_at, updated_at
FROM clients
WHERE id = ?;
`

	conn, err := c.conn.GetConn()

	if err != nil {
		log.Printf("[storage.Client] Error getting connection: %s", err)
		return nil, err
	}

	rows, err := conn.Query(get, id)

	if err != nil {
		log.Printf("[storage.Client] Error executing query: %s -> error: %s", get, err)
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {
		client := &domain.Client{}
		err = rows.Scan(&client.ID, &client.Name, &client.Surname, &client.Email, &client.BirthDate, &client.Enable, &client.CreatedAt, &client.UpdateAt)
		if err != nil {
			log.Printf("[storage.Client] Error scanning row: %s", err)
			return nil, err
		}
		return client, nil
	}

	return nil, errors.New("client not found")
}

func (c *Client) GetAll() ([]*domain.Client, error) {
	get := `
SELECT id, name, surname, email, birth_date, enable, created_at, updated_at
FROM clients
ORDER BY CREATED_AT DESC;
`
	conn, err := c.conn.GetConn()

	if err != nil {
		log.Printf("[storage.Client] Error getting connection: %s", err)
		return nil, err
	}

	rows, err := conn.Query(get)

	if err != nil {
		log.Printf("[storage.Client] Error executing query: %s -> error: %s", get, err)
		return nil, err
	}

	defer rows.Close()

	clients := make([]*domain.Client, 0)

	for rows.Next() {
		client := &domain.Client{}
		err = rows.Scan(&client.ID, &client.Name, &client.Surname, &client.Email, &client.BirthDate, &client.Enable, &client.CreatedAt, &client.UpdateAt)
		if err != nil {
			log.Printf("[storage.Client] Error scanning row: %s", err)
			return nil, err
		}
		clients = append(clients, client)
	}

	return clients, nil
}
