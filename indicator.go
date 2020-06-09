package covidtracker

import "time"

type Indicator struct {
	ID          IndicID   `bson:"_id" json:"id"`
	ExtractDate time.Time `bson:"extractDate" json:"extractDate"`
	Department  string    `bson:"dep" json:"dep"`
	Color       string    `bson:"color" json:"color"`
}

type IndicID string

type IndicService interface {
	RefreshIndicator() ([]*Indicator, error)
}

type IndicDAL interface {
	Get(dep string, date time.Time) (*Indicator, error)
	GetRange(dep int, begin, end time.Time) ([]*Indicator, error)
	Upsert(...*Indicator) error
}
