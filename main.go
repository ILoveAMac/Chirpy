package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	server := &http.Server{
		Handler: mux,
	}
	server.Addr = ":8080"

	mux.Handle("/", http.FileServer(http.Dir(".")))

	// Start the server
	server.ListenAndServe()
}
