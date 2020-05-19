package covidtracker

import "time"

type Segment struct {
	ID             SegID
	Origin         string
	Destination    string
	DateTime       time.Time
	Transportation Transportation
	RiskLevel      float32
}

type SegID string
