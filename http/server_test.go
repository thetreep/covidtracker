package http_test

import (
	"bytes"
	gohttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/thetreep/covidtracker/http"
)

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
		assertBody(t, resp, "OK")
		assertStatus(t, resp, gohttp.StatusOK)
	})

	//TODO test middlewares + routing
}

func assertStatus(t *testing.T, resp *gohttp.Response, expect int) {
	t.Helper()
	if got, want := resp.StatusCode, expect; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}

func assertBody(t *testing.T, resp *gohttp.Response, expect string) {
	t.Helper()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if got, want := buf.String(), expect; got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}
