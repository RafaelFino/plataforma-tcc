package storage

import (
	"cart/internal/config"
	"cart/internal/domain"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Cart struct {
	conn *DbConnection
}

func NewCart(config *config.Config) *Cart {
	ret := &Cart{
		conn: NewDbConnection(config),
	}

	return ret
}

func (c *Cart) Close() error {
	if c.conn == nil {
		log.Printf("[storage.Cart] Database is already closed")
		return nil
	}

	return c.conn.Close()
}

func (c *Cart) CreateCart(cart *domain.Cart) error {
	carts, err := c.conn.GetCollection("carts")

	if err != nil {
		log.Printf("[storage.Cart] Error getting collection: %s collection name: %s", err, "carts")
		return err
	}

	insertResult, err := carts.InsertOne(context.Background(), cart)

	if err != mongo.ErrNilCursor {
		return err
	}

	if err != nil {
		log.Printf("[storage.Cart] Error creating cart: %s", err)
		return err
	}

	log.Printf("[storage.Cart] Cart created: %s", insertResult.InsertedID)

	clientCarts, err := c.conn.GetCollection("client_carts")

	if err != nil {
		log.Printf("[storage.Cart] Error getting collection: %s collection name: %s", err, "client_carts")
		return err
	}

	match := bson.M{"_id": cart.Client.ID}
	change := bson.M{"$push": bson.M{"carts": bson.M{"$each": cart.ID}}}

	_, err = clientCarts.UpdateOne(context.Background(), match, change)

	if err != nil {
		log.Printf("[storage.Cart] Error updating client_carts: %s clientID: %s cartID: %s", err, cart.Client.ID, cart.ID)
	}

	return err
}

func (c *Cart) Get(cartId string) (*domain.Cart, error) {
	collection, err := c.conn.GetCollection("carts")

	if err != nil {
		return nil, err
	}

	var cart domain.Cart

	err = collection.FindOne(context.Background(), domain.Cart{ID: cartId}).Decode(&cart)

	if err != nil {
		log.Printf("[storage.Cart] Error getting cart: %s", err)
		return nil, err
	}

	log.Printf("[storage.Cart] Cart found: %s", cart.ID)

	return &cart, nil
}

func (c *Cart) UpdateCart(cart *domain.Cart) error {
	collection, err := c.conn.GetCollection("carts")

	if err != nil {
		log.Printf("[storage.Cart] Error getting collection: %s", err)
		return err
	}

	_, err = collection.ReplaceOne(context.Background(), domain.Cart{ID: cart.ID}, cart)

	if err != nil {
		log.Printf("[storage.Cart] Error updating cart: %s", err)
		return err
	}

	log.Printf("[storage.Cart] Cart updated: %s -> %+v", cart.ID, cart)
	return nil
}

func (c *Cart) GetByClient(clientId string) ([]*domain.Cart, error) {
	collection, err := c.conn.GetCollection("client_carts")

	if err != nil {
		log.Printf("[storage.Cart] Error getting collection: %s (client_carts)", err)
		return nil, err
	}

	cursor, err := collection.Find(context.Background(), bson.M{"_id": clientId})

	if err != nil {
		log.Printf("[storage.Cart] Error getting carts: %s", err)
		return nil, err
	}

	cartIds := []string{}

	for cursor.Next(context.Background()) {
		var result bson.M

		err := cursor.Decode(&result)

		if err != nil {
			log.Printf("[storage.Cart] Error decoding cart: %s", err)
			return nil, err
		}

		cartIds = append(cartIds, result["carts"].([]string)...)
	}

	cursor, err = c.conn.GetCollection("carts")

	if err != nil {
		log.Printf("[storage.Cart] Error getting collection: %s (carts)", err)
		return nil, err
	}

	cursor, err = collection.Find(context.Background(), bson.M{"_id": bson.M{"$in": cartIds}})

	carts := []*domain.Cart{}
	
	for cursor.Next(context.Background()) {
		var cart domain.Cart

		err = cursor.Decode(&cart)
		
		if err != nil {
			log.Printf("[storage.Cart] Error parsing cart: %s", err)
			return nil, err
		}

		carts = append(carts, cart)
	}

	log.Printf("[storage.Cart] Carts found: %d", len(carts))

	return carts, nil
}
