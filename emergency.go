package covidtracker

import "time"

//Emergency regroups the stats about visit at emergency room
type Emergency struct {
	ID         EmergencyID `bson:"_id" json:"id"`
	Department int         `bson:"dep" json:"dep"`
	NoticeDate time.Time   `bson:"noticeDate" json:"noticeDate"`

	AgeGroup string `bson:"ageGroup" json:"ageGroup"`

	//Count is the number of visits
	Count int32 `bson:"count" json:"count"`
	//Cov19SuspCount is the number of suspicious covid19 patient amoung the visits
	Cov19SuspCount int32 `bson:"cov19SuspCount" json:"cov19SuspCount"`

	//Cov19SuspicionHosp is the amount of hospitalized for covid-19 suspicion amoung the visits
	Cov19SuspHosp int32 `bson:"cov19SuspHospitalized" json:"cov19SuspHospitalized"`
	//TotalSOSMedAct is the amount of medical act reported by SOS Medecin
	TotalSOSMedAct int32 `bson:"totalSosMedAct" json:"totalSosMedAct"`
	//TotalSOSMedAct is the amount of medical act reported by SOS Medecin concerning the COVID-19
	SOSMedCov19SuspAct int32 `bson:"cov19SosMedAct" json:"sosMedMaleAct"`
}

type EmergencyID string

type EmergencyService interface {
	RefreshEmergency() ([]*Emergency, error)
}

type EmergencyDAL interface {
	Get(dep int, date time.Time) ([]*Emergency, error)
}
