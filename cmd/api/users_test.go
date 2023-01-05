package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/callsamu/expenses-api/internal/data"
	"github.com/callsamu/expenses-api/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterUsersHandler(t *testing.T) {
	app, mock := newTestApplication(t)
	tsrv := newTestServer(app.routes())

	t.Run("handles basic requests", func(t *testing.T) {
		input := map[string]string{
			"name":     "foobar",
			"email":    "foobar@example.com",
			"password": "my password",
		}
		inputJSON, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}

		response := tsrv.request(t, http.MethodPost, "/v1/users/register", inputJSON)
		require.Equal(t, http.StatusCreated, response.StatusCode)

		mailData, ok := mock.mailer.data.(map[string]any)
		require.NotNil(t, mailData)
		require.True(t, ok)

		assert.Equal(t, mocks.MockPlaintext, mailData["Token"])

		var output struct {
			User data.User `json:"user"`
		}
		err = json.NewDecoder(response.Body).Decode(&output)
		if err != nil {
			t.Fatal(err)
		}

		assert.EqualValues(t, 3, output.User.ID)
		assert.Equal(t, "foobar", output.User.Name)
		assert.False(t, output.User.Activated)
	})

	t.Run("validates requests", func(t *testing.T) {
		input := map[string]string{
			"name":     "ha",
			"email":    "arexample.com",
			"password": "short",
		}
		inputJSON, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}

		response := tsrv.request(t, http.MethodPost, "/v1/users/register", inputJSON)
		require.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)
	})

	t.Run("returns bad request on invalid json processing", func(t *testing.T) {
		input := []byte(`{ "name" "ha" "email": "ha@example.com"}`)
		response := tsrv.request(t, http.MethodPost, "/v1/users/register", input)
		require.Equal(t, http.StatusBadRequest, response.StatusCode)
	})

	t.Run("returns validation failed if email is duplicated", func(t *testing.T) {
		input := map[string]string{
			"name":     "foo",
			"email":    "foo@example.com",
			"password": "mypassword",
		}
		JSON, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}
		response := tsrv.request(t, http.MethodPost, "/v1/users/register", JSON)
		require.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)
	})
}

func TestActivateUsersHandler(t *testing.T) {
	app, _ := newTestApplication(t)
	tsrv := newTestServer(app.routes())

	t.Run("activates user when token is found", func(t *testing.T) {
		input := map[string]string{
			"token": mocks.MockActivationToken.Plaintext,
		}
		inputJSON, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}

		response := tsrv.request(t, http.MethodPut, "/v1/users/activated", inputJSON)
		require.Equal(t, http.StatusOK, response.StatusCode)

		var output struct {
			User data.User `json:"user"`
		}
		err = json.NewDecoder(response.Body).Decode(&output)
		if err != nil {
			t.Fatal(err)
		}
		assert.EqualValues(t, 1, output.User.ID)
		assert.True(t, output.User.Activated)

	})
	t.Run("return validation error response when token is not found", func(t *testing.T) {
		input := map[string]string{
			"token": mocks.MockPlaintext,
		}
		inputJSON, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}

		response := tsrv.request(t, http.MethodPut, "/v1/users/activated", inputJSON)
		require.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)
	})
	t.Run("return validation error response when token is invalid", func(t *testing.T) {
		input := map[string]string{
			"token": "nqwejqwej",
		}
		inputJSON, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}

		response := tsrv.request(t, http.MethodPut, "/v1/users/activated", inputJSON)
		require.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)
	})
}
