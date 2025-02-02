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
	mux.HandleFunc("GET /healthz", handlerReadiness)
	mux.HandleFunc("GET /metrics", cfg.handlerGetHits)
	mux.HandleFunc("POST /reset", cfg.handlerResetHits)
	server := http.Server{Addr: ":8080", Handler: mux}
	server.ListenAndServe()
}
