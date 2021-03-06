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

package mock

import (
	"time"

	"github.com/thetreep/covidtracker"
)

type HospDAL struct {
	GetFn      func(dep string, date time.Time) (*covidtracker.Hospitalization, error)
	GetInvoked bool

	GetRangeFn      func(dep string, start, end time.Time) ([]*covidtracker.Hospitalization, error)
	GetRangeInvoked bool

	UpsertFn      func(...*covidtracker.Hospitalization) error
	UpsertInvoked bool
}

var _ covidtracker.HospDAL = &HospDAL{}

func (m *HospDAL) Get(dep string, date time.Time) (*covidtracker.Hospitalization, error) {
	m.GetInvoked = true
	return m.GetFn(dep, date)
}

func (m *HospDAL) GetRange(dep string, start, end time.Time) ([]*covidtracker.Hospitalization, error) {
	m.GetRangeInvoked = true
	return m.GetRangeFn(dep, start, end)
}

func (m *HospDAL) Upsert(hosps ...*covidtracker.Hospitalization) error {
	m.UpsertInvoked = true
	return m.UpsertFn(hosps...)
}
