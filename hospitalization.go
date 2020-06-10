package covidtracker

import (
	"time"
)

//Hospitalization defines the usefull data about hospitalization
type Hospitalization struct {
	ID         HospID    `bson:"_id" json:"id"`
	Department string    `bson:"dep" json:"dep"`
	Date       time.Time `bson:"date" json:"date"`

	//Count is the number of patient hospitalized
	Count int `bson:"count" json:"count"`
	//CriticalCount is the number of patient in resuscitation or critical care
	CriticalCount int `bson:"critical" json:"critical"`
	//ReturnHomeCount is the number of patient that returned home
	ReturnHomeCount int `bson:"returnHome" json:"returnHome"`
	//DeathCount is the number of deaths
	DeathCount int `bson:"deaths" json:"deaths"`
}

type HospID string

type HospService interface {
	RefreshHospitalization() ([]*Hospitalization, error)
}

type HospDAL interface {
	Get(dep string, date time.Time) (*Hospitalization, error)
	GetRange(dep string, begin, end time.Time) ([]*Hospitalization, error)
	Upsert(...*Hospitalization) error
}
