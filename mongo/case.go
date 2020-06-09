package mongo

import (
	"fmt"
	"time"

	"github.com/thetreep/covidtracker"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Ensure RiskService implements covidtracker.RiskService and Accessor
var _ covidtracker.CaseDAL = &CaseDAL{}
var _ Accessor = &CaseDAL{}

// CaseDAL represents a service for managing suppliers.
type CaseDAL struct {
	client     *Client
	collection *mongo.Collection
}

func (s *CaseDAL) Client() *Client {
	return s.client
}

func (s *CaseDAL) Collection() *mongo.Collection {
	return s.collection
}

func (s *CaseDAL) Get(dep int, date time.Time) ([]*covidtracker.Case, error) {
	var result []*covidtracker.Case
	if err := s.findCases(dateFilter(date), result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *CaseDAL) GetRange(dep int, start, end time.Time) ([]*covidtracker.Case, error) {
	var result []*covidtracker.Case
	if err := s.findCases(dateRangeFilter(start, end), result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *CaseDAL) Upsert(cases ...*covidtracker.Case) error {

	if len(cases) == 0 {
		return covidtracker.Errorf("cannot upsert empty cases")
	}

	for _, c := range cases {

		existing, err := s.Get(c.Department, c.NoticeDate)
		if err != nil || existing == nil { // unexisting, add it
			bsonID := primitive.NewObjectID()
			c.ID = covidtracker.CaseID(bsonID.Hex())
			if _, insErr := s.collection.InsertOne(s.client.Ctx, c); insErr != nil {
				return fmt.Errorf("upsert: error while inserting new case: %s", insErr)
			}
		} else { // existing, update only appropriate fields
			c.ID = existing.ID // do not update these fields
			sub.Number = existing.Number
			sub.CurrentCount = existing.CurrentCount
			sub.CreationDate = existing.CreationDate
			sub.UpdateDate = now
			if _, updErr := s.collection.UpdateOne(s.client.Ctx, bson.M{"_id": existing.ID}, bson.M{"$set": sub}); updErr != nil {
				return fmt.Errorf("upsert: error while updating subsidy with number %q: %s", sub.Number, updErr)
			}
		}

	}

	return nil
}

func (s *CaseDAL) findCases(filter bson.D, cases []*covidtracker.Case) error {
	cur, err := s.collection.Find(s.client.Ctx, filter, nil)
	if err != nil {
		return covidtracker.Errorf("error while getting case: %s", err)
	}
	defer cur.Close(s.client.Ctx)

	for cur.Next(s.client.Ctx) {
		var c *covidtracker.Case
		err := cur.Decode(&c)
		if err != nil {
			return covidtracker.Errorf("error while decoding element: %s", err)
		}
		cases = append(cases, c)
	}

	if err := cur.Err(); err != nil {
		return covidtracker.Errorf("error while reading database: %s", err)
	}

	return nil
}
