package job

import (
	"context"
	"time"

	"github.com/thetreep/covidtracker"
)

type Job struct {
	Ctx context.Context
	Now func() time.Time

	RiskDAL       covidtracker.RiskDAL
	RiskJob       RiskJob
	ParametersDAL covidtracker.RiskParametersDAL
}

// NewJob creates a new job
func NewJob() *Job {
	j := &Job{Now: time.Now, Ctx: context.Background()}
	j.RiskJob.job = j
	return j
}

// Risk returns the risk job associated with the client
func (j *Job) Risk() covidtracker.RiskJob { return &j.RiskJob }
