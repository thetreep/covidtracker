package job

import (
	"math/rand"
	"time"

	"github.com/thetreep/covidtracker"
)

// RiskJob represents a job for computing risk
type RiskJob struct {
	job *Job
}

var _ covidtracker.RiskJob = &RiskJob{}

func (j *RiskJob) ComputeRisk(segs []covidtracker.Segment, protects []covidtracker.Protection) (*covidtracker.Risk, error) {

	level := rand.Float64()
	r := &covidtracker.Risk{
		RiskLevel:       level,
		ConfidenceLevel: 1. - level,
		NoticeDate:      time.Now(),
		Report: covidtracker.Report{
			Pluses: []covidtracker.Statement{
				{Value: "zone verte sur tout le trajet"},
				{Value: "trajet de moins de 100 km"},
			},
			Minuses: []covidtracker.Statement{
				{Value: "Long séjour"},
			},
			Advices: []covidtracker.Statement{
				{Value: "Prévoyez 3 masques pour votre trajet"},
				{Value: "Nettoyez-vous les mains régulièrement"},
			},
		},
	}

	for _, seg := range segs {
		r.BySegments = append(r.BySegments, covidtracker.RiskSegment{Segment: &seg, RiskLevel: level, ConfidenceLevel: 1. - level})
	}

	return r, nil
}
