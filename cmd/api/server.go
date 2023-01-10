package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := http.Server{
		Handler:      app.routes(),
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		Addr:         fmt.Sprintf(":%d", app.config.port),
	}

	shutdownError := make(chan error, 1)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit
		app.logger.Info().
			Str("signal", s.String()).
			Msg("shutting down server")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
			return
		}

		app.wg.Wait()
		shutdownError <- nil
	}()

	app.logger.Info().
		Str("environment", app.config.env).
		Int("port", app.config.port).
		Msg("starting server")

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Info().
		Str("address", srv.Addr).
		Msg("server stopped")

	return nil
}