package cds

import (
	"github.com/thetreep/covidtracker"
)

var (
	Service cdsAPI
)

type cdsAPI interface {
	HotelsByPrefix(p string) ([]covidtracker.Hotel, error)
}

type tracedCDSService struct {
	service cdsAPI
}

func (s *tracedCDSService) HotelsByPrefix(p string) ([]covidtracker.Hotel, error) {
	out, err := s.service.HotelsByPrefix(p)
	return out, err
}
