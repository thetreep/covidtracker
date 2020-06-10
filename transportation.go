package covidtracker

import (
	"time"
)

type Transportation string

const (
	TGV              Transportation = "tgv"
	TER              Transportation = "ter"
	Aircraft         Transportation = "aircraft"
	Car              Transportation = "car"
	CarSolo          Transportation = "car-solo"
	CarDuo           Transportation = "car-duo"
	CarGroup         Transportation = "car-group"
	TaxiSolo         Transportation = "taxi-solo"
	TaxiGroup        Transportation = "taxi-group"
	PublicTransports Transportation = "public-transports"
	Scooter          Transportation = "scooter"
	Bike             Transportation = "bike"
)

func (t *Transportation) Duration(departure, arrival time.Time) TransportationDuration {
	if t == nil || departure.IsZero() || arrival.IsZero() {
		return Normal
	}
	duration := arrival.Sub(departure)
	if duration <= 0 {
		return Normal
	}
	switch *t {
	case Aircraft:
		if duration <= 2*time.Hour {
			return Short
		} else if duration <= 6*time.Hour {
			return Normal
		} else {
			return Long
		}
	case TER, TGV:
		if duration <= 1*time.Hour {
			return Short
		} else if duration <= 3*time.Hour {
			return Normal
		} else {
			return Long
		}
	case CarSolo, CarDuo, CarGroup:
		if duration <= 30*time.Minute {
			return Short
		} else if duration <= 3*time.Hour {
			return Normal
		} else {
			return Long
		}
	case TaxiSolo, TaxiGroup:
		if duration <= 20*time.Minute {
			return Short
		} else if duration <= 45*time.Minute {
			return Normal
		} else {
			return Long
		}
	case PublicTransports:
		if duration <= 20*time.Minute {
			return Short
		} else if duration <= 1*time.Hour {
			return Normal
		} else {
			return Long
		}
	default:
		return Normal
	}
}

type TransportationDuration string

const (
	Short  TransportationDuration = "short"
	Normal TransportationDuration = "normal"
	Long   TransportationDuration = "long"
)
