package mock

import (
	"time"

	"github.com/thetreep/covidtracker"
)

type IndicDAL struct {
	GetFn      func(dep string, date time.Time) (*covidtracker.Indicator, error)
	GetInvoked bool

	GetRangeFn      func(dep int, start, end time.Time) ([]*covidtracker.Indicator, error)
	GetRangeInvoked bool

	UpsertFn      func(...*covidtracker.Indicator) error
	UpsertInvoked bool
}

var _ covidtracker.IndicDAL = &IndicDAL{}

func (m *IndicDAL) Get(dep string, date time.Time) (*covidtracker.Indicator, error) {
	m.GetInvoked = true
	return m.GetFn(dep, date)
}

func (m *IndicDAL) GetRange(dep int, start, end time.Time) ([]*covidtracker.Indicator, error) {
	m.GetRangeInvoked = true
	return m.GetRangeFn(dep, start, end)
}

func (m *IndicDAL) Upsert(inds ...*covidtracker.Indicator) error {
	m.UpsertInvoked = true
	return m.UpsertFn(inds...)
}
