package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	client   *mongo.Client
	database *mongo.Database
}

// NewMongoDBClient creates a new MongoDB client
func NewMongoDBClient(uri, dbName string, maxPoolSize, minPoolSize int) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(uint64(maxPoolSize)).
		SetMinPoolSize(uint64(minPoolSize))

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return &Client{
		client:   client,
		database: client.Database(dbName),
	}, nil
}

// Close closes the MongoDB connection
func (c *Client) Close(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}

// GetDatabase returns the database instance
func (c *Client) GetDatabase() *mongo.Database {
	return c.database
}

// GetCollection returns a collection instance
func (c *Client) GetCollection(name string) *mongo.Collection {
	return c.database.Collection(name)
}
