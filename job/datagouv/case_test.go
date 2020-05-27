package datagouv_test

import (
	"context"
	"testing"

	"github.com/thetreep/covidtracker"
	"github.com/thetreep/covidtracker/datagouv"
)

func TestRefreshCase(t *testing.T) {
	s := datagouv.Service{ctx: context.Background()}
	t.Run("file exist and format ok", func(t *testing.T) {
		s.BasePath = datagouv.Datagouv
		cases, err := s.RefreshCase()
		if err != nil {
			t.Fatal(err)
		}
		if len(cases) == 0 {
			t.Fatal("unexpected empty cases")
		}
	})

	t.Run("parsing", func(t *testing.T) {
		s.BasePath = datagouv.Datagouv

		ts := DatagouvServer(t)
		defer ts.Close()

		cases, err := s.RefreshCase()
		if err != nil {
			t.Fatal(err)
		}
		if len(cases) == 0 {
			t.Fatal("unexpected empty cases")
		}

		expected := []*covidtracker.Case{
			{},
		}
	})
}
