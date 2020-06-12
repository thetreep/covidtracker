package mock

import (
	"github.com/thetreep/covidtracker"
)

type RiskParameters struct {
	GetDefaultFn      func() (*covidtracker.RiskParameters, error)
	GetDefaultInvoked bool

	InsertFn      func(p *covidtracker.RiskParameters) error
	InsertInvoked bool
}

var (
	_ covidtracker.RiskParametersDAL = &RiskParameters{}
)

func (h *RiskParameters) Reset() {
	h.GetDefaultInvoked = false
	h.InsertInvoked = false
}
func (h *RiskParameters) GetDefault() (*covidtracker.RiskParameters, error) {
	h.GetDefaultInvoked = true
	return h.GetDefaultFn()
}
func (h *RiskParameters) Insert(p *covidtracker.RiskParameters) error {
	h.InsertInvoked = true
	return h.InsertFn(p)
}
