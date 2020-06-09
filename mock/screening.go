package mock

import (
	"time"

	"github.com/thetreep/covidtracker"
)

type ScreeningDAL struct {
	GetFn      func(dep int, date time.Time) (*covidtracker.Screening, error)
	GetInvoked bool

	GetRangeFn      func(dep int, start, end time.Time) ([]*covidtracker.Screening, error)
	GetRangeInvoked bool

	UpsertFn      func(...*covidtracker.Screening) error
	UpsertInvoked bool
}

var _ covidtracker.ScreeningDAL = &ScreeningDAL{}

func (m *ScreeningDAL) Get(dep int, date time.Time) (*covidtracker.Screening, error) {
	m.GetInvoked = true
	return m.GetFn(dep, date)
}

func (m *ScreeningDAL) GetRange(dep int, start, end time.Time) ([]*covidtracker.Screening, error) {
	m.GetRangeInvoked = true
	return m.GetRangeFn(dep, start, end)
}

func (m *ScreeningDAL) Upsert(scrs ...*covidtracker.Screening) error {
	m.UpsertInvoked = true
	return m.UpsertFn(scrs...)
}
