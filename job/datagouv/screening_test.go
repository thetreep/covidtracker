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

func TestRefreshScreening(t *testing.T) {
	t.Run("parsing", func(t *testing.T) {
		api, server := DatagouvServers(t)
		defer func() {
			api.Close()
			server.Close()
		}()

		s := datagouv.NewService(context.Background(), &logger.Logger{})
		s.BasePath = api.URL

		scrs, err := s.RefreshScreening()
		if err != nil {
			t.Fatal(err)
		}
		if len(scrs) == 0 {
			t.Fatal("unexpected empty screening")
		}

		timeFn := func(s string) time.Time {
			t, _ := time.Parse("2006-01-02", s)
			return t
		}

		expected := []*covidtracker.Screening{
			&covidtracker.Screening{
				Department:    "976",
				NoticeDate:    timeFn("2020-05-23"),
				Count:         0,
				PositiveCount: 0,
				PositiveRate:  0,
			},
			&covidtracker.Screening{
				Department:    "92",
				NoticeDate:    timeFn("2020-04-27"),
				Count:         47,
				PositiveCount: 7,
				PositiveRate:  13,
			},
			&covidtracker.Screening{
				Department:    "95",
				NoticeDate:    timeFn("2020-03-26"),
				Count:         28,
				PositiveCount: 24,
				PositiveRate:  16,
			},
			&covidtracker.Screening{
				Department:    "03",
				NoticeDate:    timeFn("2020-03-24"),
				Count:         0,
				PositiveCount: 0,
				PositiveRate:  0,
			},
			&covidtracker.Screening{
				Department:    "01",
				NoticeDate:    timeFn("2020-03-10"),
				Count:         12,
				PositiveCount: 1,
				PositiveRate:  1,
			},
		}
		test.Compare(t, scrs, expected, "unexpected screenings")
	})
}
