package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	// Let us see process without graceful will behave
	doProcessGraceful()

	// Then, let us work on graceful implementation
	// Can delete doProcess function to make playground simpler
	// doProcessGraceful()
}

func doProcess() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			time.Sleep(1 * time.Second)
			fmt.Println("Hello in the first loop")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			time.Sleep(1 * time.Second)
			fmt.Println("Hello in the second loop")
		}
	}()

	wg.Wait()
	fmt.Println("Process cleanup...") // This won't get called
}

func doProcessGraceful() {
	// TODO: setup context and its cancel function
	ctx, cancel := context.WithCancel(context.Background())

	// TODO: setup SIGTERM listener
	go func() {
		// Listen for the termination signal
		stop := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

		// Block until termination signal received
		<-stop

		cancel()
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		// TODO: convert into select case syntax and listen to context cancellation
		for {
			select {
			case <-ctx.Done(): // Block until cancel() is called
				fmt.Println("Break the first loop")
				return
			case <-time.After(1 * time.Second):
				fmt.Println("Hello in the first loop")
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		// TODO: convert into select case syntax and listen to context cancellation
		for {
			select {
			case <-ctx.Done(): // Block until cancel() is called
				fmt.Println("Break the second loop")
				return
			case <-time.After(1 * time.Second):
				fmt.Println("Hello in the second loop")
			}
		}
	}()

	// Wait for ongoing process to finish
	wg.Wait()
	fmt.Println("Process cleanup...") // This should get called
}
