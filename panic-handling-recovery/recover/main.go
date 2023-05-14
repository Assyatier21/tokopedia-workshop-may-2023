package main

import (
	"fmt"
	"runtime/debug"
)

func firstFunction() {
	fmt.Println("First function called")
	secondFunction()
	fmt.Println("First function finished")
}

func secondFunction() {
	fmt.Println("Second function called")
	panic("Panic happens")
	fmt.Println("Second function finished")
}

func doRecover() {
	fmt.Println("Panic example in Go")
	defer func() {
		// Recover from panic to stop termination of the application
		if r := recover(); r != nil {
			fmt.Println("Message: ", r)
			debug.PrintStack()
		}

	}()
	firstFunction()
	fmt.Println("All process finished")
}

func main() {
	doRecover()
}
