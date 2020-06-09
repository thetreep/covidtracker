package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func dateFilter(d time.Time) bson.M {
	start := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
	end := time.Date(d.Year(), d.Month(), d.Day()+1, 0, 0, 0, 0, time.UTC)
	return bson.M{
		"$gte": start,
		"$lt":  end,
	}
}

func dateRangeFilter(start, end time.Time) bson.M {
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	end = time.Date(end.Year(), end.Month(), end.Day()+1, 0, 0, 0, 0, time.UTC)
	return bson.M{
		"$gte": start,
		"$lt":  end,
	}
}
