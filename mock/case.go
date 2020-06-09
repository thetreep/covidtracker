package mock

import (
	"time"

	"github.com/thetreep/covidtracker"
)

type CaseDAL struct {
	GetFn      func(dep int, date time.Time) ([]*covidtracker.Case, error)
	GetInvoked bool

	GetRangeFn      func(dep int, start, end time.Time) ([]*covidtracker.Case, error)
	GetRangeInvoked bool

	UpsertFn      func([]*covidtracker.Case) error
	UpsertInvoked bool
}

var _ covidtracker.CaseDAL = &CaseDAL{}

func (m *CaseDAL) Get(dep int, date time.Time) ([]*covidtracker.Case, error) {
	m.GetInvoked = true
	return m.GetFn(dep, date)
}

func (m *CaseDAL) GetRange(dep int, start, end time.Time) ([]*covidtracker.Case, error) {
	m.GetRangeInvoked = true
	return m.GetRangeFn(dep, start, end)
}

func (m *CaseDAL) Upsert(cases []*covidtracker.Case) error {
	m.UpsertInvoked = true
	return m.UpsertFn(cases)
}
