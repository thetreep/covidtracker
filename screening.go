package covidtracker

import "time"

type Screening struct {
	ID         HospID    `bson:"_id" json:"id"`
	Department int       `bson:"dep" json:"dep"`
	NoticeDate time.Time `bson:"noticeDate" json:"noticeDate"`

	AgeGroup string `bson:"ageGroup" json:"ageGroup"`

	Count         int32 `bson:"count" json:"count"`
	PositiveCount int32 `bson:"positiveCount" json:"positiveCount"`
	PositiveRate  int32 `bson:"positiveRate" json:"positiveRate"`
}

type ScreeningID string

type ScreeningService interface {
	RefreshScreening() ([]*Screening, error)
}

type ScreeningDAL interface {
	Get(dep int, date time.Time) ([]*Screening, error)
}
