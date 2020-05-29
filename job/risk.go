package job

import (
	"fmt"
	"time"

	"github.com/thetreep/covidtracker"
)

// RiskJob represents a job for computing risk
type RiskJob struct {
	job *Job
}

var _ covidtracker.RiskJob = &RiskJob{}

func (j *RiskJob) ComputeRisk(segs []covidtracker.Segment, protects []covidtracker.Protection) (*covidtracker.Risk, error) {
	r := &covidtracker.Risk{
		NoticeDate: time.Now(),
	}
	for i, seg := range segs {
		segRisk, err := j.computeSegmentRisk(seg, protects)
		if err != nil {
			return nil, fmt.Errorf("cannot compute risk for segment %d: %s", i, err)
		}
		r.BySegments = append(r.BySegments, segRisk)
	}
	if err := j.aggregateSegmentRisk(r); err != nil {
		return nil, fmt.Errorf("cannot aggregate risk of %d segments: %s", len(segs), err)
	}
	if err := j.computeReport(r, protects); err != nil {
		return nil, fmt.Errorf("cannot compute report: %s", err)
	}
	return r, nil
}

func (j *RiskJob) computeSegmentRisk(seg covidtracker.Segment, protects []covidtracker.Protection) (covidtracker.RiskSegment, error) {
	var riskLevel float64
	var maskProtect float64
	var gelProtect float64

	// basic coef for mask protection
	for _, prot := range protects {
		switch prot.Type {
		case covidtracker.MaskSewn:
			maskProtect = max(maskProtect, 0.35)
		case covidtracker.MaskSurgical:
			maskProtect = max(maskProtect, 0.80)
		case covidtracker.MaskFFPX:
			maskProtect = max(maskProtect, 0.90)
		case covidtracker.Gel:
			gelProtect = 0.95
		}
	}

	// @todo: use real values here !

	// @todo: use real values here !
	switch seg.Transportation {
	case covidtracker.Aircraft:
		riskLevel = 0.40
		maskProtect *= 0.9
		gelProtect *= 0.05
	case covidtracker.TER:
		riskLevel = 0.35
		maskProtect *= 0.9
		gelProtect *= 0.1
	case covidtracker.TGV:
		riskLevel = 0.32
		maskProtect *= 0.9
		gelProtect *= 0.1
	case covidtracker.CarSolo:
		riskLevel = 0.01
		maskProtect *= 0.0
		gelProtect *= 0.1
	case covidtracker.CarDuo:
		riskLevel = 0.4
		maskProtect *= 0.7
		gelProtect *= 0.1
	case covidtracker.CarGroup:
		riskLevel = 0.5
		maskProtect *= 0.8
		gelProtect *= 0.1
	case covidtracker.TaxiSolo:
		riskLevel = 0.32
		maskProtect *= 0.7
		gelProtect *= 0.1
	case covidtracker.TaxiGroup:
		riskLevel = 0.4
		maskProtect *= 0.8
		gelProtect *= 0.1
	case covidtracker.PublicTransports:
		riskLevel = 0.4
		maskProtect *= 0.9
		gelProtect *= 0.2
	case covidtracker.Scooter:
		riskLevel = 0.01
		maskProtect *= 0.1
		gelProtect *= 0.05
	case covidtracker.Bike:
		riskLevel = 0.005
		maskProtect *= 0.1
		gelProtect *= 0.05
	default:
		return covidtracker.RiskSegment{}, fmt.Errorf("invalid transportation mode %q", seg.Transportation)
	}
	riskLevel *= (1 - maskProtect)
	riskLevel *= (1 - gelProtect)

	return covidtracker.RiskSegment{Segment: &seg, RiskLevel: riskLevel, ConfidenceLevel: 1 - riskLevel}, nil

}

func (j *RiskJob) aggregateSegmentRisk(risk *covidtracker.Risk) error {
	if risk == nil || len(risk.BySegments) == 0 {
		return fmt.Errorf("invalid risk for segments, no segment found")
	}
	var probasSegment []float64
	for _, seg := range risk.BySegments {
		probasSegment = append(probasSegment, seg.RiskLevel)
	}
	risk.RiskLevel = probaUnionIndepSlice(0, probasSegment)
	risk.ConfidenceLevel = 1 - risk.RiskLevel
	return nil
}

func (j *RiskJob) computeReport(risk *covidtracker.Risk, protects []covidtracker.Protection) error {
	if risk == nil || len(risk.BySegments) == 0 {
		return fmt.Errorf("invalid risk for segments, no segment found")
	}
	risk.Report = covidtracker.Report{
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
	}

	return nil
}

// probaUnionIndep compute the probability of the union of two independent events
// p(a OR b) = p(a) + p(b) - p(a AND b) = p(a) + p(b) - p(a)*p(b)
func probaUnionIndep(pa, pb float64) float64 {
	return pa + pb - pa*pb
}

// probaUnionIndepSlice compute recursively the probability of the union of a slice of mutually independent events
// p(a1 OR a2 OR ... OR an) = p(a1 OR (a2 OR ... OR an))
func probaUnionIndepSlice(fromIndex int, probas []float64) float64 {
	if fromIndex >= len(probas) {
		return 0.
	}
	if fromIndex == len(probas)-1 {
		return probas[fromIndex]
	}
	return probaUnionIndep(probas[fromIndex], probaUnionIndepSlice(fromIndex+1, probas))
}

func max(a, b float64) float64 {
	if a >= b {
		return a
	}
	return b
}
