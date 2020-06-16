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

type Screening struct {
	ID         ScreeningID `bson:"_id" json:"id"`
	Department string      `bson:"dep" json:"dep"`
	NoticeDate time.Time   `bson:"noticeDate" json:"noticeDate"`

	Count         int `bson:"count" json:"count"`
	PositiveCount int `bson:"positiveCount" json:"positiveCount"`
	PositiveRate  int `bson:"positiveRate" json:"positiveRate"`
}

type ScreeningID string

type ScreeningService interface {
	RefreshScreening() ([]*Screening, error)
}

type ScreeningDAL interface {
	Get(dep string, date time.Time) (*Screening, error)
	GetRange(dep string, begin, end time.Time) ([]*Screening, error)
	Upsert(...*Screening) error
}
