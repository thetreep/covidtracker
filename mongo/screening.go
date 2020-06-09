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
var _ covidtracker.ScreeningDAL = &ScreeningDAL{}
var _ Accessor = &ScreeningDAL{}

// ScreeningDAL represents a service for managing suppliers.
type ScreeningDAL struct {
	client     *Client
	collection *mongo.Collection
}

func (s *ScreeningDAL) Client() *Client {
	return s.client
}

func (s *ScreeningDAL) Collection() *mongo.Collection {
	return s.collection
}

func (s *ScreeningDAL) Get(dep string, date time.Time) (*covidtracker.Screening, error) {
	var result *covidtracker.Screening

	err := s.collection.FindOne(s.client.Ctx, bson.M{"dep": dep, "noticeDate": dateFilter(date)}).Decode(&result)
	switch err {
	case mongo.ErrNoDocuments:
		return nil, fmt.Errorf("no screening document found with dep=%s date=%s", dep, date.Format("2006-02-01"))
	case nil:
		return result, nil
	default:
		return nil, fmt.Errorf("error while getting screening: %s", err)
	}

	return result, nil
}

func (s *ScreeningDAL) GetRange(dep int, start, end time.Time) ([]*covidtracker.Screening, error) {
	var result []*covidtracker.Screening

	cur, err := s.collection.Find(s.client.Ctx, bson.M{"dep": dep, "noticeDate": dateRangeFilter(start, end)})
	if err != nil {
		return nil, covidtracker.Errorf("error while getting case: %s", err)
	}
	defer cur.Close(s.client.Ctx)

	for cur.Next(s.client.Ctx) {
		var c *covidtracker.Screening
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

func (s *ScreeningDAL) Upsert(scrs ...*covidtracker.Screening) error {

	if len(scrs) == 0 {
		return covidtracker.Errorf("cannot upsert empty scrs")
	}

	for _, scr := range scrs {
		existing, err := s.Get(scr.Department, scr.NoticeDate)
		if err != nil || existing == nil { // not exist => add it
			bsonID := primitive.NewObjectID()
			scr.ID = covidtracker.ScreeningID(bsonID.Hex())
			if _, insErr := s.collection.InsertOne(s.client.Ctx, scr); insErr != nil {
				return errors.Wrap(insErr, "inserting new case")
			}
		} else { // existing, update only appropriate fields
			scr.ID = existing.ID
			if _, updErr := s.collection.UpdateOne(s.client.Ctx, bson.M{"_id": existing.ID}, bson.M{"$set": scr}); updErr != nil {
				return errors.Wrapf(err, "updating case dep=%s date%s", scr.Department, scr.NoticeDate.Format("2006-02-01"))
			}
		}

	}

	return nil
}