package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
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

type testServer struct {
	*httptest.Server
}

func newTestServer(handler http.Handler) *testServer {
	return &testServer{httptest.NewServer(handler)}
}

func (ts *testServer) GET(t *testing.T, uri string) *http.Response {
	response, err := ts.Client().Get(ts.URL + uri)
	if err != nil {
		t.Fatal(err)
	}

	return response
}
