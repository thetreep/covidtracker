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

package covidtracker_test

import (
	"testing"

	"github.com/thetreep/covidtracker"
	"github.com/thetreep/toolbox/test"
)

func TestGeocoding(t *testing.T) {

	t.Run("check postcode", func(t *testing.T) {

		tcases := map[string]struct {
			code string
			exp  bool
		}{
			"123 nok":      {code: "123", exp: false},
			"a123 nok":     {code: "a123", exp: false},
			"123a nok":     {code: "123a", exp: false},
			"78139 ok":     {code: "78139", exp: true},
			"00230 ok":     {code: "00230", exp: true},
			"90001 ok":     {code: "90001", exp: true},
			"12900013 nok": {code: "12900013", exp: false},
			"empty nok":    {code: "", exp: false},
		}

		for n, tcase := range tcases {
			got, want := covidtracker.Geo{
				Properties: covidtracker.Properties{
					GeoCoding: &covidtracker.GeoCoding{PostCode: tcase.code},
				},
			}.Check() == nil, tcase.exp
			test.Compare(t, got, want, n)
		}
	})
}
