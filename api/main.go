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

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "server encountered an error: %s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	fmt.Println("setting up server")

	fmt.Printf("PID: %d\n", os.Getpid())

	// Create a new serve mux to act as our route multiplexer
	mux := http.NewServeMux()

	// Create a healthcheck endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("healthcheck called")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"healthy"}`))
	})

	// Create a new http server
	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	// setting up graceful shutdown logic
	fmt.Println("setting up graceful shutdown")

	// create a channel to receive errors
	errChan := make(chan error)

	ctx, done := context.WithCancel(ctx)
	defer done()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

		// block until a signal is received
		fmt.Println("waiting for shutdown signal")
		<-sig
		fmt.Println("received shutdown signal")

		// create a context with a timeout to allow the process to be killed if the
		// timeout is reached
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		// shutting down server
		fmt.Println("shutting down server")
		if err := server.Shutdown(ctx); err != nil {
			errChan <- fmt.Errorf("[in main.run] failed to shutdown http server: %w", err)
			return
		}

		fmt.Println("server shutdown complete")

		done()
	}()

	// starting server
	fmt.Println("listening on localhost:8080")
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("[in main.run] failed to listen and serve: %w", err)
	}

	// waiting for done signal or error
	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		fmt.Println("exiting application")
		return nil
	}
}
