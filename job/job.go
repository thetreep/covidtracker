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

package job

import (
	"context"
	"time"

	"github.com/thetreep/covidtracker"
)

type Job struct {
	Ctx context.Context
	Now func() time.Time

	RiskDAL           covidtracker.RiskDAL
	RiskJob           RiskJob
	RiskParametersDAL covidtracker.RiskParametersDAL
	EmergencyDAL      covidtracker.EmergencyDAL

	HotelDAL covidtracker.HotelDAL
	HotelJob HotelJob

	RefreshJob RefreshJob

	logger covidtracker.Logfer
}

// NewJob creates a new job
func NewJob(log covidtracker.Logfer, refresher covidtracker.Refresher) *Job {
	j := &Job{Now: time.Now, Ctx: context.Background()}
	j.RiskJob.job = j
	j.RefreshJob = RefreshJob{
		job:       j,
		Refresher: refresher,
	}
	j.HotelJob.job = j
	j.logger = log
	return j
}

// Risk returns the risk job associated with the client
func (j *Job) Risk() covidtracker.RiskJob { return &j.RiskJob }

// Refresh returns the refresher job associated with the client
func (j *Job) Refresh() covidtracker.RefreshJob { return &j.RefreshJob }

// Hotels returns hotels search by the client
func (j *Job) Hotels() covidtracker.HotelJob { return &j.HotelJob }
