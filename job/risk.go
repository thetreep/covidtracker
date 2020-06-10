package job

import (
	"fmt"
	"math"
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
		maskProtect float64
		gelProtect  float64
	)

	// basic coef for mask protection
	//@todo: update coefs according to scope
	for _, prot := range protects {
		switch prot.Type {
		case covidtracker.MaskSewn:
			maskProtect = math.Max(maskProtect, 0.71)
		case covidtracker.MaskSurgical:
			maskProtect = math.Max(maskProtect, 0.85)
		case covidtracker.MaskFFPX:
			maskProtect = math.Max(maskProtect, 0.99)
		case covidtracker.Gel:
			// @todo gelProtect = 0.99
			gelProtect = 0.1
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

	parameters, err := j.job.RiskParametersDAL.GetDefault()
	if err != nil {
		return risk, err
	}
	scope := covidtracker.ParameterScope{Transportation: seg.Transportation, Duration: seg.Transportation.Duration(seg.Departure, seg.Arrival)}
	params, ok := parameters.ParametersByScope[scope]
	if !ok {
		j.job.logger.Info(j.job.Ctx, "no parameters for this scope with transportation and duration, fallback on transportation 'Normal'")
		params, ok = parameters.ParametersByScope[covidtracker.ParameterScope{Transportation: seg.Transportation, Duration: covidtracker.Normal}]
		if !ok {
			return risk, covidtracker.Errorf("no parameters are defined for scope %v", scope)
		}
	}
	for _, p := range params.Pluses {
		addPlus(string(seg.Transportation), p)
	}
	for _, m := range params.Minuses {
		addMinus(string(seg.Transportation), m)
	}
	for _, a := range params.Advices {
		addAdvice(string(seg.Transportation), a)
	}

	if duration > 4*time.Hour {
		addAdvice(string(covidtracker.Mask), "Votre voyage est long, emportez plusieurs masques")
	}

	// @todo: fill here with actual numbers for department
	probaContagious := probaContagiousPerson([]int{250, 350, 250, 120, 100}, 1000000)

	// infected if infected with contact OR with direct OR with indirect contact
	riskLevel := probaUnionIndep(
		probaInfected(params.NbContact, probaContagious, params.ProbaContagionContact),
		probaUnionIndep(
			probaInfected(params.NbDirect, probaContagious, params.ProbaContagionDirect),
			probaInfected(params.NbIndirect, probaContagious, params.ProbaContagionIndirect),
		),
	)

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

	// aggregate and remove duplicates
	pluses, minuses, advices := make(map[covidtracker.Statement]struct{}), make(map[covidtracker.Statement]struct{}), make(map[covidtracker.Statement]struct{})
	for _, seg := range risk.BySegments {
		for _, p := range seg.Report.Pluses {
			pluses[p] = struct{}{}
		}
		for _, m := range seg.Report.Minuses {
			minuses[m] = struct{}{}
		}
		for _, a := range seg.Report.Advices {
			advices[a] = struct{}{}
		}
	}
	for p := range pluses {
		risk.Report.Pluses = append(risk.Report.Pluses, p)
	}
	for m := range minuses {
		risk.Report.Minuses = append(risk.Report.Minuses, m)
	}
	for a := range advices {
		risk.Report.Advices = append(risk.Report.Advices, a)
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

// probaInfected computes the probability of being infected from nbPersons with a probability of being contagious probaContagious
// and a probability of contagion in the related segment of probaContagion
// being infected = being infected by person p1 OR by p2 OR ... OR by pn
// being infected by person pi = person i is contagious AND contagion happens. We do the approximation that the level of contagion of this person makes more or less likely the successful contagion i.e. probaContagion and probaContagious are independent
// thus we have P = probaUnionEquiprobableIndep(probaContagious*probaContagion, nbPersons)
func probaInfected(nbPersons int, probaContagious, probaContagion float64) float64 {
	return probaUnionEquiprobableIndep(probaContagion*probaContagious, nbPersons)
}

// Compute the probability that a random person in a specific place (for, example a French department) is contagious of Covid-19
// nbSuspiciousCase is the number of suspicious Covid-19 in the place for the last 14 days
// totalPopulation is the number of people living in this place
//
// We consider (with approximation) that people with no or minor symptoms do not go to the emergency and that all people with major symptoms go to the emergency
// Approximation => We take into account only the 14 last days to compute the total number of people infected that are still infectious.
//
func probaContagiousPerson(nbSuspiciousCasesEmergency []int, totalPop int) float64 {
	nbEmergency := 0.
	for _, nb := range nbSuspiciousCasesEmergency {
		nbEmergency += float64(nb)
	}

	// nbSympt is the estimate number of symptomatic people infected with Covid-19.
	// We suppose that the people with suspicious case at emergency are people with major symptoms
	// Thus we generalise this number of people with symptoms to number of symptomatic people actually infected with Covid-19, extrapolating data the following study:
	// source (preprint): Fontanet, A., Tondeur, L., Madec, Y., Grant, R., Besombes, C., Jolly, N., ... & Temmam, S. (2020). Cluster of COVID-19 in northern France: A retrospective closed cohort study.
	// we consider an infected person the one with at least one positive serology
	// In such study, we have 321 with major symptoms from which 121 are positive + 21 positives with minor symptoms (considered here not taking into account in the suspicious cases at emergency)
	nbSympt := nbEmergency * (121. + 21.) / 321.

	// asymptomatic proportion = proportion of asymptomatically infected individuals among the total number of infected individuals.
	// c.f. asymptomatic proportion src : Mizumoto, K., Kagaya, K., Zarebski, A., & Chowell, G. (2020). Estimating the asymptomatic proportion of coronavirus disease 2019 (COVID-19) cases on board the Diamond Princess cruise ship, Yokohama, Japan, 2020. Eurosurveillance, 25(10), 2000180.
	asymptProp := 17.9 / 100.

	// nbAsympt is the estimate number of asymptomatic people infected with Covid-19
	nbAsympt := float64(nbSympt*asymptProp) / (1. - asymptProp) // asymptProp = nbAsympt / (nbAsympt + nbSympt)
	nbInfected := nbSympt + nbAsympt

	// Probabilty that a random person in a specific place is infected with Covid-19
	probaInfected := nbInfected / float64(totalPop)

	// Probabilty that a random person in a specific place is infectious with Covid-19
	probaInfectious := probaInfected * 1. // For now, we take probaInfectious = probaInfected
	return probaInfectious
}

// probaUnionEquiprobableIndep compute the probability of the union of n independent equiprobable events with probability p
// P = (p-1)*(1-p)^(n-1) + 1
func probaUnionEquiprobableIndep(p float64, n int) float64 {
	if n == 0 {
		return 0.
	}
	return (p-1)*math.Pow(1-p, float64(n-1)) + 1
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
