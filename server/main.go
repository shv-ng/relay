package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func run() error {
	// Start server on PORT provided by default its 8000
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// GET /health endpoint return empty body with status code 200
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {})

	// GET / endpoint say hello world
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World")
	})

	return http.ListenAndServe(":"+port, nil)
}
