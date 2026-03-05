package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()

	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	fileServer := http.FileServer(http.Dir("."))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))

	mux.HandleFunc("GET /healthz", serverHealth)
	mux.HandleFunc("GET /metrics", apiCfg.getHits)
	mux.HandleFunc("POST /reset", apiCfg.resetHits)

	// Start the server
	server.ListenAndServe()
}

func serverHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) getHits(w http.ResponseWriter, r *http.Request) {
	// 1. Get the current value
	hits := cfg.fileserverHits.Load()

	// 2. Set the header (as requested in the instructions)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// 3. Use fmt.Fprintf to format the string "Hits: x" directly into the writer
	fmt.Fprintf(w, "Hits: %d", hits)
}

func (cfg *apiConfig) resetHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
	cfg.fileserverHits.Store(0)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
