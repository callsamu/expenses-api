package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/callsamu/expenses-api/internal/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			password:   "mypassword",
			wantStatus: http.StatusCreated,
		},
		{
			name:       "Validates requests",
			username:   "",
			email:      "invalid@example.com",
			password:   "tooshort",
			wantStatus: http.StatusUnprocessableEntity,
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
			require.Equal(t, ts.wantStatus, response.StatusCode)
			if ts.wantStatus != http.StatusCreated {
				return
			}

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

func TestActivateUsersHandler(t *testing.T) {
	app, mock := newTestApplication(t)
	tsrv := newTestServer(app.routes())

	t.Run("activates user when token is found", func(t *testing.T) {
		input := map[string]string{
			"token": "abcdefghijklmnopqrstuvwxyz",
		}

		inputJSON, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}

		columns := sqlmock.NewRows([]string{"id", "created_at", "name", "email", "password_hash", "activated", "version"})

		password := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

		mock.ExpectQuery("SELECT").WillReturnRows(columns.AddRow(1, time.Now(), "foo", "foo@example.com", password, false, 1))
		mock.ExpectQuery("UPDATE users").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow(1))
		mock.ExpectExec("DELETE FROM tokens").WillReturnResult(sqlmock.NewResult(1, 1))

		t.Log(string(inputJSON))
		response := tsrv.request(t, http.MethodPut, "/v1/users/activated", inputJSON)
		require.Equal(t, http.StatusOK, response.StatusCode)

		var output struct {
			User data.User `json:"user"`
		}
		err = json.NewDecoder(response.Body).Decode(&output)

		assert.EqualValues(t, 1, output.User.ID)
		assert.True(t, output.User.Activated)

		if err = mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})

	t.Run("return validation error response when tokens is not found", func(t *testing.T) {
		input := map[string]string{
			"token": "abcdefghijklmnopqrstuvwxyz",
		}

		inputJSON, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}

		mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)
		response := tsrv.request(t, http.MethodPut, "/v1/users/activated", inputJSON)
		require.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)
		if err = mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})

	t.Run("return validation error response when tokens is invalid", func(t *testing.T) {
		input := map[string]string{
			"token": "too short",
		}

		inputJSON, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}

		response := tsrv.request(t, http.MethodPut, "/v1/users/activated", inputJSON)
		require.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)
	})
}
