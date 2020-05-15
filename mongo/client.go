package mongo

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/thetreep/covidtracker"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Client encapsulates the mongo client
type Client struct {
	MongoURI string
	// Returns the current time.
	Now func() time.Time
	Ctx context.Context

	// DAL
	risk RiskDAL

	mongo    *mongo.Client
	database *mongo.Database
}

// NewClient creates a new client with mongodb scheme : mongodb://xxxx
func NewClient(mongoURI string) *Client {
	c := &Client{Now: time.Now, MongoURI: mongoURI, Ctx: context.Background()}
	c.risk.client = c
	return c
}

// Open opens and initializes the Mongo database.
func (c *Client) Open() error {
	mClient, err := mongo.NewClient(options.Client().ApplyURI(c.MongoURI))
	if err != nil {
		return fmt.Errorf("error while creating mongo client: %s", err)
	}
	ctx, cancel := context.WithTimeout(c.Ctx, 10*time.Second)
	defer cancel()
	if err = mClient.Connect(ctx); err != nil {
		return fmt.Errorf("error while connecting mongo client: %s", err)
	}
	if err = mClient.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("error while pinging mongo server: %s", err)
	}
	// Mongo client is up and server is reachable
	c.mongo = mClient

	mongoDatabase := os.Getenv("THETREEP_COVIDTRACKER_DATABASE")
	if mongoDatabase == "" {
		mongoDatabase = "thetreep-covidtracker"
	}

	c.database = c.mongo.Database(mongoDatabase)
	c.risk.collection = c.database.Collection("risk")

	return nil
}

// Close disconnect the underlying mongo database.
func (c *Client) Close() error {
	return c.mongo.Disconnect(c.Ctx)
}

// Risk returns the dal for risk
func (c *Client) Risk() covidtracker.RiskDAL { return &c.risk }

type Accessor interface {
	Client() *Client
	Collection() *mongo.Collection
}
