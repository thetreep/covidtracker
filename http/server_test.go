/*
	This file is part of covidtracker also known as EviteCovid .

    Copyright 2020 the Treep

    covdtracker is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    covidtracker is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with covidtracker.  If not, see <https://www.gnu.org/licenses/>.
*/

package http_test

import (
	"bytes"
	"encoding/json"
	gohttp "net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/thetreep/covidtracker"
	"github.com/thetreep/covidtracker/http"
	"github.com/thetreep/toolbox/test"
)

type ApiResponse struct {
	Err string
}

func TestServer(t *testing.T) {

	tServer := httptest.NewServer(http.NewServer().Handlers())
	defer func() {
		tServer.Close()
	}()

	t.Run("test ping", func(t *testing.T) {
		resp, err := gohttp.Get(tServer.URL + "/ping")
		if err != nil {
			t.Fatal(err)
		}
		assertString(t, resp, "OK")
		assertStatus(t, resp, gohttp.StatusOK)
	})

	t.Run("missing api secret", func(t *testing.T) {
		resp, err := gohttp.Get(tServer.URL + "/toto")
		if err != nil {
			t.Fatal(err)
		}
		assertJSON(t, resp, ApiResponse{Err: covidtracker.ErrMissingAPISecret.Error()})
		assertStatus(t, resp, gohttp.StatusUnauthorized)
	})

	prvSecret := os.Getenv("THETREEP_COVIDTRACKER_SECRET")
	defer os.Setenv("THETREEP_COVIDTRACKER_SECRET", prvSecret)

	os.Setenv("THETREEP_COVIDTRACKER_SECRET", "a.wonderfull.secret")

	t.Run("bad api secret", func(t *testing.T) {
		req, err := gohttp.NewRequest(gohttp.MethodGet, tServer.URL+"/toto", nil)
		req.Header.Add("api-secret", "toto")

		resp, err := gohttp.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		assertJSON(t, resp, ApiResponse{Err: covidtracker.ErrInvalidAPISecret.Error()})
		assertStatus(t, resp, gohttp.StatusUnauthorized)
	})

	t.Run("not found", func(t *testing.T) {
		req, err := gohttp.NewRequest(gohttp.MethodGet, tServer.URL+"/toto", nil)
		req.Header.Add("api-secret", "a.wonderfull.secret")

		resp, err := gohttp.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		assertStatus(t, resp, gohttp.StatusNotFound)
	})

	//TODO test middlewares + routing
}

func assertStatus(t *testing.T, resp *gohttp.Response, expect int) {
	t.Helper()
	if got, want := resp.StatusCode, expect; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}

func assertJSON(t *testing.T, resp *gohttp.Response, want ApiResponse) {
	t.Helper()

	var got ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	test.Compare(t, got, want, "unexpected json reponse")
}

func assertString(t *testing.T, resp *gohttp.Response, expect string) {
	t.Helper()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if got, want := buf.String(), expect; got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}
