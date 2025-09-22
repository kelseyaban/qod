package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
	"os"
	"errors"
	"os/signal"
	"syscall"
	"context"
)

func (app *application) serve() error {

	srv := &http.Server {
	Addr: fmt.Sprintf(":%d", app.config.port),
	Handler:  app.routes(),
	IdleTimeout: time.Minute,
	ReadTimeout: 5 * time.Second,
	WriteTimeout: 10 * time.Second,
	ErrorLog: slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}
	//channel to keep track of any errors during the shutdown process
	shutdownError := make(chan error)
	// create a goroutine that runs in the background listening for the shutdown signals
	  go func() {
		 quit := make(chan os.Signal, 1)  // receive the shutdown signal
		 signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // signal occurred
		 s := <-quit   // blocks until a signal is received
		 // message about shutdown in process
		 app.logger.Info("shutting down server", "signal", s.String())
		// create a context
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// initiate the shutdown. If all okay this returns nil
		shutdownError <- srv.Shutdown(ctx)
		}()
	
   	app.logger.Info("starting server", "addr", srv.Addr, "env", app.config.env)

	   err := srv.ListenAndServe()
	   if !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		// check the error channel to see if there were shutdown errors
		err = <-shutdownError
		if err != nil {
			 return err
	   }
	   // graceful shutdown was successful
	   app.logger.Info("stopped server", "address", srv.Addr)
	 
	   return nil

}