package main

import (
	"io"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/callsamu/pfapi/internal/data"
)

func newTestApplication(t *testing.T) (*application, *sqlmock.Sqlmock) {
	cfg := config{
		port: 4000,
		env:  "testing",
	}

	log := log.New(io.Discard, "", 0)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	app := &application{
		logger: log,
		config: cfg,
		models: data.NewModels(db),
	}

	return app, &mock
}
