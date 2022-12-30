package main

import (
	"io"
	"log"
	"testing"
)

func newTestApplication(t *testing.T) *application {
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
