package job

import (
	"fmt"

	"github.com/thetreep/covidtracker"
	"github.com/thetreep/covidtracker/job/cds"
)

type HotelJob struct {
	job *Job
}

var _ covidtracker.HotelJob = &HotelJob{}

func (j *HotelJob) HotelsByPrefix(prefix string) ([]*covidtracker.Hotel, error) {
	hotels, err := cds.Service.HotelsByPrefix(prefix)
	if err != nil {
		return nil, fmt.Errorf("cannot get hotels: %s", err)
	}
	var hotelsFr []*covidtracker.Hotel
	for _, h := range hotels {
		// We skip hotels that are not in France
		if h.Country != "FR" {
			continue
		}
		hotelsFr = append(hotelsFr, h)
	}
	return hotelsFr, nil
}
