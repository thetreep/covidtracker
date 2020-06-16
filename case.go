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

type Case struct {
	ID         CaseID    `bson:"_id" json:"id"`
	Department string    `bson:"dep" json:"dep"`
	NoticeDate time.Time `bson:"noticeDate" json:"noticeDate"`

	//HospServiceCountRelated is the number of hospital services reporting at least one case
	HospServiceCountRelated int `bson:"hospServiceCountRelated" json:"hospServiceCountRelated"`
}

type CaseID string

type CaseService interface {
	RefreshCase() ([]*Case, error)
}

type CaseDAL interface {
	Get(dep string, date time.Time) (*Case, error)
	GetRange(dep string, begin, end time.Time) ([]*Case, error)
	Upsert(...*Case) error
}
