package datagouv_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/thetreep/covidtracker/job/datagouv"
)

type DatagouvAPI struct{}

func DatagouvServers(t *testing.T) (api *httptest.Server, opendata *httptest.Server) {
	t.Helper()

	opendata = httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//try to serve file
			fmt.Fprint(w, fileContent("./testdata/"+r.URL.Path))
		}),
	)

	api = httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			file := "./testdata/"
			n := strings.Replace(r.URL.Path, "/api/1/datasets/", "", -1)
			switch n {
			case string(datagouv.EmergencyDataset):
				file += "emergency-dataset.json"
			case string(datagouv.HospDataset):
				file += "hospitalization-dataset.json"
			case string(datagouv.ScreeningDataset):
				file += "screening-dataset.json"
			case string(datagouv.IndicDataset):
				file += "indicator-dataset.json"
			default:
				t.Fatalf("unexpected path %q", r.URL.Path)
			}
			fmt.Fprint(w, strings.Replace(fileContent(file), "{{apiServerHost}}", opendata.URL, -1))
		}),
	)
	return
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
