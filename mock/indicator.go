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

type IndicDAL struct {
	GetFn      func(dep string, date time.Time) (*covidtracker.Indicator, error)
	GetInvoked bool

	GetRangeFn      func(dep int, start, end time.Time) ([]*covidtracker.Indicator, error)
	GetRangeInvoked bool

	UpsertFn      func(...*covidtracker.Indicator) error
	UpsertInvoked bool
}

var _ covidtracker.IndicDAL = &IndicDAL{}

func (m *IndicDAL) Get(dep string, date time.Time) (*covidtracker.Indicator, error) {
	m.GetInvoked = true
	return m.GetFn(dep, date)
}

func (m *IndicDAL) GetRange(dep int, start, end time.Time) ([]*covidtracker.Indicator, error) {
	m.GetRangeInvoked = true
	return m.GetRangeFn(dep, start, end)
}

func (m *IndicDAL) Upsert(inds ...*covidtracker.Indicator) error {
	m.UpsertInvoked = true
	return m.UpsertFn(inds...)
}
