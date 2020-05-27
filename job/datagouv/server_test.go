package datagouv_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func DatagouvServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			file := "./testdata/"
			switch r.URL.Path {
			case EmergencyURL:
				file += "emergency.csv"
			case CaseURL:
				file += "case.csv"
			case HospitalizationURL:
				file += "hospitalization.csv"
			case ScreeningURL:
				file += "screening.csv"
			case IndicatorURL:
				file += "indicator.csv"
			default:
				t.Fatal("unexpected path %q", r.URL.Path)
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

func assertRessourcexist(t *testing.T, url string) {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		return t.Fatal(err)
	}
	defer resp.Body.Close()

	t.Logf("get %s returns %s", url, resp.Body.Status)

	if got, want := resp.StatusCode, http.StatusOK; resp.StatusCode != http.StatusOK {
		t.Fatal("unexpected status code : got %s want %s", got, want)
	}

}
