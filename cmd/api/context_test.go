package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/callsamu/expenses-api/internal/data"
	"github.com/callsamu/expenses-api/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestContextSetUser(t *testing.T) {
	app, _ := newTestApplication(t)
	request, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	user := mocks.MockActivatedUser
	request = app.contextSetUser(request, user)
	retrieved, ok := request.Context().Value(userContextKey).(*data.User)
	assert.True(t, ok)
	assert.Equal(t, user, retrieved)
}

func TestContextGetUser(t *testing.T) {
	app, _ := newTestApplication(t)

	t.Run("gets user", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		user := mocks.MockActivatedUser
		ctx := context.WithValue(request.Context(), userContextKey, user)
		request = request.WithContext(ctx)

		retrieved := app.contextGetUser(request)
		assert.Equal(t, user, retrieved)
	})

	t.Run("panics if there was no user called", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		assert.Panics(t, func() { app.contextGetUser(request) })
	})
}