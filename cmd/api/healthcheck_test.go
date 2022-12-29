package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	response := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/v1/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	cfg := config{
		port: 4001,
		env:  "testing",
	}

	log := log.New(io.Discard, "", 0)

	app := &application{
		logger: log,
		config: cfg,
	}

	app.healthcheckHandler(response, request)
	body := response.Body.String()

	want := "available"
	if !strings.Contains(body, want) {
		t.Errorf("want body to contain \"%s\"", want)
	}

	want = cfg.env
	if !strings.Contains(body, want) {
		t.Errorf("want body to contain \"%s\"", want)
	}

	want = version
	if !strings.Contains(body, want) {
		t.Errorf("want body to contain \"%s\"", want)
	}
}
