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

func (ts *testServer) request(t *testing.T, method string, url string, body []byte) *http.Response {
	request, err := http.NewRequest(method, ts.URL+url, bytes.NewBuffer(body))
	response, err := ts.Client().Do(request)
	if err != nil {
		t.Fatal(err)
	}

	return response
}
