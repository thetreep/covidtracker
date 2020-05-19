package covidtracker

import "time"

type Case struct {
	ID         CaseID    `bson:"_id" json:"id"`
	Department uint8     `bson:"dep" json:"dep"`
	NoticeDate time.Time `bson:"noticeDate" json:"noticeDate"`

	//HospServiceCountReport is the number of hospital services reporting at least one case
	HospServiceCountRelated int32 `bson:"hospServiceCountRelated" json:"hospServiceCountRelated"`
}

type CaseID string

type CaseService interface {
	RefreshCase() ([]*Case, error)
}

type CaseDAL interface {
	Get(dep uint8, date time.Time) ([]*Case, error)
}
