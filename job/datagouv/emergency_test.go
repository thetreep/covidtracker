package datagouv_test

import (
	"context"
	"testing"
	"time"

	"github.com/thetreep/covidtracker"
	"github.com/thetreep/covidtracker/job/datagouv"
	"github.com/thetreep/toolbox/test"
)

func TestRefreshEmergency(t *testing.T) {
	t.Run("parsing", func(t *testing.T) {
		ts := DatagouvServer(t)
		defer ts.Close()

		s := datagouv.Service{Ctx: context.Background(), BasePath: ts.URL}

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
				Department:         1,
				NoticeDate:         timeFn("2020-02-24"),
				AgeGroup:           "0",
				Count:              357,
				Cov19SuspCount:     0,
				Cov19SuspHosp:      0,
				TotalSOSMedAct:     0,
				SOSMedCov19SuspAct: 0,
			},
			&covidtracker.Emergency{
				Department:         62,
				NoticeDate:         timeFn("2020-05-25"),
				AgeGroup:           "B",
				Count:              113,
				Cov19SuspCount:     1,
				Cov19SuspHosp:      0,
				TotalSOSMedAct:     0,
				SOSMedCov19SuspAct: 0,
			},
			&covidtracker.Emergency{
				Department:         72,
				NoticeDate:         timeFn("2020-05-25"),
				AgeGroup:           "A",
				Count:              67,
				Cov19SuspCount:     0,
				Cov19SuspHosp:      0,
				TotalSOSMedAct:     0,
				SOSMedCov19SuspAct: 0,
			},
			&covidtracker.Emergency{
				Department:         976,
				NoticeDate:         timeFn("2020-05-25"),
				AgeGroup:           "E",
				Count:              2,
				Cov19SuspCount:     0,
				Cov19SuspHosp:      0,
				TotalSOSMedAct:     0,
				SOSMedCov19SuspAct: 0,
			},
		}

		test.Compare(t, sos, expected, "unexpected cases")

	})
}
