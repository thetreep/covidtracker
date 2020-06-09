package datagouv_test

import (
	"context"
	"testing"
	"time"

	"github.com/thetreep/covidtracker"
	"github.com/thetreep/covidtracker/job/datagouv"
	"github.com/thetreep/toolbox/test"
)

func TestRefreshIndicator(t *testing.T) {
	t.Run("parsing", func(t *testing.T) {
		ts := DatagouvServer(t)
		defer ts.Close()

		s := datagouv.Service{Ctx: context.Background(), BasePath: ts.URL}

		inds, err := s.RefreshIndicator()
		if err != nil {
			t.Fatal(err)
		}
		if len(inds) == 0 {
			t.Fatal("unexpected empty indicators")
		}

		timeFn := func(s string) time.Time {
			t, _ := time.Parse("2006-01-02", s)
			return t
		}

		expected := []*covidtracker.Indicator{
			{
				ExtractDate: timeFn("2020-05-07"),
				Department:  59,
				Color:       "rouge",
			},
			{
				ExtractDate: timeFn("2020-05-07"),
				Department:  78,
				Color:       "rouge",
			},
			{
				ExtractDate: timeFn("2020-05-04"),
				Department:  88,
				Color:       "rouge",
			},
			{
				ExtractDate: timeFn("2020-05-03"),
				Department:  15,
				Color:       "orange",
			},
		}

		test.Compare(t, inds, expected, "unexpected indicator")

	})
}
