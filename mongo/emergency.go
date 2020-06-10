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
var _ covidtracker.EmergencyDAL = &EmergencyDAL{}
var _ Accessor = &EmergencyDAL{}

// EmergencyDAL represents a service for managing suppliers.
type EmergencyDAL struct {
	client     *Client
	collection *mongo.Collection
}

func (s *EmergencyDAL) Client() *Client {
	return s.client
}

func (s *EmergencyDAL) Collection() *mongo.Collection {
	return s.collection
}

func (s *EmergencyDAL) Get(dep string, date time.Time) (*covidtracker.Emergency, error) {
	var result *covidtracker.Emergency

	err := s.collection.FindOne(s.client.Ctx, bson.M{"dep": dep, "passageDate": dateFilter(date)}).Decode(&result)
	switch err {
	case mongo.ErrNoDocuments:
		return nil, fmt.Errorf("no emergency document found with dep=%s date=%s", dep, date.Format("2006-02-01"))
	case nil:
		return result, nil
	default:
		return nil, fmt.Errorf("error while getting emergency: %s", err)
	}
}

func (s *EmergencyDAL) GetRange(dep string, start, end time.Time) ([]*covidtracker.Emergency, error) {
	var result []*covidtracker.Emergency

	cur, err := s.collection.Find(s.client.Ctx, bson.M{"dep": dep, "passageDate": dateRangeFilter(start, end)})
	if err != nil {
		return nil, covidtracker.Errorf("error while getting case: %s", err)
	}
	defer cur.Close(s.client.Ctx)

	for cur.Next(s.client.Ctx) {
		var c *covidtracker.Emergency
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

func (s *EmergencyDAL) Upsert(ems ...*covidtracker.Emergency) error {

	if len(ems) == 0 {
		return covidtracker.Errorf("cannot upsert empty ems")
	}

	for _, em := range ems {
		existing, err := s.Get(em.Department, em.PassageDate)
		if err != nil || existing == nil { // not exist => add it
			bsonID := primitive.NewObjectID()
			em.ID = covidtracker.EmergencyID(bsonID.Hex())
			if _, insErr := s.collection.InsertOne(s.client.Ctx, em); insErr != nil {
				return errors.Wrap(insErr, "inserting new case")
			}
		} else { // existing, update only appropriate fields
			em.ID = existing.ID
			if _, updErr := s.collection.UpdateOne(s.client.Ctx, bson.M{"_id": existing.ID}, bson.M{"$set": em}); updErr != nil {
				return errors.Wrapf(err, "updating case dep=%s date%s", em.Department, em.PassageDate.Format("2006-02-01"))
			}
		}

	}

	return nil
}
