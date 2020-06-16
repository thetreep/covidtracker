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

var _ covidtracker.HotelDAL = &HotelDAL{}
var _ Accessor = &HotelDAL{}

type HotelDAL struct {
	client     *Client
	collection *mongo.Collection
}

func (s *HotelDAL) Client() *Client {
	return s.client
}

func (s *HotelDAL) Collection() *mongo.Collection {
	return s.collection
}

// Get returns a hotel by ID.
func (s *HotelDAL) Get(id covidtracker.HotelID) (*covidtracker.Hotel, error) {
	var result *covidtracker.Hotel
	if err := s.collection.FindOne(s.client.Ctx, bson.M{"_id": string(id)}).Decode(&result); err != nil {
		return nil, covidtracker.Errorf("error while getting hotel: %s", err)
	}
	return result, nil
}

//Insert CdsApi result
func (s *HotelDAL) Insert(hotels []*covidtracker.Hotel) ([]*covidtracker.Hotel, error) {
	if hotels == nil || len(hotels) == 0 {
		return nil, covidtracker.ErrDocRequired("hotels")
	}
	var (
		err     error
		session mongo.Session
		ctx     = context.Background()
	)

	if session, err = s.client.mongo.StartSession(); err != nil {
		return nil, covidtracker.Errorf("error while mongo start session to insert %d hotels: %s", len(hotels), err)
	}
	if err = session.StartTransaction(); err != nil {
		return nil, covidtracker.Errorf("error while mongo start transaction to insert %d hotels: %s", len(hotels), err)
	}

	var result []*covidtracker.Hotel
	if err = mongo.WithSession(ctx, session, func(sessCtx mongo.SessionContext) error {
		for _, hotel := range hotels {
			var resp *covidtracker.Hotel
			if notFound := s.collection.FindOne(s.client.Ctx, bson.M{"name": hotel.Name, "address": hotel.Address}).Decode(&resp); notFound == nil {
				result = append(result, resp)
				continue
			}
			bsonID := primitive.NewObjectID()
			hotel.ID = covidtracker.HotelID(bsonID.Hex())
			_, err := s.collection.InsertOne(s.client.Ctx, hotel)
			if err != nil {
				session.AbortTransaction(ctx)
				return covidtracker.Errorf("error while inserting hotel: %s", err)
			}
			result = append(result, hotel)
		}

		if err := session.CommitTransaction(sessCtx); err != nil {
			return covidtracker.Errorf("error while mongo commit transaction to insert %d hotels: %s", len(hotels), err)
		}
		return nil
	}); err != nil {
		return nil, covidtracker.Errorf("error while mongo session to insert %d hotels: %s", len(hotels), err)
	}
	session.EndSession(ctx)
	return result, nil
}
