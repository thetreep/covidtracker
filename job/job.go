package job

import (
	"context"
	"time"

	"github.com/thetreep/covidtracker"
)

type Job struct {
	Ctx context.Context
	Now func() time.Time

	RiskDAL covidtracker.RiskDAL
	RiskJob RiskJob

	HotelJob HotelJob
}

// NewJob creates a new job
func NewJob() *Job {
	j := &Job{Now: time.Now, Ctx: context.Background()}
	j.RiskJob.job = j
	j.HotelJob.job = j
	return j
}

// Risk returns the risk job associated with the client
func (j *Job) Risk() covidtracker.RiskJob { return &j.RiskJob }

// Hotels returns hotels search by the client
func (j *Job) Hotels() covidtracker.HotelJob { return &j.HotelJob }
