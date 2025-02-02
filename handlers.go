package main

import (
	"fmt"
	"net/http"
)

func handlerReadiness(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("OK"))
}

func (cfg *apiConfig) handlerGetHits(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	fmt.Println("Yos")
	writer.WriteHeader(http.StatusOK)
	hits := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	writer.Write([]byte(hits))
}

func (cfg *apiConfig) handlerResetHits(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	writer.Write([]byte("OK"))
}
