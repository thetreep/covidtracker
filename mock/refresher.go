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

import "github.com/thetreep/covidtracker"

type Refresher struct {
	//Jobs
	RefreshFn      func(caseD covidtracker.CaseDAL, emD covidtracker.EmergencyDAL, hospD covidtracker.HospDAL, indicD covidtracker.IndicDAL, scrD covidtracker.ScreeningDAL) error
	RefreshInvoked bool
}

var (
	_ covidtracker.RefreshJob = &Refresher{}
)

func (m *Refresher) Refresh(caseD covidtracker.CaseDAL, emD covidtracker.EmergencyDAL, hospD covidtracker.HospDAL, indicD covidtracker.IndicDAL, scrD covidtracker.ScreeningDAL) error {
	m.RefreshInvoked = true
	return m.RefreshFn(caseD, emD, hospD, indicD, scrD)
}
