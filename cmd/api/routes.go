package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	r := httprouter.New()

	r.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	r.NotFound = http.HandlerFunc(app.notFoundResponse)

	return r
}
