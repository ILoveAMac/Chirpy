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

	fileServer := http.FileServer(http.Dir("."))
	mux.Handle("/app/", http.StripPrefix("/app", fileServer))

	mux.HandleFunc("/healthz", serverHealth)

	// Start the server
	server.ListenAndServe()
}

func serverHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
