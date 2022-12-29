package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	log := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		logger: log,
		config: cfg,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)

	srv := http.Server{
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		Addr:         fmt.Sprintf(":%d", cfg.port),
	}

	app.logger.Printf("starting %s server at port :%d", cfg.env, cfg.port)
	err := srv.ListenAndServe()
	if err != nil {
		app.logger.Fatal(err)
	}
}
