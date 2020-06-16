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

type EmergencyDAL struct {
	GetFn      func(dep string, date time.Time) (*covidtracker.Emergency, error)
	GetInvoked bool

	GetRangeFn      func(dep string, start, end time.Time) ([]*covidtracker.Emergency, error)
	GetRangeInvoked bool

	UpsertFn      func(...*covidtracker.Emergency) error
	UpsertInvoked bool
}

var _ covidtracker.EmergencyDAL = &EmergencyDAL{}

func (m *EmergencyDAL) Get(dep string, date time.Time) (*covidtracker.Emergency, error) {
	m.GetInvoked = true
	return m.GetFn(dep, date)
}
func (m *EmergencyDAL) GetRange(dep string, start, end time.Time) ([]*covidtracker.Emergency, error) {
	m.GetRangeInvoked = true
	return m.GetRangeFn(dep, start, end)
}

func (m *EmergencyDAL) Upsert(ems ...*covidtracker.Emergency) error {
	m.UpsertInvoked = true
	return m.UpsertFn(ems...)
}
