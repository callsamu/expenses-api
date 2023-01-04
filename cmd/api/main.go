package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/callsamu/expenses-api/internal/data"
	"github.com/callsamu/expenses-api/internal/mailer"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

type application struct {
	config config
	logger *log.Logger
	models data.Models
	mailer interface {
		Send(recipient string, template string, data any) error
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max idle time")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "", "SMTP Host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 2525, "SMTP Port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "", "SMTP Username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "", "SMTP Password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "", "SMTP Sender")

	flag.Parse()

	log := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := &application{
		logger: log,
		config: cfg,
		models: data.NewModels(db),
		mailer: mailer.New(
			cfg.smtp.host,
			cfg.smtp.port,
			cfg.smtp.username,
			cfg.smtp.password,
			cfg.smtp.sender,
		),
	}

	srv := http.Server{
		Handler:      app.routes(),
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		Addr:         fmt.Sprintf(":%d", cfg.port),
	}

	app.logger.Printf("starting %s server at port :%d", cfg.env, cfg.port)
	err = srv.ListenAndServe()
	if err != nil {
		app.logger.Fatal(err)
	}
}
