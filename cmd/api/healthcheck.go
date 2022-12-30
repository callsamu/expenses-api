package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"version":     version,
		"status":      "available",
		"environment": app.config.env,
	}

	json, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, string(json))
}
