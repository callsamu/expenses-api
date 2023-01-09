package main

import (
	"fmt"
	"net/http"
	"time"
)

func (app *application) serve() error {

	srv := http.Server{
		Handler:      app.routes(),
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		Addr:         fmt.Sprintf(":%d", app.config.port),
	}

	app.logger.Printf("starting %s server at port :%d", app.config.env, app.config.port)
	err := srv.ListenAndServe()

	return err
}