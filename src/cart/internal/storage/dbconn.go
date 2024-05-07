package storage

import (
	"cart/internal/config"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DbConnection struct {
	client        *mongo.Client
	clientOptions *options.ClientOptions
	url           string
	DBName        string
}

func NewDbConnection(config *config.Config) *DbConnection {
	return &DbConnection{
		url:    config.DBURL,
		DBName: config.DBName,
	}
}

func (d *DbConnection) Connect() error {
	if d.client != nil {
		return nil
	}

	log.Printf("[DbConnection] Connecting to %s", d.url)

	clientOptions := options.Client().ApplyURI(d.url)

	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Printf("[DbConnection] Error connecting to database: %s", err)
		return err
	}

	d.client = client
	d.clientOptions = clientOptions

	return nil
}

func (d *DbConnection) Close() error {
	if d.client == nil {
		log.Printf("[DbConnection] Database is already closed")
		return nil
	}

	err := d.client.Disconnect(context.Background())

	if err != nil {
		log.Printf("[DbConnection] Error closing connection: %s", err)
		return err
	}

	d.client = nil
	d.clientOptions = nil

	log.Printf("[DbConnection] Connection closed for %s", d.url)

	return nil
}

func (d *DbConnection) GetCollection(name string) (*mongo.Collection, error) {
	if d.client == nil {
		log.Printf("[DbConnection] Database is not connected")

		err := d.Connect()

		if err != nil {
			log.Printf("[DbConnection] Error connecting to database: %s", err)
			return nil, err
		}
	}

	return d.client.Database(d.DBName).Collection(name), nil
}
