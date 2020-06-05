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
	if err := j.aggregateReport(r, protects); err != nil {
		return nil, fmt.Errorf("cannot compute report: %s", err)
	}
	return r, nil
}

func (j *RiskJob) computeSegmentRisk(seg covidtracker.Segment, protects []covidtracker.Protection) (covidtracker.RiskSegment, error) {
	var (
		riskLevel   float64
		maskProtect float64
		gelProtect  float64
	)

	// basic coef for mask protection
	for _, prot := range protects {
		switch prot.Type {
		case covidtracker.MaskSewn:
			maskProtect = max(maskProtect, 0.71)
		case covidtracker.MaskSurgical:
			maskProtect = max(maskProtect, 0.85)
		case covidtracker.MaskFFPX:
			maskProtect = max(maskProtect, 0.99)
		case covidtracker.Gel:
			gelProtect = 0.99
		}
	}
	risk := covidtracker.RiskSegment{Segment: &seg}
	addPlus := func(category, msg string) {
		risk.Report.Pluses = append(risk.Report.Pluses, covidtracker.Statement{Category: category, Value: msg})
	}
	addMinus := func(category, msg string) {
		risk.Report.Minuses = append(risk.Report.Minuses, covidtracker.Statement{Category: category, Value: msg})
	}
	addAdvice := func(category, msg string) {
		risk.Report.Advices = append(risk.Report.Advices, covidtracker.Statement{Category: category, Value: msg})
	}

	duration := seg.Arrival.Sub(seg.Departure)

	// @todo: use real values here !
	switch seg.Transportation {
	case covidtracker.Aircraft:
		if duration != 0 && duration <= 4*time.Hour {
			addPlus(string(covidtracker.Aircraft), "Segment avion court")
			riskLevel = 0.55
			maskProtect *= 0.8
			gelProtect *= 0.05
		} else {
			if duration != 0 {
				addMinus(string(covidtracker.Aircraft), "Segment avion long")
			}
			riskLevel = 0.65
			maskProtect *= 0.75
			gelProtect *= 0.10
		}
	case covidtracker.TER, covidtracker.TGV:
		if duration != 0 && duration <= 2*time.Hour {
			riskLevel = 0.45
			maskProtect *= 0.8
			gelProtect *= 0.1
			addPlus(string(covidtracker.Aircraft), "Segment train court")
		} else {
			riskLevel = 0.55
			maskProtect *= 0.7
			gelProtect *= 0.2
			if duration != 0 {
				addMinus(string(covidtracker.Aircraft), "Segment train long")
			}
		}
	case covidtracker.CarSolo:
		addPlus(string(covidtracker.Car), "Vous êtes seul(e) dans la voiture")
		if duration != 0 && duration <= 3*time.Hour {
			addPlus(string(covidtracker.Car), "Segment voiture court")
			riskLevel = 0.01
			maskProtect *= 0.0
			gelProtect *= 0.1
		} else {
			if duration != 0 {
				addMinus(string(covidtracker.Car), "Segment voiture long, il faudra probablement s'arrêter à une pompe à essence")
			}
			riskLevel = 0.30
			maskProtect *= 0.75
			gelProtect *= 0.5
		}
	case covidtracker.CarDuo:
		addMinus(string(covidtracker.Car), "Vous êtes plusieurs dans la voiture")
		if duration != 0 && duration <= 3*time.Hour {
			addPlus(string(covidtracker.Car), "Segment voiture court")
			riskLevel = 0.6
			maskProtect *= 0.7
			gelProtect *= 0.1
		} else {
			if duration != 0 {
				addMinus(string(covidtracker.Car), "Segment voiture long, il faudra probablement s'arrêter à une pompe à essence")
			}

			riskLevel = 0.8
			maskProtect *= 0.8
			gelProtect *= 0.2
		}
	case covidtracker.CarGroup:
		addMinus(string(covidtracker.Car), "Vous êtes plusieurs dans la voiture")
		if duration != 0 && duration <= 3*time.Hour {
			addPlus(string(covidtracker.Car), "Segment voiture court")
			riskLevel = 0.7
			maskProtect *= 0.8
			gelProtect *= 0.1
		} else {
			if duration != 0 {
				addMinus(string(covidtracker.Car), "Segment voiture long, il faudra probablement s'arrêter à une pompe à essence")
			}
			riskLevel = 0.7
			maskProtect *= 0.8
			gelProtect *= 0.1
		}
	case covidtracker.TaxiSolo:
		riskLevel = 0.45
		maskProtect *= 0.7
		gelProtect *= 0.1
	case covidtracker.TaxiGroup:
		addMinus(string(covidtracker.TaxiGroup), "Vous êtes plusieurs passagers dans le taxi")
		riskLevel = 0.6
		maskProtect *= 0.8
		gelProtect *= 0.1
	case covidtracker.PublicTransports:
		if duration != 0 && duration <= 20*time.Hour {
			addPlus(string(covidtracker.PublicTransports), "Segment de transports en commun court")
			riskLevel = 0.5
			maskProtect *= 0.8
			gelProtect *= 0.25
		} else {
			addMinus(string(covidtracker.PublicTransports), "Segment de transports en commun long")
			riskLevel = 0.6
			maskProtect *= 0.8
			gelProtect *= 0.25
		}
	case covidtracker.Scooter:
		riskLevel = 0.01
		maskProtect *= 0.1
		gelProtect *= 0.05
	case covidtracker.Bike:
		riskLevel = 0.005
		maskProtect *= 0.1
		gelProtect *= 0.05
	default:
		return risk, fmt.Errorf("invalid transportation mode %q", seg.Transportation)
	}
	if duration > 4*time.Hour {
		addAdvice(string(covidtracker.Mask), "Votre voyage est long, emportez plusieurs masques")
	}
	riskLevel *= (1 - maskProtect)
	riskLevel *= (1 - gelProtect)
	risk.RiskLevel = riskLevel
	risk.ConfidenceLevel = 1 - riskLevel

	return risk, nil
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

func (j *RiskJob) aggregateReport(risk *covidtracker.Risk, protects []covidtracker.Protection) error {
	if risk == nil || len(risk.BySegments) == 0 {
		return fmt.Errorf("invalid risk for segments, no segment found")
	}
	risk.Report = covidtracker.Report{}
	for _, seg := range risk.BySegments {
		risk.Report.Pluses = append(risk.Report.Pluses, seg.Report.Pluses...)
		risk.Report.Minuses = append(risk.Report.Minuses, seg.Report.Minuses...)
		risk.Report.Advices = append(risk.Report.Advices, seg.Report.Advices...)
	}
	var (
		hasMask bool
		hasGel  bool
	)
	for _, prot := range protects {
		switch prot.Type {
		case covidtracker.MaskSewn, covidtracker.MaskSurgical, covidtracker.MaskFFPX:
			hasMask = true
		case covidtracker.Gel:
			hasGel = true
		}
	}
	if hasMask {
		risk.Report.Pluses = append(risk.Report.Pluses, covidtracker.Statement{Value: "Vous portez un masque", Category: "mask"})
	} else {
		risk.Report.Minuses = append(risk.Report.Minuses, covidtracker.Statement{Value: "Vous ne portez pas de masque", Category: "mask"})
		risk.Report.Advices = append(risk.Report.Advices, covidtracker.Statement{Value: "Portez un masque", Category: "mask"})
	}

	if hasGel {
		risk.Report.Pluses = append(risk.Report.Pluses, covidtracker.Statement{Value: "Vous utilisez du gel hydroalcoolique", Category: "gel"})
	} else {
		risk.Report.Minuses = append(risk.Report.Minuses, covidtracker.Statement{Value: "Vous n'avez pas de gel hydroalcoolique", Category: "gel"})
		risk.Report.Advices = append(risk.Report.Advices, covidtracker.Statement{Value: "Ayez du gel hydroalcoolique et lavez vous les main régulièrement", Category: "mask"})
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
