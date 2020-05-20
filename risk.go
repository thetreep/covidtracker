package covidtracker

import (
	"time"
)

//Risk is the definition of risk and confidence level of a trip
type Risk struct {
	ID              RiskID        `bson:"_id" json:"id"`
	NoticeDate      time.Time     `bson:"noticeDate" json:"noticeDate"`
	ConfidenceLevel float64       `bson:"confidenceLevel" json:"confidenceLevel"`
	RiskLevel       float64       `bson:"riskLevel" json:"riskLevel"`
	BySegments      []RiskSegment `bson:"bySegments" json:"bySegments"`
	Pluses          []string      `bson:"pluses" json:"pluses"`
	Minuses         []string      `bson:"minuses" json:"minuses"`
	Advices         []string      `bson:"advices" json:"advices"`
}

//RiskID identifies a Risk
type RiskID string

//RiskSegment is the risk and the confidence level for a given segment
type RiskSegment struct {
	ID RiskSegID `bson:"_id" json:"id"`

	*Segment `bson:"seg" json:"seg"`

	RiskLevel       float64 `bson:"riskLevel" json:"riskLevel"`
	ConfidenceLevel float64 `bson:"confidenceLevel" json:"confidenceLevel"`
}

//RiskSegID identifies a RiskSegment
type RiskSegID string

//RiskDAL defines the data access layer of risk data
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
