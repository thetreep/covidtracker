package datagouv

type API interface {
	Get(query interface{}) ([]*Risk, error)
}

var _ API = &covidtracker.RiskAPI
