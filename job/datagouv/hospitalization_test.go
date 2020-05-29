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
	t.Run("file exist", func(t *testing.T) {
		assertRessourceExist(t, datagouv.DatagouvBase+string(datagouv.HospitalizationURL))
	})

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
				Department:      1,
				NoticeDate:      timeFn("2020-03-18"),
				Sex:             "0",
				Count:           2,
				CriticalCount:   0,
				ReturnHomeCount: 1,
				DeathCount:      0,
			},
			&covidtracker.Hospitalization{
				ID:              "",
				Department:      93,
				NoticeDate:      timeFn("2020-03-18"),
				Sex:             "0",
				Count:           92,
				CriticalCount:   33,
				ReturnHomeCount: 15,
				DeathCount:      5,
			},
			&covidtracker.Hospitalization{
				ID:              "",
				Department:      42,
				NoticeDate:      timeFn("2020-05-17"),
				Sex:             "2",
				Count:           127,
				CriticalCount:   5,
				ReturnHomeCount: 374,
				DeathCount:      67,
			},
			&covidtracker.Hospitalization{
				ID:              "",
				Department:      22,
				NoticeDate:      timeFn("2020-05-20"),
				Sex:             "1",
				Count:           24,
				CriticalCount:   2,
				ReturnHomeCount: 81,
				DeathCount:      18,
			},
			&covidtracker.Hospitalization{
				ID:              "",
				Department:      976,
				NoticeDate:      timeFn("2020-05-26"),
				Sex:             "2",
				Count:           26,
				CriticalCount:   6,
				ReturnHomeCount: 106,
				DeathCount:      6,
			},
		}

		test.Compare(t, hosp, expected, "unexpected hospitalization")

	})
}
