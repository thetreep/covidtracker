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
	risk            RiskDAL
	hotel           HotelDAL
	covCase         CaseDAL
	emergency       EmergencyDAL
	hospitalization HospDAL
	indic           IndicDAL
	screening       ScreeningDAL

	mongo    *mongo.Client
	database *mongo.Database
}

// NewClient creates a new client with mongodb scheme : mongodb://xxxx
func NewClient(mongoURI string) *Client {
	c := &Client{Now: time.Now, MongoURI: mongoURI, Ctx: context.Background()}
	c.risk.client = c
	c.hotel.client = c
	c.covCase.client = c
	c.emergency.client = c
	c.hospitalization.client = c
	c.indic.client = c
	c.screening.client = c
	return c
}

// Open opens and initializes the Mongo database.
func (c *Client) Open() error {
	opts := options.Client().ApplyURI(c.MongoURI).SetServerSelectionTimeout(10 * time.Second).SetSocketTimeout(10 * time.Second)
	if user, pwd := os.Getenv("THETREEP_COVIDTRACKER_MONGO_USER"), os.Getenv("THETREEP_COVIDTRACKER_MONGO_PASSWORD"); user != "" || pwd != "" {
		opts.SetAuth(options.Credential{
			AuthMechanism: "SCRAM-SHA-256",
			Username:      user,
			Password:      pwd,
		})
	}
	mClient, err := mongo.NewClient(opts)
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
	c.hotel.collection = c.database.Collection("hotels")
	c.covCase.collection = c.database.Collection("case")
	c.emergency.collection = c.database.Collection("emergency")
	c.hospitalization.collection = c.database.Collection("hospitalization")
	c.indic.collection = c.database.Collection("indicator")
	c.screening.collection = c.database.Collection("screening")

	return nil
}

// Close disconnect the underlying mongo database.
func (c *Client) Close() error {
	return c.mongo.Disconnect(c.Ctx)
}

// Risk returns the dal for risk
func (c *Client) Risk() covidtracker.RiskDAL { return &c.risk }

// Hotel returns the dal for hotel
func (c *Client) Hotel() covidtracker.HotelDAL { return &c.hotel }

// Case returns the dal for hospital service with at least one declared case
func (c *Client) Case() covidtracker.CaseDAL {
	return &c.covCase
}

// Emergency returns the dal for emergency data
func (c *Client) Emergency() covidtracker.EmergencyDAL {
	return &c.emergency
}

// Hospitalization returns the dal for hospitalization data
func (c *Client) Hospitalization() covidtracker.HospDAL {
	return &c.hospitalization
}

// Indicator returns the dal for indicator data
func (c *Client) Indicator() covidtracker.IndicDAL {
	return &c.indic
}

// Screening returns the dal for screening data
func (c *Client) Screening() covidtracker.ScreeningDAL {
	return &c.screening
}

type Accessor interface {
	Client() *Client
	Collection() *mongo.Collection
}
