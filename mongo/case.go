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

func (s *CaseDAL) Get(dep string, date time.Time) (*covidtracker.Case, error) {
	var result *covidtracker.Case

	err := s.collection.FindOne(s.client.Ctx, bson.M{"dep": dep, "noticeDate": dateFilter(date)}).Decode(&result)
	switch err {
	case mongo.ErrNoDocuments:
		return nil, fmt.Errorf("no document found with dep=%s date=%s", dep, date.Format("2006-02-01"))
	case nil:
		return result, nil
	default:
		return nil, fmt.Errorf("error while getting case: %s", err)
	}

	return result, nil
}

func (s *CaseDAL) GetRange(dep int, start, end time.Time) ([]*covidtracker.Case, error) {
	var result []*covidtracker.Case

	cur, err := s.collection.Find(s.client.Ctx, bson.M{"dep": dep, "noticeDate": dateRangeFilter(start, end)})
	if err != nil {
		return nil, covidtracker.Errorf("error while getting case: %s", err)
	}
	defer cur.Close(s.client.Ctx)

	for cur.Next(s.client.Ctx) {
		var c *covidtracker.Case
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

func (s *CaseDAL) Upsert(cases ...*covidtracker.Case) error {

	if len(cases) == 0 {
		return covidtracker.Errorf("cannot upsert empty cases")
	}

	for _, c := range cases {
		existing, err := s.Get(c.Department, c.NoticeDate)
		if err != nil || existing == nil { // not exist => add it
			bsonID := primitive.NewObjectID()
			c.ID = covidtracker.CaseID(bsonID.Hex())
			if _, insErr := s.collection.InsertOne(s.client.Ctx, c); insErr != nil {
				return errors.Wrap(insErr, "inserting new case")
			}
		} else { // existing, update only appropriate fields
			c.ID = existing.ID
			if _, updErr := s.collection.UpdateOne(s.client.Ctx, bson.M{"_id": existing.ID}, bson.M{"$set": c}); updErr != nil {
				return errors.Wrapf(err, "updating case dep=%s date%s", c.Department, c.NoticeDate.Format("2006-02-01"))
			}
		}

	}

	return nil
}
