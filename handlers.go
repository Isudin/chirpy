package main

import (
	"context"
	"fmt"
	"net/http"
)

func handlerReadiness(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("OK"))
}

func (cfg *apiConfig) handlerGetHits(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	hits := fmt.Sprintf(`
	<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>`, cfg.fileserverHits.Load())
	writer.Write([]byte(hits))
}

func (cfg *apiConfig) handlerReset(writer http.ResponseWriter, request *http.Request) {
	if cfg.platform != "dev" {
		writer.WriteHeader(403)
		return
	}

	cfg.fileserverHits.Store(0)
	err := cfg.queries.DeleteAllUsers(context.Background())
	if err != nil {
		respondError(writer, http.StatusInternalServerError, "Error deleting users", err)
		return
	}

	err = cfg.queries.DeleteAllChirps(context.Background())
	if err != nil {
		respondError(writer, http.StatusInternalServerError, "Error deleting chirps", err)
		return
	}

	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("OK"))
}
