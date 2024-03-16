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

func (app *application) serve() error{
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	go func(){
		quit := make(chan os.Signal, 1)

		// instructs the go runtime to start sending signals to the specified types
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		//Update the log entry to say "shutting down server instead of "caught signal
		app.logger.PrintInfo("shutting down server", map[string]string{
			"signal": s.String(),
		})

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		// Call Shutdown() on our server, passing in the context we just made
		// Shutdown() will return nil if err then return error
		shutdownError <- srv.Shutdown(ctx)

	}()

	app.logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env": app.config.env,
	})

	err := srv.ListenAndServe()
	// when gracful shutdown happen then server will return ErrServerClosed message
	if !errors.Is(err, http.ErrServerClosed){
		return err
	}

	err = <-shutdownError
	if err != nil{
		return err
	}

	app.logger.PrintInfo("stopped server", map[string]string{
		"add": srv.Addr,
	})

	return nil
}
