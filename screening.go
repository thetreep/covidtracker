package covidtracker

import "time"

type Screening struct {
	ID         HospID    `bson:"_id" json:"id"`
	Department int       `bson:"dep" json:"dep"`
	NoticeDate time.Time `bson:"noticeDate" json:"noticeDate"`

	AgeGroup AgeGroup `bson:"ageGroup" json:"ageGroup"`

	Count         int `bson:"count" json:"count"`
	PositiveCount int `bson:"positiveCount" json:"positiveCount"`
	PositiveRate  int `bson:"positiveRate" json:"positiveRate"`
}

type ScreeningID string

type ScreeningService interface {
	RefreshScreening() ([]*Screening, error)
}

type ScreeningDAL interface {
	Get(dep int, date time.Time) ([]*Screening, error)
	GetRange(dep int, begin, end time.Time) ([]*Screening, error)
	Upsert(scrs []*Screening) error
}
