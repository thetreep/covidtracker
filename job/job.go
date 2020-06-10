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

	RefreshJob RefreshJob
	HotelJob   HotelJob

	logger covidtracker.Logfer
}

// NewJob creates a new job
func NewJob(refresher covidtracker.Refresher) *Job {
	j := &Job{Now: time.Now, Ctx: context.Background()}
	j.RiskJob.job = j
	j.RefreshJob = RefreshJob{
		job:       j,
		Refresher: refresher,
	}
	j.HotelJob.job = j
	return j
}

// Risk returns the risk job associated with the client
func (j *Job) Risk() covidtracker.RiskJob { return &j.RiskJob }

// Refresh returns the refresher job associated with the client
func (j *Job) Refresh() covidtracker.RefreshJob { return &j.RefreshJob }

// Hotels returns hotels search by the client
func (j *Job) Hotels() covidtracker.HotelJob { return &j.HotelJob }
