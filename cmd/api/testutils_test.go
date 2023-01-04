package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/callsamu/expenses-api/internal/data"
)

type mockMailer struct {
	recipient string
	template  string
	data      any

	calls int
}

type mocks struct {
	mailer *mockMailer
	db     sqlmock.Sqlmock
}

func (m *mockMailer) Send(recipient string, template string, data any) error {
	m.recipient = recipient
	m.template = template
	m.data = data

	m.calls += 1

	return nil
}

func newTestApplication(t *testing.T) (*application, mocks) {
	cfg := config{
		port: 4000,
		env:  "testing",
	}

	log := log.New(os.Stdout, "", 0)

	db, mockDB, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	mockMailer := &mockMailer{}

	app := &application{
		logger: log,
		config: cfg,
		models: data.NewModels(db),
		mailer: mockMailer,
	}

	mocks := mocks{
		db:     mockDB,
		mailer: mockMailer,
	}

	return app, mocks
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
