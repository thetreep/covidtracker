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
	Report          Report        `bson:"report" json:"report"`
}

type Report struct {
	Minuses []Statement `bson:"minuses" json:"minuses"`
	Pluses  []Statement `bson:"pluses" json:"pluses"`
	Advices []Statement `bson:"advices" json:"advices"`
}

type Statement struct {
	Value    string `bson:"value" json:"value"`
	Category string `bson:"category" json:"category"`
}

//RiskID identifies a Risk
type RiskID string

//RiskSegment is the risk and the confidence level for a given segment
type RiskSegment struct {
	ID RiskSegID `bson:"_id" json:"id"`

	*Segment `bson:"segment" json:"segment"`

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
	ComputeRisk(segs []Segment, protects []Protection) (*Risk, error)
}
