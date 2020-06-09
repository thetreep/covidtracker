package job

import (
	"fmt"

	"github.com/thetreep/covidtracker"
	"github.com/thetreep/covidtracker/job/cds"
)

var (
	Service cdsAPI
)

type cdsAPI interface {
	HotelsByPrefix(p string) ([]covidtracker.Hotel, error)
}

type HotelJob struct {
	job *Job
}

var _ covidtracker.HotelJob = &HotelJob{}

func (j *HotelJob) HotelsByPrefix(prefix string) ([]covidtracker.Hotel, error) {
	hotels, err := cds.Service.HotelsByPrefix(prefix)
	if err != nil {
		return nil, fmt.Errorf("cannot get hotels: %s", err)
	}
	return hotels, nil
}
