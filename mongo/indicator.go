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
var _ covidtracker.IndicDAL = &IndicDAL{}
var _ Accessor = &IndicDAL{}

// IndicDAL represents a service for managing suppliers.
type IndicDAL struct {
	client     *Client
	collection *mongo.Collection
}

func (s *IndicDAL) Client() *Client {
	return s.client
}

func (s *IndicDAL) Collection() *mongo.Collection {
	return s.collection
}

func (s *IndicDAL) Get(dep string, date time.Time) (*covidtracker.Indicator, error) {
	var result *covidtracker.Indicator

	err := s.collection.FindOne(s.client.Ctx, bson.M{"dep": dep, "extractDate": dateFilter(date)}).Decode(&result)
	switch err {
	case mongo.ErrNoDocuments:
		return nil, fmt.Errorf("no indicator document found with dep=%s date=%s", dep, date.Format("2006-02-01"))
	case nil:
		return result, nil
	default:
		return nil, fmt.Errorf("error while getting indicator: %s", err)
	}

	return result, nil
}

func (s *IndicDAL) GetRange(dep int, start, end time.Time) ([]*covidtracker.Indicator, error) {
	var result []*covidtracker.Indicator

	cur, err := s.collection.Find(s.client.Ctx, bson.M{"dep": dep, "extractDate": dateRangeFilter(start, end)})
	if err != nil {
		return nil, covidtracker.Errorf("error while getting case: %s", err)
	}
	defer cur.Close(s.client.Ctx)

	for cur.Next(s.client.Ctx) {
		var c *covidtracker.Indicator
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

func (s *IndicDAL) Upsert(inds ...*covidtracker.Indicator) error {

	if len(inds) == 0 {
		return covidtracker.Errorf("cannot upsert empty inds")
	}

	for _, ind := range inds {
		existing, err := s.Get(ind.Department, ind.ExtractDate)
		if err != nil || existing == nil { // not exist => add it
			bsonID := primitive.NewObjectID()
			ind.ID = covidtracker.IndicID(bsonID.Hex())
			if _, insErr := s.collection.InsertOne(s.client.Ctx, ind); insErr != nil {
				return errors.Wrap(insErr, "inserting new case")
			}
		} else { // existing, update only appropriate fields
			ind.ID = existing.ID
			if _, updErr := s.collection.UpdateOne(s.client.Ctx, bson.M{"_id": existing.ID}, bson.M{"$set": ind}); updErr != nil {
				return errors.Wrapf(err, "updating case dep=%s date%s", ind.Department, ind.ExtractDate.Format("2006-02-01"))
			}
		}

	}

	return nil
}
