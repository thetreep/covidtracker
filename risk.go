package covidtracker

import (
	"time"
)

//Risk is the definition of risk and confidence level of a trip
type Risk struct {
	ID              RiskID    `bson:"_id" json:"id"`
	NoticeDate      time.Time `bson:"notice_date" json:"notice_date"`
	ConfidenceScore float64   `bson:"confidence_score" json:"confidence_score"`
	RiskLevel       float64   `bson:"risk_level" json:"risk_level"`

	//TODO add routes and protections ?
	// Routes          []Route      `bson:"routes" json:"routes"`
	// Protections     []Protection `bson:"protections" json:"protections"`
}

//RiskID identifies a Risk
type RiskID string

//RiskAPI defines the data access layer of risk data
type RiskDAL interface {
	Get(id RiskID) (*Risk, error)
	Insert(r ...*Risk) error
}

//RiskJob defines the job to implements risk data logic
type RiskJob interface {
	ComputeRisk() (*Risk, error)
}

//RiskAPI defines the api to get risk data
type RiskAPI interface {
	Get(query interface{}) ([]*Risk, error)
}
