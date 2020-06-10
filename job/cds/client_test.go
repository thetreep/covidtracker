package cds

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/thetreep/covidtracker"
	"github.com/thetreep/toolbox/test"
)

var count int

func init() {
	AgentDutyCode = "THE_TREEP_API"
	User = "cdsuser"
	Password = "superpassword"
}

func TestHotelsByPrefix(t *testing.T) {
	t.Run("Find hotels with success response", func(t *testing.T) {
		c, mux, _, teardown := setup(t)
		defer teardown()
		mux.HandleFunc("/Hotels", func(w http.ResponseWriter, r *http.Request) {
			checkAuth(t, r)
			checkURLParams(t, r, map[string]string{
				"agentDutyCode": AgentDutyCode,
				"prefix":        "ibis Budget Marseille Timone",
			})
			defer r.Body.Close()
			fmt.Fprint(w, fileContent("./testdata/HotelByPrefix-resp.json"))
		})
		hotels, err := c.HotelsByPrefix("ibis Budget Marseille Timone")
		if err != nil {
			t.Fatal("une erreur est survenue", err)
		}
		if got, want := len(hotels), 1; got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
		fmt.Println(hotels[0])
		exphotels := []*covidtracker.Hotel{
			&covidtracker.Hotel{
				Name:     "Ibis Budget Marseille Timone",
				Address:  "69, Boulevard Sakakini",
				City:     "Marseille",
				ZipCode:  "13005",
				Country:  "FR",
				ImageURL: "https://bookings.cdsgroupe.com/photos/Search/FR/ACC/251/ACC2514.jpg",
				SanitaryInfos: []string{
					"Enregistrement & Règlement en ligne",
					"Distanciation sociale & sens de circulation",
					"Formation des équipes internes aux mesures internes",
					"Horaires & Messages de nettoyages des chambres",
					"Renforcement du nettoyage du linge",
					"Procédures de nettoyage renforcées des points de contacts en chambre",
					"Port du masque par le personnel",
					"Mise à disposition de gel hydro-alcoolique",
					"Nettoyage renforcé des lieux de passage",
				},
				SanitaryNote: 7.0,
				SanitaryNorm: "Accor - All Safe",
			}}
		test.Compare(t, hotels[0], exphotels[0])
	})
}

func fileContent(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("cannot open file %q", filename))
	}
	return string(content)
}

func parseBody(body io.Reader) string {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		panic(fmt.Sprintf("cannot parse body %q", body))
	}
	return string(data)
}

func setup(t *testing.T) (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	// HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	apiHandler := http.NewServeMux()
	server := httptest.NewServer(apiHandler)
	apiHandler.Handle("/", mux)

	// client is configured to use test server.
	endpoint = server.URL
	client = newClient(nil)
	mux.HandleFunc("/Authenticate", func(w http.ResponseWriter, r *http.Request) {
		t.Helper()
		if got, want := r.Header.Get("Content-Type"), "application/json"; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
		expect := fileContent("./testdata/Authenticate-req.json")
		if expect == "" {
			t.Fatal("error during parsing request")
		}
		defer r.Body.Close()
		body := parseBody(r.Body)
		fmt.Println(string(fileContent("./testdata/Authenticate-resp.json")))
		if got, want := body, expect; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
		fmt.Fprint(w, fileContent("./testdata/Authenticate-resp.json"))
	})

	return client, mux, server.URL, server.Close
}
func checkAuth(t *testing.T, r *http.Request) {
	t.Helper()
	auth := r.Header.Get("Authorization")
	if len(auth) == 0 {
		t.Fatal("expect auth, got none")
	}
	if !strings.EqualFold(auth[:7], "bearer ") {
		t.Fatalf("expect Bearer auth, got %q", auth[:7])
	}
	if got, want := auth[7:], "my-cds-auth-token"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func checkURLParams(t *testing.T, r *http.Request, expectParams map[string]string) {
	t.Helper()
	params := r.URL.Query()
	if got, want := len(params), len(expectParams); got != want {
		t.Fatalf("got %d URL params (%#v), want %d", got, params, want)
	}
	for k, v := range expectParams {
		if got, want := params.Get(k), v; got != want {
			t.Fatalf("URL param %s: got %s, want %s", k, got, want)
		}
	}
}
