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

func TestRefreshHospitalization(t *testing.T) {
	t.Run("parsing", func(t *testing.T) {
		api, server := DatagouvServers(t)
		defer func() {
			api.Close()
			server.Close()
		}()

		s := datagouv.NewService(context.Background(), &logger.Logger{})
		s.BasePath = api.URL

		hosp, err := s.RefreshHospitalization()
		if err != nil {
			t.Fatal(err)
		}
		if len(hosp) == 0 {
			t.Fatal("unexpected empty cases")
		}

		timeFn := func(s string) time.Time {
			t, _ := time.Parse("2006-01-02", s)
			return t
		}

		expected := []*covidtracker.Hospitalization{
			&covidtracker.Hospitalization{
				Department:      "976",
				Date:            timeFn("2020-05-26"),
				Count:           26,
				CriticalCount:   6,
				ReturnHomeCount: 106,
				DeathCount:      6,
			},
			&covidtracker.Hospitalization{
				Department:      "22",
				Date:            timeFn("2020-05-20"),
				Count:           24,
				CriticalCount:   2,
				ReturnHomeCount: 81,
				DeathCount:      18,
			},
			&covidtracker.Hospitalization{
				Department:      "42",
				Date:            timeFn("2020-05-17"),
				Count:           127,
				CriticalCount:   5,
				ReturnHomeCount: 374,
				DeathCount:      67,
			},
			&covidtracker.Hospitalization{
				Department:      "01",
				Date:            timeFn("2020-03-18"),
				Count:           13,
				CriticalCount:   0,
				ReturnHomeCount: 2,
				DeathCount:      0,
			},
			&covidtracker.Hospitalization{
				Department:      "93",
				Date:            timeFn("2020-03-18"),
				Count:           184,
				CriticalCount:   45,
				ReturnHomeCount: 20,
				DeathCount:      10,
			},
		}

		test.Compare(t, hosp, expected, "unexpected hospitalization")

	})
}
