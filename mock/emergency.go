package mock

import (
	"time"

	"github.com/thetreep/covidtracker"
)

type EmergencyDAL struct {
	GetFn      func(dep int, date time.Time) ([]*covidtracker.Emergency, error)
	GetInvoked bool

	GetRangeFn      func(dep int, start, end time.Time) ([]*covidtracker.Emergency, error)
	GetRangeInvoked bool

	UpsertFn      func([]*covidtracker.Emergency) error
	UpsertInvoked bool
}

var _ covidtracker.EmergencyDAL = &EmergencyDAL{}

func (m *EmergencyDAL) Get(dep int, date time.Time) ([]*covidtracker.Emergency, error) {
	m.GetInvoked = true
	return m.GetFn(dep, date)
}
func (m *EmergencyDAL) GetRange(dep int, start, end time.Time) ([]*covidtracker.Emergency, error) {
	m.GetRangeInvoked = true
	return m.GetRangeFn(dep, start, end)
}

func (m *EmergencyDAL) Upsert(ems []*covidtracker.Emergency) error {
	m.UpsertInvoked = true
	return m.UpsertFn(ems)
}
