/*
	This file is part of covidtracker also known as EviteCovid .

    Copyright 2020 the Treep

    covdtracker is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    covidtracker is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with covidtracker.  If not, see <https://www.gnu.org/licenses/>.
*/

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
