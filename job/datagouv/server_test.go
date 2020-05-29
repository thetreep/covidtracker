package datagouv_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thetreep/covidtracker/job/datagouv"
)

func DatagouvServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			file := "./testdata/"
			switch r.URL.Path {
			case string(datagouv.EmergencyURL):
				file += "emergency.csv"
			case string(datagouv.CaseURL):
				file += "case.csv"
			case string(datagouv.HospitalizationURL):
				file += "hospitalization.csv"
			case string(datagouv.ScreeningURL):
				file += "screening.csv"
			case string(datagouv.IndicatorURL):
				file += "indicator.csv"
			default:
				t.Fatalf("unexpected path %q", r.URL.Path)
			}
			fmt.Fprint(w, fileContent(file))
		}),
	)
}

func fileContent(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("cannot open file %q", filename))
	}
	return string(content)
}

func assertRessourceExist(t *testing.T, url string) {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	t.Logf("get %s returns %s", url, resp.Status)

	if got, want := resp.StatusCode, http.StatusOK; resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code : got %d want %d", got, want)
	}
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
