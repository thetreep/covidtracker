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

func (s *HospDAL) Get(dep int, date time.Time) (*covidtracker.Hospitalization, error) {
	var result *covidtracker.Hospitalization

	err := s.collection.FindOne(s.client.Ctx, bson.M{"dep": dep, "date": dateFilter(date)}).Decode(&result)
	switch err {
	case mongo.ErrNoDocuments:
		return nil, fmt.Errorf("no hospitalization document found with dep=%d date=%s", dep, date.Format("2006-02-01"))
	case nil:
		return result, nil
	default:
		return nil, fmt.Errorf("error while getting hospitalization: %s", err)
	}

	return result, nil
}

func (s *HospDAL) GetRange(dep int, start, end time.Time) ([]*covidtracker.Hospitalization, error) {
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
				return errors.Wrapf(err, "updating case dep=%d date%s", h.Department, h.Date.Format("2006-02-01"))
			}
		}

	}

	return nil
}
