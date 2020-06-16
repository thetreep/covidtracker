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

package datagouv_test

import (
	"context"
	"testing"
	"time"

	"github.com/thetreep/covidtracker"
	"github.com/thetreep/covidtracker/job/datagouv"
	"github.com/thetreep/covidtracker/logger"
	"github.com/thetreep/toolbox/test"
)

func TestRefreshIndicator(t *testing.T) {
	t.Run("parsing", func(t *testing.T) {
		api, server := DatagouvServers(t)
		defer func() {
			api.Close()
			server.Close()
		}()

		s := datagouv.NewService(context.Background(), &logger.Logger{})
		s.BasePath = api.URL

		inds, err := s.RefreshIndicator()
		if err != nil {
			t.Fatal(err)
		}
		if len(inds) == 0 {
			t.Fatal("unexpected empty indicators")
		}

		timeFn := func(s string) time.Time {
			t, _ := time.Parse("2006-01-02", s)
			return t
		}

		expected := []*covidtracker.Indicator{
			{
				ExtractDate: timeFn("2020-05-07"),
				Department:  "59",
				Color:       "rouge",
			},
			{
				ExtractDate: timeFn("2020-05-07"),
				Department:  "78",
				Color:       "rouge",
			},
			{
				ExtractDate: timeFn("2020-05-04"),
				Department:  "88",
				Color:       "rouge",
			},
			{
				ExtractDate: timeFn("2020-05-03"),
				Department:  "15",
				Color:       "orange",
			},
		}

		test.Compare(t, inds, expected, "unexpected indicator")

	})
}
