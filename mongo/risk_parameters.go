package mongo

import (
	"github.com/thetreep/covidtracker"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Ensure ParametersService implements covidtracker.ParametersService and Accessor
var _ covidtracker.RiskParametersDAL = &RiskParametersDAL{}
var _ Accessor = &RiskParametersDAL{}

type RiskParametersDAL struct {
	client     *Client
	collection *mongo.Collection
}

func (s *RiskParametersDAL) Client() *Client {
	return s.client
}

func (s *RiskParametersDAL) Collection() *mongo.Collection {
	return s.collection
}

// GetDefault returns the default parameters from db
func (s *RiskParametersDAL) GetDefault() (*covidtracker.RiskParameters, error) {
	var result *covidtracker.RiskParameters
	if err := s.collection.FindOne(s.client.Ctx, bson.M{"default": true}).Decode(&result); err != nil && err != mongo.ErrNoDocuments {
		return nil, covidtracker.Errorf("error while getting parameters: %s", err)
	} else if err == mongo.ErrNoDocuments {
		return nil, covidtracker.ErrNoParametersDefined
	}
	return result, nil
}

// Insert creates new parameters.
func (s *RiskParametersDAL) Insert(params *covidtracker.RiskParameters) error {
	if params == nil {
		return covidtracker.ErrDocRequired("parameters")
	}

	var (
		err     error
		session mongo.Session
	)

	if session, err = s.client.mongo.StartSession(); err != nil {
		return covidtracker.Errorf("error while mongo start session to create parameters: %s", err)
	}
	if err = session.StartTransaction(); err != nil {
		return covidtracker.Errorf("error while mongo start transaction to create parameters: %s", err)
	}
	_, err = s.collection.InsertOne(s.client.Ctx, params)
	if err != nil {
		return covidtracker.Errorf("error while inserting parameters: %s", err)
	}

	return nil
}
