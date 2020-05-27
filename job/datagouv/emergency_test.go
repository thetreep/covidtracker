package datagouv_test

import (
	"context"
	"testing"

	"github.com/thetreep/covidtracker/datagouv"
)

func TestRefreshEmergency(t *testing.T) {
	s := datagouv.Service{ctx: context.Background()}
	t.Run("download file and parse", func(t *testing.T) {
		ems, err := s.RefreshEmergency()
		if err != nil {
			t.Fatal(err)
		}
		if len(ems) == 0 {
			t.Fatal("unexpected empty cases")
		}
	})
}
