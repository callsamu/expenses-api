package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	response := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/v1/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	app := newTestApplication()

	app.healthcheckHandler(response, request)
	body := response.Body.String()

	assertBodyContains(t, body, "available")
	assertBodyContains(t, body, app.config.env)
	assertBodyContains(t, body, version)
}
