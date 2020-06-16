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
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/thetreep/covidtracker"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Ensure RiskService implements covidtracker.RiskService and Accessor
var _ covidtracker.HospDAL = &HospDAL{}
var _ Accessor = &HospDAL{}

// HospDAL represents a service for managing suppliers.
type HospDAL struct {
	client     *Client
	collection *mongo.Collection
}

func (s *HospDAL) Client() *Client {
	return s.client
}

func (s *HospDAL) Collection() *mongo.Collection {
	return s.collection
}

func (s *HospDAL) Get(dep string, date time.Time) (*covidtracker.Hospitalization, error) {
	var result *covidtracker.Hospitalization

	err := s.collection.FindOne(s.client.Ctx, bson.M{"dep": dep, "date": dateFilter(date)}).Decode(&result)
	switch err {
	case mongo.ErrNoDocuments:
		return nil, fmt.Errorf("no hospitalization document found with dep=%s date=%s", dep, date.Format("2006-02-01"))
	case nil:
		return result, nil
	default:
		return nil, fmt.Errorf("error while getting hospitalization: %s", err)
	}
}

func (s *HospDAL) GetRange(dep string, start, end time.Time) ([]*covidtracker.Hospitalization, error) {
	var result []*covidtracker.Hospitalization

	cur, err := s.collection.Find(s.client.Ctx, bson.M{"dep": dep, "date": dateRangeFilter(start, end)})
	if err != nil {
		return nil, covidtracker.Errorf("error while getting case: %s", err)
	}
	defer cur.Close(s.client.Ctx)

	for cur.Next(s.client.Ctx) {
		var c *covidtracker.Hospitalization
		err := cur.Decode(&c)
		if err != nil {
			return nil, covidtracker.Errorf("error while decoding element: %s", err)
		}
		result = append(result, c)
	}
	if err := cur.Err(); err != nil {
		return nil, covidtracker.Errorf("error while reading database: %s", err)
	}

	return result, nil
}

func (s *HospDAL) Upsert(hosps ...*covidtracker.Hospitalization) error {

	if len(hosps) == 0 {
		return covidtracker.Errorf("cannot upsert empty hosps")
	}

	for _, h := range hosps {
		existing, err := s.Get(h.Department, h.Date)
		if err != nil || existing == nil { // not exist => add it
			bsonID := primitive.NewObjectID()
			h.ID = covidtracker.HospID(bsonID.Hex())
			if _, insErr := s.collection.InsertOne(s.client.Ctx, h); insErr != nil {
				return errors.Wrap(insErr, "inserting new case")
			}
		} else { // existing, update only appropriate fields
			h.ID = existing.ID
			if _, updErr := s.collection.UpdateOne(s.client.Ctx, bson.M{"_id": existing.ID}, bson.M{"$set": h}); updErr != nil {
				return errors.Wrapf(err, "updating case dep=%s date%s", h.Department, h.Date.Format("2006-02-01"))
			}
		}

	}

	return nil
}
