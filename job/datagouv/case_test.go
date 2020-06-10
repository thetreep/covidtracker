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

func TestRefreshCase(t *testing.T) {

	t.Run("parsing", func(t *testing.T) {
		api, server := DatagouvServers(t)
		defer func() {
			api.Close()
			server.Close()
		}()

		s := datagouv.NewService(context.Background(), &logger.Logger{})
		s.BasePath = api.URL

		cases, err := s.RefreshCase()
		if err != nil {
			t.Fatal(err)
		}
		if len(cases) == 0 {
			t.Fatal("unexpected empty cases")
		}

		timeFn := func(s string) time.Time {
			t, _ := time.Parse("2006-01-02", s)
			return t
		}

		expected := []*covidtracker.Case{
			&covidtracker.Case{
				ID:                      "",
				Department:              "973",
				NoticeDate:              timeFn("2020-05-19"),
				HospServiceCountRelated: 3,
			},
			&covidtracker.Case{
				ID:                      "",
				Department:              "85",
				NoticeDate:              timeFn("2020-05-16"),
				HospServiceCountRelated: 8,
			},
			&covidtracker.Case{
				ID:                      "",
				Department:              "57",
				NoticeDate:              timeFn("2020-05-13"),
				HospServiceCountRelated: 30,
			},
			&covidtracker.Case{
				Department:              "67",
				NoticeDate:              timeFn("2020-05-10"),
				HospServiceCountRelated: 28,
			},
		}

		test.Compare(t, cases, expected, "unexpected cases")

	})
}
