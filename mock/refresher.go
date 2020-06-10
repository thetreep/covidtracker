package mock

import "github.com/thetreep/covidtracker"

type Refresher struct {
	//Jobs
	RefreshFn      func(caseD covidtracker.CaseDAL, emD covidtracker.EmergencyDAL, hospD covidtracker.HospDAL, indicD covidtracker.IndicDAL, scrD covidtracker.ScreeningDAL) error
	RefreshInvoked bool
}

var (
	_ covidtracker.RefreshJob = &Refresher{}
)

func (m *Refresher) Refresh(caseD covidtracker.CaseDAL, emD covidtracker.EmergencyDAL, hospD covidtracker.HospDAL, indicD covidtracker.IndicDAL, scrD covidtracker.ScreeningDAL) error {
	m.RefreshInvoked = true
	return m.RefreshFn(caseD, emD, hospD, indicD, scrD)
}
