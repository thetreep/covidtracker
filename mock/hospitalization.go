package mock

import (
	"time"

	"github.com/thetreep/covidtracker"
)

type HospDAL struct {
	GetFn      func(dep int, date time.Time) ([]*covidtracker.Hospitalization, error)
	GetInvoked bool

	GetRangeFn      func(dep int, start, end time.Time) ([]*covidtracker.Hospitalization, error)
	GetRangeInvoked bool

	UpsertFn      func([]*covidtracker.Hospitalization) error
	UpsertInvoked bool
}

var _ covidtracker.HospDAL = &HospDAL{}

func (m *HospDAL) Get(dep int, date time.Time) ([]*covidtracker.Hospitalization, error) {
	m.GetInvoked = true
	return m.GetFn(dep, date)
}

func (m *HospDAL) GetRange(dep int, start, end time.Time) ([]*covidtracker.Hospitalization, error) {
	m.GetRangeInvoked = true
	return m.GetRangeFn(dep, start, end)
}

func (m *HospDAL) Upsert(hosps []*covidtracker.Hospitalization) error {
	m.UpsertInvoked = true
	return m.UpsertFn(hosps)
}
