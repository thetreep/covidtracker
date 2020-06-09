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
	Report          Report  `bson:"report" json:"report"`
}

type Parameters struct {
	// Use to splecify that these are the default parameters
	IsDefault bool `bson:"default" json:"default"`

	// The parameters associated to a scope
	ParametersByScope map[ParameterScope]RiskParameter `bson:"parameters_by_scope" json:"parameters_by_scope"`
}

type RiskParameter struct {
	// The number of persons with direct projection possible
	NbDirect int `bson:"nb_direct" json:"nb_direct"`

	// The probability of contagion via direct projection with an infectious person
	ProbaContagionDirect float64 `bson:"proba_contagion_direct" json:"proba_contagion_direct"`

	// The number of persons with direct contact with the person
	NbContact int `bson:"nb_contact" json:"nb_contact"`

	// The probability of contagion via direct contact with an infectious person
	ProbaContagionContact float64 `bson:"proba_contagion_contact" json:"proba_contagion_contact"`

	// The number of persons with indirect contact
	NbIndirect int `bson:"nb_indirect" json:"nb_indirect"`

	// The probability of contagion via indirect contact with an infectious person
	ProbaContagionIndirect float64 `bson:"proba_contagion_indirect" json:"proba_contagion_indirect"`

	// The Pluses of this kind of segment
	Pluses []string `bson:"pluses" json:"pluses"`

	// The Minuses of this kind of segment
	Minuses []string `bson:"minuses" json:"minuses"`

	// The Advices of this kind of segment
	Advices []string `bson:"advices" json:"advices"`
}

type ParameterScope struct {
	Transportation Transportation         `bson:"transportation" json:"transportation"`
	Duration       TransportationDuration `bson:"duration" json:"duration"`
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

//ParametersDAL defines the data access layer of risk parameters
type ParametersDAL interface {
	GetDefault() (*Parameters, error)
	Insert(p *Parameters) error
}
