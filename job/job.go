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
	RiskJob *RiskJob
}

// NewJob creates a new job
func NewJob() *Job {
	j := &Job{Now: time.Now, Ctx: context.Background()}
	j.RiskJob.job = j
	return j
}
