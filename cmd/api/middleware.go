package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/callsamu/expenses-api/internal/data"
	"github.com/callsamu/expenses-api/internal/validator"
)

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			r = app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]
		v := validator.New()
		if data.ValidateTokenPlaintext(v, token); !v.Valid() {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		user, err := app.models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
				return
			}
		}

		r = app.contextSetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		trustedOrigins := app.config.cors.trustedOrigins

		origin := r.Header.Get("Origin")
		requestAllowMethod := r.Header.Get("Access-Control-Request-Method")

		if len(trustedOrigins) == 1 && trustedOrigins[0] == "*" {
			w.Header().Set("Access-Control-Allow-Origin", "*")

			// Check if is preflight request
			if origin != "" && r.Method == http.MethodOptions && requestAllowMethod != "" {
				w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PATCH, PUT, DELETE")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
				w.WriteHeader(http.StatusOK)

				return
			}
		} else {
			w.Header().Set("Vary", "Origin")

			if origin != "" {
				for i := range trustedOrigins {
					if origin == trustedOrigins[i] {
						w.Header().Set("Access-Control-Allow-Origin", origin)

						// Check if is preflight request
						if r.Method == http.MethodOptions && requestAllowMethod != "" {
							w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PATCH, PUT, DELETE")
							w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
							w.WriteHeader(http.StatusOK)

							return
						}
					}
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}