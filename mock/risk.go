package mock

import "github.com/thetreep/covidtracker"

type Risk struct {
	//Jobs
	ComputeRiskFn      func(segs []covidtracker.Segment, protects []covidtracker.Protection) (*covidtracker.Risk, error)
	ComputeRiskInvoked bool

	//DAL
	GetFn      func(id covidtracker.RiskID) (*covidtracker.Risk, error)
	GetInvoked bool

	InsertFn      func(r ...*covidtracker.Risk) error
	InsertInvoked bool

	//API
	EstimateFn      func(query interface{}) ([]*covidtracker.Risk, error)
	EstimateInvoked bool
}

var (
	_ covidtracker.RiskJob = &Risk{}
	_ covidtracker.RiskDAL = &Risk{}
)

func (r *Risk) Reset() {
	r.InsertInvoked = false
	r.ComputeRiskInvoked = false
	r.GetInvoked = false
	r.EstimateInvoked = false
}
func (r *Risk) ComputeRisk(segs []covidtracker.Segment, protects []covidtracker.Protection) (*covidtracker.Risk, error) {
	r.ComputeRiskInvoked = true
	return r.ComputeRiskFn(segs, protects)
}
func (r *Risk) Get(id covidtracker.RiskID) (*covidtracker.Risk, error) {
	r.GetInvoked = true
	return r.GetFn(id)
}
func (r *Risk) Insert(risks ...*covidtracker.Risk) error {
	r.InsertInvoked = true
	return r.InsertFn(risks...)
}
func (r *Risk) Estimate(query interface{}) ([]*covidtracker.Risk, error) {
	r.EstimateInvoked = true
	return r.EstimateFn(query)
}
