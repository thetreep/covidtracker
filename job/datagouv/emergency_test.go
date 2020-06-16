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

func TestRefreshEmergency(t *testing.T) {
	t.Run("parsing", func(t *testing.T) {
		api, server := DatagouvServers(t)
		defer func() {
			api.Close()
			server.Close()
		}()

		s := datagouv.NewService(context.Background(), &logger.Logger{})
		s.BasePath = api.URL

		sos, err := s.RefreshEmergency()
		if err != nil {
			t.Fatal(err)
		}
		if len(sos) == 0 {
			t.Fatal("unexpected empty cases")
		}

		timeFn := func(s string) time.Time {
			t, _ := time.Parse("2006-01-02", s)
			return t
		}

		expected := []*covidtracker.Emergency{
			&covidtracker.Emergency{
				Department:         "62",
				PassageDate:        timeFn("2020-05-25"),
				Count:              113,
				Cov19SuspCount:     1,
				Cov19SuspHosp:      0,
				TotalSOSMedAct:     0,
				SOSMedCov19SuspAct: 0,
			},
			&covidtracker.Emergency{
				Department:         "72",
				PassageDate:        timeFn("2020-05-25"),
				Count:              67,
				Cov19SuspCount:     0,
				Cov19SuspHosp:      0,
				TotalSOSMedAct:     0,
				SOSMedCov19SuspAct: 0,
			},
			&covidtracker.Emergency{
				Department:         "976",
				PassageDate:        timeFn("2020-05-25"),
				Count:              2,
				Cov19SuspCount:     0,
				Cov19SuspHosp:      0,
				TotalSOSMedAct:     0,
				SOSMedCov19SuspAct: 0,
			},
			&covidtracker.Emergency{
				Department:         "01",
				PassageDate:        timeFn("2020-02-24"),
				Count:              357,
				Cov19SuspCount:     0,
				Cov19SuspHosp:      0,
				TotalSOSMedAct:     0,
				SOSMedCov19SuspAct: 0,
			},
		}

		test.Compare(t, sos, expected, "unexpected cases")

	})
}
