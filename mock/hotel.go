package mock

import (
	"github.com/thetreep/covidtracker"
)

type Hotel struct {
	//Jobs
	HotelsByPrefixFn      func(prefix string) ([]*covidtracker.Hotel, error)
	HotelsByPrefixInvoked bool

	//DAL
	GetFn      func(id covidtracker.HotelID) (*covidtracker.Hotel, error)
	GetInvoked bool

	InsertFn      func(hotels ...*covidtracker.Hotel) error
	InsertInvoked bool

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
	h.GetInvoked = false
	h.InsertInvoked = false
}

func (h *Hotel) HotelsByPrefix(prefix string) ([]*covidtracker.Hotel, error) {
	h.HotelsByPrefixInvoked = true
	return h.HotelsByPrefixFn(prefix)
}

func (h *Hotel) Get(id covidtracker.HotelID) (*covidtracker.Hotel, error) {
	h.GetInvoked = true
	return h.GetFn(id)
}
func (h *Hotel) Insert(hotels ...*covidtracker.Hotel) error {
	h.InsertInvoked = true
	return h.InsertFn(hotels...)
}

func (h *Hotel) Search(query interface{}) ([]*covidtracker.Hotel, error) {
	h.SearchInvoked = true
	return h.SearchFn(query)
}
