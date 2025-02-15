package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func panicHandleHTTP(command http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Defer the process of recovery
		defer func() {
			// Recover from panic to stop termination of the application
			if r := recover(); r != nil {
				fmt.Println("[PanicHandleHTTP] panic message: ", r)
				debug.PrintStack()
			}
			// TODO: setup recover function to recover from a panic
		}()

		// Execute HTTP function that has been wrapped
		command(w, r)
	}
}

func registerRoutes(server *http.Server) {
	// Create endpoint to test panic process and call HTTP wrapper function to wrap our process
	// TODO: call HTTP wrapper function using http.Handle()
	// TODO: write message to the client and execute panic to trigger termination of the application
	http.Handle("/", panicHandleHTTP(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World\n"))
		panic("panic testing on /")
	}))
	server.Handler = http.DefaultServeMux
}

func doHTTPPanicRecovery() {
	port := 8000
	httpServer := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	// Register our HTTP endpoint
	registerRoutes(httpServer)

	fmt.Println("HTTP server listening on port", port)
	err := httpServer.ListenAndServe()
	if err != nil {
		fmt.Println("error when ListenAndServe")
		return
	}
}

func main() {
	doHTTPPanicRecovery()
}
