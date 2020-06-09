package datagouv_test

import (
	"context"
	"testing"
	"time"

	"github.com/thetreep/covidtracker"
	"github.com/thetreep/covidtracker/job/datagouv"
	"github.com/thetreep/toolbox/test"
)

func TestRefreshScreening(t *testing.T) {
	t.Run("parsing", func(t *testing.T) {
		api, server := DatagouvServers(t)
		defer func() {
			api.Close()
			server.Close()
		}()

		s := datagouv.Service{Ctx: context.Background(), BasePath: api.URL}

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
