package covidtracker

type RefreshJob interface {
	Refresh(CaseDAL, EmergencyDAL, HospDAL, IndicDAL, ScreeningDAL) error
}

type Refresher interface {
	CaseService
	EmergencyService
	HospService
	IndicService
	ScreeningService
}
