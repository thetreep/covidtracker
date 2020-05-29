package datagouv_test

import (
	"context"
	"testing"
	"time"

	"github.com/thetreep/covidtracker"
	"github.com/thetreep/covidtracker/job/datagouv"
	"github.com/thetreep/toolbox/test"
)

func TestRefreshCase(t *testing.T) {
	t.Run("file exist", func(t *testing.T) {
		assertRessourceExist(t, datagouv.DatagouvBase+string(datagouv.CaseURL))
	})

	t.Run("parsing", func(t *testing.T) {
		ts := DatagouvServer(t)
		defer ts.Close()

		s := datagouv.Service{Ctx: context.Background(), BasePath: ts.URL}

		cases, err := s.RefreshCase()
		if err != nil {
			t.Fatal(err)
		}
		if len(cases) == 0 {
			t.Fatal("unexpected empty cases")
		}
		// pretty.Log(cases)

		timeFn := func(s string) time.Time {
			t, _ := time.Parse("2006-01-02", s)
			return t
		}

		expected := []*covidtracker.Case{
			&covidtracker.Case{
				Department:              67,
				NoticeDate:              timeFn("2020-05-10"),
				HospServiceCountRelated: 28,
			},
			&covidtracker.Case{
				ID:                      "",
				Department:              57,
				NoticeDate:              timeFn("2020-05-13"),
				HospServiceCountRelated: 30,
			},
			&covidtracker.Case{
				ID:                      "",
				Department:              85,
				NoticeDate:              timeFn("2020-05-16"),
				HospServiceCountRelated: 8,
			},
			&covidtracker.Case{
				ID:                      "",
				Department:              973,
				NoticeDate:              timeFn("2020-05-19"),
				HospServiceCountRelated: 3,
			},
		}

		test.Compare(t, cases, expected, "unexpected cases")

	})
}
