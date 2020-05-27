package datagouv_test

import (
	"path"
	"testing"

	"github.com/kr/pretty"
	"github.com/thetreep/covidtracker/job/datagouv"
)

func TestRefreshCase(t *testing.T) {
	t.Run("file exist", func(t *testing.T) {
		assertRessourceExist(t, path.Join(datagouv.DatagouvBase, string(datagouv.CaseURL)))
	})

	t.Run("parsing", func(t *testing.T) {
		ts := DatagouvServer(t)
		defer ts.Close()

		cases, err := s.RefreshCase()
		if err != nil {
			t.Fatal(err)
		}
		if len(cases) == 0 {
			t.Fatal("unexpected empty cases")
		}
		pretty.Log(cases)
		t.Fatal("toto")
		// expected := []*covidtracker.Case{
		// 	{},
		// 	{},
		// 	{},
		// 	{},
		// }

	})
}
