package covidtracker

import "time"

type Case struct {
	ID         CaseID    `bson:"_id" json:"id"`
	Department string    `bson:"dep" json:"dep"`
	NoticeDate time.Time `bson:"noticeDate" json:"noticeDate"`

	//HospServiceCountRelated is the number of hospital services reporting at least one case
	HospServiceCountRelated int `bson:"hospServiceCountRelated" json:"hospServiceCountRelated"`
}

type CaseID string

type CaseService interface {
	RefreshCase() ([]*Case, error)
}

type CaseDAL interface {
	Get(dep string, date time.Time) (*Case, error)
	GetRange(dep string, begin, end time.Time) ([]*Case, error)
	Upsert(...*Case) error
}
