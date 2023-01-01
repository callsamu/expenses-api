package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler(t *testing.T) {
	app, _ := newTestApplication(t)
	ts := newTestServer(app.routes())
	response := ts.request(t, http.MethodGet, "/v1/healthcheck", nil)

	var input struct {
		Status     string `json:"status"`
		SystemInfo struct {
			Version     string `json:"version"`
			Environment string `json:"environment"`
		} `json:"system_info"`
	}

	err := json.NewDecoder(response.Body).Decode(&input)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "available", input.Status)
	assert.Equal(t, version, input.SystemInfo.Version)
	assert.Equal(t, app.config.env, input.SystemInfo.Environment)

	contentType := response.Header.Get("Content-Type")
	assert.Equal(t, "application/json", contentType, "incorrect content-type header")
}
