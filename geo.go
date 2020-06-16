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

import (
	"fmt"
	"regexp"
)

type Geo struct {
	Properties Properties `json:"properties"`
	Type       string     `json:"type"`
	Geometry   Geometry   `json:"geometry"`
}

type Properties struct {
	GeoCoding *GeoCoding `json:"geocoding,omitempty"`
	//props fr
	Name        *string  `json:"nom,omitempty"`
	PostalCode  *string  `json:"code,omitempty"`
	DepCode     *string  `json:"codeDepartement,omitempty"`
	RegionCode  *string  `json:"codeRegion,omitempty"`
	Population  *int     `json:"population,omitempty"`
	PostalCodes []string `json:"codesPostaux,omitempty"`
}

type Geometry struct {
	Coordinates []float64 `json:"coordinates"`
	Type        string    `json:"type"`
}

type GeoCoding struct {
	Type        string      `json:"type"`
	Accuracy    int         `json:"accuracy"`
	Label       string      `json:"label"`
	Name        string      `json:"name"`
	HouseNumber string      `json:"housenumber"`
	Street      string      `json:"street"`
	Locality    string      `json:"locality"`
	PostCode    string      `json:"postcode"`
	City        string      `json:"city"`
	District    *string     `json:"district,omitempty"`
	County      *string     `json:"county,omitempty"`
	State       *string     `json:"state,omitempty"`
	Country     string      `json:"country,omitempty"`
	Admin       AdminLevels `json:"admin"`
	Geohash     string      `json:"geohash"`
}

type AdminLevels struct {
	Level2 string `json:"level2"`
	Level4 string `json:"level4"`
	Level6 string `json:"level6"`
}

var regex = regexp.MustCompile(`^(([0-8][0-9])|(9[0-5])|(2[ab]))[0-9]{3}$`)

//Check checks format / value of Geo
func (g Geo) Check() error {

	if g.Properties.GeoCoding == nil {
		return nil
	}

	if regex == nil {
		return fmt.Errorf("postal code regexp cannot be created")
	}

	if !regex.MatchString(g.Properties.GeoCoding.PostCode) {
		return fmt.Errorf("postal code %q is invalid", g.Properties.GeoCoding.PostCode)
	}

	return nil
}

func (g Geo) Dep() (string, error) {

	if g.Properties.DepCode != nil {
		return *g.Properties.DepCode, nil
	}

	if g.Properties.GeoCoding != nil {
		if err := g.Check(); err != nil {
			return "", err
		}
		return g.Properties.GeoCoding.PostCode[:2], nil
	}

	return "", fmt.Errorf("department missing")
}
