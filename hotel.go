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

package covidtracker

import "fmt"

type HotelID string

type Hotel struct {
	ID            HotelID  `bson:"_id" json:"id"`
	Name          string   `bson:"name" json:"name"`
	Address       string   `bson:"address" json:"address"`
	City          string   `bson:"city" json:"city"`
	ZipCode       string   `bson:"zipCode" json:"zipCode"`
	Country       string   `bson:"country" json:"country"`
	ImageURL      string   `bson:"imageUrl" json:"imageUrl"`
	SanitaryInfos []string `bson:"sanitaryInfos" json:"sanitaryInfos"`
	SanitaryNote  float64  `bson:"sanitaryNote" json:"sanitaryNote"`
	SanitaryNorm  string   `bson:"sanitaryNorm" json:"sanitaryNorm"`
}

func (h *Hotel) Dep() (string, error) {
	if len(h.ZipCode) > 2 {
		return h.ZipCode[:2], nil
	}

	return "", fmt.Errorf("department missing")
}

type HotelDAL interface {
	Get(id HotelID) (*Hotel, error)
	Insert(hotels []*Hotel) ([]*Hotel, error)
}

type HotelJob interface {
	HotelsByPrefix(city string, prefix string) ([]*Hotel, error)
}
