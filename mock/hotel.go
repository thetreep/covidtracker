package mock

import (
	"github.com/thetreep/covidtracker"
)

type Hotel struct {
	//Jobs
	HotelsByPrefixFn      func(prefix string) ([]*covidtracker.Hotel, error)
	HotelsByPrefixInvoked bool

	//API
	SearchFn      func(query interface{}) ([]*covidtracker.Hotel, error)
	SearchInvoked bool
}

var (
	_ covidtracker.HotelJob = &Hotel{}
)

func (h *Hotel) Reset() {
	h.HotelsByPrefixInvoked = false
	h.SearchInvoked = false
}

func (h *Hotel) HotelsByPrefix(prefix string) ([]*covidtracker.Hotel, error) {
	h.HotelsByPrefixInvoked = true
	return h.HotelsByPrefixFn(prefix)
}

func (h *Hotel) Search(query interface{}) ([]*covidtracker.Hotel, error) {
	h.SearchInvoked = true
	return h.SearchFn(query)
}
