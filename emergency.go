package covidtracker

import "time"

//Emergency regroups the stats about visit at emergency room
type Emergency struct {
	ID          EmergencyID `bson:"_id" json:"id"`
	Department  string      `bson:"dep" json:"dep"`
	PassageDate time.Time   `bson:"passageDate" json:"passageDate"`

	//Count is the number of visits
	Count int `bson:"count" json:"count"`
	//Cov19SuspCount is the number of suspicious covid19 patient amoung the visits
	Cov19SuspCount int `bson:"cov19SuspCount" json:"cov19SuspCount"`

	//Cov19SuspicionHosp is the amount of hospitalized for covid-19 suspicion amoung the visits
	Cov19SuspHosp int `bson:"cov19SuspHospitalized" json:"cov19SuspHospitalized"`
	//TotalSOSMedAct is the amount of medical act reported by SOS Medecin
	TotalSOSMedAct int `bson:"totalSosMedAct" json:"totalSosMedAct"`
	//TotalSOSMedAct is the amount of medical act reported by SOS Medecin concerning the COVID-19
	SOSMedCov19SuspAct int `bson:"cov19SosMedAct" json:"sosMedMaleAct"`
}

type EmergencyID string

type EmergencyService interface {
	RefreshEmergency() ([]*Emergency, error)
}

type EmergencyDAL interface {
	Get(dep string, date time.Time) (*Emergency, error)
	GetRange(dep string, begin, end time.Time) ([]*Emergency, error)
	Upsert(...*Emergency) error
}
