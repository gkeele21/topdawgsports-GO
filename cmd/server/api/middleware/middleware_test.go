package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MordFustang21/nova"
)

func Test_middleware(t *testing.T) {
	m := nova.New()
	Register(m)

	m.Get("/test", func(req *nova.Request) error {
		return nil
	})

	ts := httptest.NewServer(m)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/test")
	if err != nil {
		t.Error("couldn't make request")
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Error("Content-Type not set")
	}

	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Error("Access-Control-Allow-Origin not set")
	}
}
