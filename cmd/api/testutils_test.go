package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/callsamu/expenses-api/internal/data"
	"github.com/callsamu/expenses-api/internal/mocks"
)

type mockMailer struct {
	recipient string
	template  string
	data      any

	calls int
}

type mock struct {
	mailer *mockMailer
	users  *mocks.UserModel
	tokens *mocks.TokenModel
}

func (m *mockMailer) Send(recipient string, template string, data any) error {
	m.recipient = recipient
	m.template = template
	m.data = data

	m.calls += 1

	return nil
}

func newTestApplication(t *testing.T) (*application, mock) {
	cfg := config{
		port: 4000,
		env:  "testing",
	}

	log := log.New(os.Stdout, "", 0)

	mockMailer := &mockMailer{}
	mockUsers := &mocks.UserModel{}
	mockTokens := &mocks.TokenModel{}

	app := &application{
		logger: log,
		config: cfg,
		mailer: mockMailer,
		models: data.Models{
			Users:  mockUsers,
			Tokens: mockTokens,
		},
	}

	mocks := mock{
		mailer: mockMailer,
		users:  mockUsers,
		tokens: mockTokens,
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
