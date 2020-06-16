/*
	This file is part of covidtracker also known as EviteCovid .

    Copyright 2020 the Treep

    covdtracker is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    covidtracker is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with covidtracker.  If not, see <https://www.gnu.org/licenses/>.
*/

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
