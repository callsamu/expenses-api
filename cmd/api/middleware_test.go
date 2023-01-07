package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/callsamu/expenses-api/internal/data"
	"github.com/callsamu/expenses-api/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthenticateMiddleware(t *testing.T) {

	validToken := mocks.MockAuthenticationToken.Plaintext
	invalidToken := mocks.MockPlaintext

	cases := []struct {
		name              string
		header            string
		wantStatus        int
		wantAnonymousUser bool
	}{
		{
			name:              "token is valid and user is authenticated",
			header:            "Bearer" + " " + validToken,
			wantStatus:        http.StatusOK,
			wantAnonymousUser: false,
		},
		{
			name:       "header is malformed",
			header:     "Bear akak",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "token is malformed",
			header:     "Bearer" + " " + "1231230809",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "token is invalid",
			header:     "Bearer" + " " + invalidToken,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:              "header is not set and user is anonymous",
			header:            "",
			wantStatus:        http.StatusOK,
			wantAnonymousUser: true,
		},
	}

	for _, ts := range cases {
		t.Run(ts.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			request, err := http.NewRequest("", "", nil)
			if err != nil {
				t.Fatal(err)
			}

			request.Header.Set("Authorization", ts.header)

			var user *data.User
			app, _ := newTestApplication(t)

			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.NotPanics(t, func() {
					user = app.contextGetUser(r)
				})
				w.Write([]byte("OK"))
			})
			app.authenticate(next).ServeHTTP(rr, request)

			response := rr.Result()

			if ts.wantStatus == http.StatusOK {
				require.NotNil(t, user, ts.name)

				if ts.wantAnonymousUser {
					assert.True(t, user.IsAnonymous(), "authenticated user should be anonymous")
				} else {
					assert.False(t, user.IsAnonymous(), "authenticated user should not be anonymous")
				}

				body, err := io.ReadAll(response.Body)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, body, []byte("OK"))
			}

			assert.Equal(t, ts.wantStatus, response.StatusCode)
		})
	}

}
