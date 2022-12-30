package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler(t *testing.T) {
	response := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "/v1/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	app := newTestApplication(t)
	app.healthcheckHandler(response, request)

	var input struct {
		Status      string `json:"status"`
		Version     string `json:"version"`
		Environment string `json:"environment"`
	}

	err = json.NewDecoder(response.Body).Decode(&input)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "available", input.Status)
	assert.Equal(t, version, input.Version)
	assert.Equal(t, app.config.env, input.Environment)
}
