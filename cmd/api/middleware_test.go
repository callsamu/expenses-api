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

func TestEnableCORS(t *testing.T) {
	app, _ := newTestApplication(t)

	cases := []struct {
		name           string
		origin         string
		trustedOrigins []string

		method        string
		requestMethod string

		wantHeader                 string
		wantAllowCORSUnsafeMethods bool
	}{
		{
			name:           "enables every origin when trustedOrigins contains a single wildcard",
			origin:         "https://bar.net",
			trustedOrigins: []string{"*"},
			method:         http.MethodGet,
			wantHeader:     "*",
		},
		{
			name:           "sends respective origin when it matches one of trustedOrigins",
			origin:         "https://bar.net",
			trustedOrigins: []string{"https://foo.com", "https://bar.net"},
			method:         http.MethodGet,
			wantHeader:     "https://bar.net",
		},
		{
			name:           "confirms preflight request when trustedOrigin is wildcard",
			origin:         "https://bar.net",
			trustedOrigins: []string{"*"},
			method:         http.MethodOptions,
			wantHeader:     "*",

			requestMethod:              "PATCH",
			wantAllowCORSUnsafeMethods: true,
		},
		{
			name:           "confirms preflight request when origin matches trustedOrigins",
			origin:         "https://bar.net",
			trustedOrigins: []string{"https://foo.com", "https://bar.net"},
			method:         http.MethodOptions,
			wantHeader:     "https://bar.net",

			requestMethod:              "PATCH",
			wantAllowCORSUnsafeMethods: true,
		},
	}

	for _, ts := range cases {
		t.Run(ts.name, func(t *testing.T) {
			app.config.cors.trustedOrigins = ts.trustedOrigins

			rr := httptest.NewRecorder()
			request, err := http.NewRequest(ts.method, "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			request.Header.Set("Origin", ts.origin)
			request.Header.Set("Access-Control-Request-Method", ts.requestMethod)

			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("OK"))
			})
			app.enableCORS(next).ServeHTTP(rr, request)
			response := rr.Result()

			header := response.Header.Get("Access-Control-Allow-Origin")
			assert.Equal(t, ts.wantHeader, header)

			if ts.method == http.MethodOptions {
				allowedMethods := response.Header.Get("Access-Control-Allow-Methods")
				assert.Equal(t, "OPTIONS, PATCH, PUT, DELETE", allowedMethods)

				allowedHeaders := response.Header.Get("Access-Control-Allow-Headers")
				assert.Equal(t, "Authorization, Content-Type", allowedHeaders)
			}

			assert.Equal(t, http.StatusOK, response.StatusCode)
		})

	}

}
