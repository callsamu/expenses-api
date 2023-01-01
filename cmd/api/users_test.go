package main

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/callsamu/pfapi/internal/data"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUsersHandler(t *testing.T) {
	app, mock := newTestApplication(t)
	tsrv := newTestServer(app.routes())

	cases := []struct {
		name       string
		username   string
		email      string
		password   string
		wantStatus int
	}{
		{
			name:       "Handles basic requests",
			username:   "foobar",
			email:      "foobar@example.com",
			password:   "password",
			wantStatus: http.StatusCreated,
		},
	}

	for _, ts := range cases {
		t.Run(ts.name, func(t *testing.T) {
			input := envelope{
				"name":     ts.username,
				"email":    ts.email,
				"password": ts.password,
			}
			inputJSON, err := json.Marshal(input)
			if err != nil {
				t.Fatal(err)
			}

			columns := sqlmock.NewRows([]string{"id", "version", "created_at"})
			mock.ExpectQuery("").WillReturnRows(columns.AddRow(1, 1, time.Now()))

			response := tsrv.request(t, http.MethodPost, "/v1/users/register", inputJSON)
			assert.Equal(t, ts.wantStatus, response.StatusCode)

			var output struct {
				User data.User `json:"user"`
			}
			err = json.NewDecoder(response.Body).Decode(&output)

			assert.EqualValues(t, 1, output.User.ID)
			assert.Equal(t, ts.username, output.User.Name)
			assert.Equal(t, ts.email, output.User.Email)
			assert.Equal(t, false, output.User.Activated)
		})
	}

	t.Run("returns bad request on invalid json processing", func(t *testing.T) {
		input := []byte(`{ "foo": "bar" "ha"}`)

		columns := sqlmock.NewRows([]string{"id", "version", "created_at"})
		mock.ExpectQuery("").WillReturnRows(columns.AddRow(1, 1, time.Now()))

		response := tsrv.request(t, http.MethodPost, "/v1/users/register", input)
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})

}
