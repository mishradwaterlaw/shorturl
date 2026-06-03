package main

import (
	"log"
	"net/http"

	"github.com/mishradwaterlaw/shorturl/handler"
	"github.com/mishradwaterlaw/shorturl/store"
)

func main() {
	memoryStore := store.NewMemoryStore()

	// make(chan string, 100) creates a buffered channel of type string.
	// The buffer capacity of 100 allows it to hold 100 messages simultaneously.
	// This acts as a decoupling queue, preventing HTTP handlers from slowing down if logging is slow.
	logCh := make(chan string, 100)

	// The `go` keyword spins off an independent thread of execution (goroutine).
	// `func() { ... }()` is an Anonymous Immediately Invoked Function Expression.
	// This creates and fires up our background worker routine alongside the main program flow.
	go func() {
		// `for msg := range logCh` creates a continuous loop that blocks and waits for data.
		// Every time a string is pushed into logCh from any HTTP request, this loop wakes up,
		// pulls the message out, executes the body, and goes back to waiting.
		for msg := range logCh {
			log.Println("REDIRECT:", msg) // Prints the message synchronously in the background.
		}
	}()

	// Inject the logging channel into the handler initialization step.
	h := handler.NewHandler(memoryStore, logCh)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", h.HandleShorten)
	mux.HandleFunc("GET /", h.HandleRedirect)

	log.Println("server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
