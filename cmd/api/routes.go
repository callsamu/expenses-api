package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	r := httprouter.New()

	r.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	r.HandlerFunc(http.MethodPost, "/v1/users/register", app.registerUserHandler)
	r.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	r.HandlerFunc(http.MethodPost, "/v1/users/authentication", app.sendAuthenticationTokenHandler)

	r.NotFound = http.HandlerFunc(app.notFoundResponse)

	return r
}
