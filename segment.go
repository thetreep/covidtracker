package covidtracker

import "time"

type Segment struct {
	ID             SegID          `bson:"_id" json:"-"`
	Origin         string         `bson:"origin" json:"origin"`
	Destination    string         `bson:"destination" json:"destination"`
	Departure      time.Time      `bson:"departure" json:"departure"`
	Arrival        time.Time      `bson:"arrival" json:"arrival"`
	Transportation Transportation `bson:"transportation" json:"transportation"`
	RiskLevel      *float32       `bson:"riskLevel" json:"riskLevel,omitempty"`
}

type SegID string
