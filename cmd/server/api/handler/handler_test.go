package handler

import (
	"github.com/labstack/echo"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_healthCheck(t *testing.T) {
	e := echo.New()
	e.GET("/healthcheck", healthCheck)

	ts := httptest.NewServer(e)
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
