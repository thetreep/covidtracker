package covidtracker

import (
	"time"
)

//Hospitalization defines the usefull data about hospitalization
type Hospitalization struct {
	ID         HospID    `bson:"_id" json:"id"`
	Department uint8     `bson:"dep" json:"dep"`
	NoticeDate time.Time `bson:"noticeDate" json:"noticeDate"`

	Sex string `bson:"sex" json:"sex"`

	//Count is the number of patient hospitalized
	Count int32 `bson:"count" json:"count"`
	//CriticalCount is the number of patient in resuscitation or critical care
	CriticalCount int32 `bson:"critical" json:"critical"`
	//ReturnHomeCount is the number of patient that returned home
	ReturnHomeCount int32 `bson:"returnHome" json:"returnHome"`
	//DeathCount is the number of deaths
	DeathCount int32 `bson:"deaths" json:"deaths"`
}

type HospID string

type HospService interface {
	RefreshHospitalization() ([]*Hospitalization, error)
}

type HospDAL interface {
	Get(dep uint8, date time.Time) ([]*Hospitalization, error)
}
