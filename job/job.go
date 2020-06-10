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
