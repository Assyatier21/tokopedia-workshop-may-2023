package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

const (
	port = 8000
)

var (
	httpServer = &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}
)

func main() {
	// Let us see HTTP serve without graceful shutdown
	doServeHTTPGraceful()

	// Now let us implement graceful shutdown.
	// Can delete doServeHTTP to make playground simpler
	// doServeHTTPGraceful()
}

func doServeHTTP() {
	fmt.Println("HTTP server listening on port", port)
	err := httpServer.ListenAndServe()
	if err != nil {
		fmt.Println("error when ListenAndServe")
		return
	}

	fmt.Println("Process cleanup...") // This won't get called
}

func doServeHTTPGraceful() {
	// TODO: create context and its cancel function
	ctx, cancel := context.WithCancel(context.Background())

	// TODO: setup SIGTERM listener
	go func() {
		// Listen for the termination signal
		stop := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

		// Block until termination signal received
		<-stop

		// Block until termination signal received
		cancel()
	}()

	// setup errgroup with context so we can listen to its cancellation
	eg, egCtx := errgroup.WithContext(ctx)

	// setup HTTP listener
	eg.Go(func() error {
		fmt.Println("HTTP server listening on port", port)
		err := httpServer.ListenAndServe()
		fmt.Println("HTTP server finish listening on port", port)
		return err
	})

	// TODO: setup HTTP graceful shutdown
	eg.Go(func() error {
		<-egCtx.Done()
		fmt.Println("Shutting down the system...")
		httpServer.Shutdown(ctx)
		fmt.Println("Graceful shutdown finished...")
		return nil
	})

	// wait for errgroup
	if err := eg.Wait(); err != nil {
		fmt.Printf("Exit reason: %s \n", err)
	}

	fmt.Println("Process cleanup...") // This should get called
}
