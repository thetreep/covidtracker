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
	"github.com/thetreep/covidtracker"
)

//RefreshJob defines the job to refresh data
type RefreshJob struct {
	job       *Job
	Refresher covidtracker.Refresher
	logger    covidtracker.Logfer
}

//Refresh refreshes the data and save it in database
func (j *RefreshJob) Refresh(caseD covidtracker.CaseDAL, emD covidtracker.EmergencyDAL, hospD covidtracker.HospDAL, indicD covidtracker.IndicDAL, scrD covidtracker.ScreeningDAL) error {
	logger := j.job.logger
	ctx := j.job.Ctx

	logger.Debug(ctx, "start refreshing data")

	var (
		cases  []*covidtracker.Case
		ems    []*covidtracker.Emergency
		indics []*covidtracker.Indicator
		hosps  []*covidtracker.Hospitalization
		scrs   []*covidtracker.Screening
		err    error
	)

	if caseD != nil {
		cases, err = j.Refresher.RefreshCase()
		if err != nil {
			return err
		}
		if err := caseD.Upsert(cases...); err != nil {
			return err
		}
	} else {
		logger.Debug(ctx, "no dal for case data is defined")
	}

	if emD != nil {
		ems, err = j.Refresher.RefreshEmergency()
		if err != nil {
			return err
		}
		if err := emD.Upsert(ems...); err != nil {
			return err
		}
	} else {
		logger.Debug(ctx, "no dal for emergency data defined")
	}

	if hospD != nil {
		hosps, err = j.Refresher.RefreshHospitalization()
		if err != nil {
			return err
		}
		if err := hospD.Upsert(hosps...); err != nil {
			return err
		}
	} else {
		logger.Debug(ctx, "no dal for hospitalization data is defined")
	}

	if indicD != nil {
		indics, err = j.Refresher.RefreshIndicator()
		if err != nil {
			return err
		}
		if err := indicD.Upsert(indics...); err != nil {
			return err
		}
	} else {
		logger.Debug(ctx, "no dal for indicator is defined")
	}

	if scrD != nil {
		scrs, err = j.Refresher.RefreshScreening()
		if err != nil {
			return err
		}
		if err := scrD.Upsert(scrs...); err != nil {
			return err
		}
	} else {
		logger.Debug(ctx, "no dal for screening data is defined")
	}
	return nil
}
