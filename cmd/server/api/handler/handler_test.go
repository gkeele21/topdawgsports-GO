package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MordFustang21/nova"
)

func Test_healthCheck(t *testing.T) {
	m := nova.New()
	m.Get("/healthcheck", healthCheck)

	ts := httptest.NewServer(m)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/healthcheck")
	if err != nil {
		t.Errorf("couldn't make request: %s", err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		t.Errorf("got error response status %d", res.StatusCode)
	}
}
