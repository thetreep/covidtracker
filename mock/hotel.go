/*
	This file is part of covidtracker also known as EviteCovid .

    Copyright 2020 the Treep

    covdtracker is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    covidtracker is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with covidtracker.  If not, see <https://www.gnu.org/licenses/>.
*/

package mock

import (
	"github.com/thetreep/covidtracker"
)

type Hotel struct {
	//Jobs
	HotelsByPrefixFn      func(city string, prefix string) ([]*covidtracker.Hotel, error)
	HotelsByPrefixInvoked bool

	//DAL
	GetFn      func(id covidtracker.HotelID) (*covidtracker.Hotel, error)
	GetInvoked bool

	InsertFn      func(h []*covidtracker.Hotel) ([]*covidtracker.Hotel, error)
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

func (h *Hotel) HotelsByPrefix(city string, prefix string) ([]*covidtracker.Hotel, error) {
	h.HotelsByPrefixInvoked = true
	return h.HotelsByPrefixFn(city, prefix)
}

func (h *Hotel) Get(id covidtracker.HotelID) (*covidtracker.Hotel, error) {
	h.GetInvoked = true
	return h.GetFn(id)
}
func (h *Hotel) Insert(hotels []*covidtracker.Hotel) ([]*covidtracker.Hotel, error) {
	h.InsertInvoked = true
	hotels, err := h.InsertFn(hotels)
	return hotels, err
}

func (h *Hotel) Search(query interface{}) ([]*covidtracker.Hotel, error) {
	h.SearchInvoked = true
	return h.SearchFn(query)
}
