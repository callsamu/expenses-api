package main

import (
	"io"
	"log"
	"strings"
	"testing"
)

func newTestApplication() *application {
	cfg := config{
		port: 4001,
		env:  "testing",
	}

	log := log.New(io.Discard, "", 0)

	app := &application{
		logger: log,
		config: cfg,
	}

	return app
}

func assertBodyContains(t *testing.T, body, substr string) {
	if strings.Contains(body, substr) {
		t.Errorf("want body to contain \"%s\"", body)
	}
}