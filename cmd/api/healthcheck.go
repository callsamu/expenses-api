package main

import (
	"log"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"version":     version,
		"status":      "available",
		"environment": app.config.env,
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		log.Fatal(err)
	}
}
