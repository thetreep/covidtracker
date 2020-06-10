package mock

import (
	"time"

	"github.com/thetreep/covidtracker"
)

type CaseDAL struct {
	GetFn      func(dep string, date time.Time) (*covidtracker.Case, error)
	GetInvoked bool

	GetRangeFn      func(dep string, start, end time.Time) ([]*covidtracker.Case, error)
	GetRangeInvoked bool

	UpsertFn      func(...*covidtracker.Case) error
	UpsertInvoked bool
}

var _ covidtracker.CaseDAL = &CaseDAL{}

func (m *CaseDAL) Get(dep string, date time.Time) (*covidtracker.Case, error) {
	m.GetInvoked = true
	return m.GetFn(dep, date)
}

func (m *CaseDAL) GetRange(dep string, start, end time.Time) ([]*covidtracker.Case, error) {
	m.GetRangeInvoked = true
	return m.GetRangeFn(dep, start, end)
}

func (m *CaseDAL) Upsert(cases ...*covidtracker.Case) error {
	m.UpsertInvoked = true
	return m.UpsertFn(cases...)
}
