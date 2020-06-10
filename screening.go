package covidtracker

import "time"

type Screening struct {
	ID         ScreeningID `bson:"_id" json:"id"`
	Department string      `bson:"dep" json:"dep"`
	NoticeDate time.Time   `bson:"noticeDate" json:"noticeDate"`

	Count         int `bson:"count" json:"count"`
	PositiveCount int `bson:"positiveCount" json:"positiveCount"`
	PositiveRate  int `bson:"positiveRate" json:"positiveRate"`
}

type ScreeningID string

type ScreeningService interface {
	RefreshScreening() ([]*Screening, error)
}

type ScreeningDAL interface {
	Get(dep string, date time.Time) (*Screening, error)
	GetRange(dep string, begin, end time.Time) ([]*Screening, error)
	Upsert(...*Screening) error
}
