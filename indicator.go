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

type Indicator struct {
	ID          IndicID   `bson:"_id" json:"id"`
	ExtractDate time.Time `bson:"extractDate" json:"extractDate"`
	Department  string    `bson:"dep" json:"dep"`
	Color       string    `bson:"color" json:"color"`
}

type IndicID string

type IndicService interface {
	RefreshIndicator() ([]*Indicator, error)
}

type IndicDAL interface {
	Get(dep string, date time.Time) (*Indicator, error)
	GetRange(dep int, begin, end time.Time) ([]*Indicator, error)
	Upsert(...*Indicator) error
}
