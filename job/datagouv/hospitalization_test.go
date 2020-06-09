package datagouv_test

import (
	"context"
	"testing"
	"time"

	"github.com/thetreep/covidtracker"
	"github.com/thetreep/covidtracker/job/datagouv"
	"github.com/thetreep/toolbox/test"
)

func TestRefreshHospitalization(t *testing.T) {
	t.Run("parsing", func(t *testing.T) {
		ts := DatagouvServer(t)
		defer ts.Close()

		s := datagouv.Service{Ctx: context.Background(), BasePath: ts.URL}

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
				ID:              "",
				Department:      976,
				Date:            timeFn("2020-05-26"),
				Count:           26,
				CriticalCount:   6,
				ReturnHomeCount: 106,
				DeathCount:      6,
			},
			&covidtracker.Hospitalization{
				ID:              "",
				Department:      22,
				Date:            timeFn("2020-05-20"),
				Count:           24,
				CriticalCount:   2,
				ReturnHomeCount: 81,
				DeathCount:      18,
			},
			&covidtracker.Hospitalization{
				ID:              "",
				Department:      42,
				Date:            timeFn("2020-05-17"),
				Count:           127,
				CriticalCount:   5,
				ReturnHomeCount: 374,
				DeathCount:      67,
			},
			&covidtracker.Hospitalization{
				ID:              "",
				Department:      1,
				Date:            timeFn("2020-03-18"),
				Count:           13,
				CriticalCount:   0,
				ReturnHomeCount: 2,
				DeathCount:      0,
			},
			&covidtracker.Hospitalization{
				ID:              "",
				Department:      93,
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
