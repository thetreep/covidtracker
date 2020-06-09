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
