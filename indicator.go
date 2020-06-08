package covidtracker

import "time"

type Indicator struct {
	ID          IndicID   `bson:"_id" json:"id"`
	ExtractDate time.Time `bson:"extractDate" json:"extractDate"`
	Department  int       `bson:"dep" json:"dep"`
	Color       string    `bson:"color" json:"color"`
}

type IndicID string

type IndicService interface {
	RefreshIndicator() ([]*Indicator, error)
}

type IndicDAL interface {
	Get(dep int, date time.Time) ([]*Indicator, error)
	GetRange(dep int, begin, end time.Time) ([]*Indicator, error)
	CreateNew([]*Indicator) error
}
