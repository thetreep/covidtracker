package job

import (
	"github.com/thetreep/covidtracker"
)

// RiskJob represents a job for computing risk
type RiskJob struct {
	job *Job
}

var _ covidtracker.RiskJob = &RiskJob{}

func (j *RiskJob) ComputeRisk(segs []covidtracker.Segment, protects []covidtracker.Protection) (*covidtracker.Risk, error) {
	return nil, nil
}
