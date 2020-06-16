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

import (
	"time"
)

//Hospitalization defines the usefull data about hospitalization
type Hospitalization struct {
	ID         HospID    `bson:"_id" json:"id"`
	Department string    `bson:"dep" json:"dep"`
	Date       time.Time `bson:"date" json:"date"`

	//Count is the number of patient hospitalized
	Count int `bson:"count" json:"count"`
	//CriticalCount is the number of patient in resuscitation or critical care
	CriticalCount int `bson:"critical" json:"critical"`
	//ReturnHomeCount is the number of patient that returned home
	ReturnHomeCount int `bson:"returnHome" json:"returnHome"`
	//DeathCount is the number of deaths
	DeathCount int `bson:"deaths" json:"deaths"`
}

type HospID string

type HospService interface {
	RefreshHospitalization() ([]*Hospitalization, error)
}

type HospDAL interface {
	Get(dep string, date time.Time) (*Hospitalization, error)
	GetRange(dep string, begin, end time.Time) ([]*Hospitalization, error)
	Upsert(...*Hospitalization) error
}
