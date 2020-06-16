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
	"context"

	"github.com/thetreep/covidtracker"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Ensure RiskService implements covidtracker.RiskService and Accessor
var _ covidtracker.RiskDAL = &RiskDAL{}
var _ Accessor = &RiskDAL{}

// RiskService represents a service for managing suppliers.
type RiskDAL struct {
	client     *Client
	collection *mongo.Collection
}

func (s *RiskDAL) Client() *Client {
	return s.client
}

func (s *RiskDAL) Collection() *mongo.Collection {
	return s.collection
}

// Get returns a risk by ID.
func (s *RiskDAL) Get(id covidtracker.RiskID) (*covidtracker.Risk, error) {
	var result *covidtracker.Risk
	if err := s.collection.FindOne(s.client.Ctx, bson.M{"_id": string(id)}).Decode(&result); err != nil {
		return nil, covidtracker.Errorf("error while getting risk: %s", err)
	}
	return result, nil
}

// Create creates new risks.
func (s *RiskDAL) Insert(risks ...*covidtracker.Risk) error {
	if risks == nil || len(risks) == 0 {
		return covidtracker.ErrDocRequired("risk")
	}

	var (
		err     error
		session mongo.Session
		ctx     = context.Background()
	)

	if session, err = s.client.mongo.StartSession(); err != nil {
		return covidtracker.Errorf("error while mongo start session to create %d risks: %s", len(risks), err)
	}
	if err = session.StartTransaction(); err != nil {
		return covidtracker.Errorf("error while mongo start transaction to create %d risks: %s", len(risks), err)
	}

	//@todo see s.collection.InsertMany(...)
	if err = mongo.WithSession(ctx, session, func(sessCtx mongo.SessionContext) error {

		for _, risk := range risks {
			bsonID := primitive.NewObjectID()
			risk.ID = covidtracker.RiskID(bsonID.Hex())
			_, err := s.collection.InsertOne(s.client.Ctx, risk)
			if err != nil {
				session.AbortTransaction(ctx)
				return covidtracker.Errorf("error while inserting operation: %s", err)
			}
		}

		if err := session.CommitTransaction(sessCtx); err != nil {
			return covidtracker.Errorf("error while mongo commit transaction to create %d risk: %s", len(risks), err)
		}

		return nil
	}); err != nil {
		return covidtracker.Errorf("error while mongo session to create %d risk: %s", len(risks), err)
	}
	session.EndSession(ctx)
	return nil
}
