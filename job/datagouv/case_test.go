package datagouv_test

import (
	"context"
	"testing"

<<<<<<< Updated upstream
	"github.com/thetreep/covidtracker"
=======
>>>>>>> Stashed changes
	"github.com/thetreep/covidtracker/datagouv"
)

func TestRefreshCase(t *testing.T) {
<<<<<<< Updated upstream
	t.Run("file exist and format ok", func(t *testing.T) {
		s := datagouv.Service{ctx: context.Background(), BasePath: datagouv.Datagouv}
	})

	t.Run("parsing", func(t *testing.T) {
		ts := DatagouvServer(t)
		defer ts.Close()

=======
	s := datagouv.Service{ctx: context.Background()}
	t.Run("file exist and format ok", func(t *testing.T) {
		cases, err := s.RefreshCase()
		if err != nil {
			t.Fatal(err)
		}
		if len(cases) == 0 {
			t.Fatal("unexpected empty cases")
		}
	})

	t.Run("parsing", func(t *testing.T) {
>>>>>>> Stashed changes
		cases, err := s.RefreshCase()
		if err != nil {
			t.Fatal(err)
		}
		if len(cases) == 0 {
			t.Fatal("unexpected empty cases")
		}
<<<<<<< Updated upstream

		expected := []*covidtracker.Case{
			{},
			{},
			{},
			{},
		}

=======
>>>>>>> Stashed changes
	})
}
