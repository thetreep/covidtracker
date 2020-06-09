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
		ts := DatagouvServer(t)
		defer ts.Close()

		s := datagouv.Service{Ctx: context.Background(), BasePath: ts.URL}

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
				Department:    1,
				NoticeDate:    timeFn("2020-03-10"),
				AgeGroup:      "0",
				Count:         0,
				PositiveCount: 0,
				PositiveRate:  0,
			},
			&covidtracker.Screening{
				Department:    3,
				NoticeDate:    timeFn("2020-03-24"),
				AgeGroup:      "C",
				Count:         0,
				PositiveCount: 0,
				PositiveRate:  0,
			},
			&covidtracker.Screening{
				Department:    92,
				NoticeDate:    timeFn("2020-04-27"),
				AgeGroup:      "E",
				Count:         47,
				PositiveCount: 7,
				PositiveRate:  13,
			},
			&covidtracker.Screening{
				Department:    95,
				NoticeDate:    timeFn("2020-03-26"),
				AgeGroup:      "D",
				Count:         18,
				PositiveCount: 13,
				PositiveRate:  13,
			},
			&covidtracker.Screening{
				Department:    976,
				NoticeDate:    timeFn("2020-05-23"),
				AgeGroup:      "E",
				Count:         0,
				PositiveCount: 0,
				PositiveRate:  0,
			},
		}
		test.Compare(t, scrs, expected, "unexpected screenings")
	})
}
