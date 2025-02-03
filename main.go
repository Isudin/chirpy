package main

import (
	"net/http"
	"sync/atomic"
)

func main() {
	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}
	mux := http.NewServeMux()
	handler := http.FileServer(http.Dir("."))
	mux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(handler)))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerGetHits)
	mux.HandleFunc("POST /admin/reset", cfg.handlerResetHits)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	server := http.Server{Addr: ":8080", Handler: mux}
	server.ListenAndServe()
}
