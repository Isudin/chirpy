package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"

	"github.com/Isudin/chirpy/internal/database"
	"github.com/google/uuid"
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

func (cfg *apiConfig) handlerCreateChirp(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	chirp := struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&chirp)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, "something went wrong", err)
		return
	}

	if len(chirp.Body) > 140 {
		respondError(writer, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	validBody := validateProfanity(chirp.Body)

	params := database.CreateChirpParams{
		Body:   validBody,
		UserID: chirp.UserID,
	}

	createdChirp, err := cfg.queries.CreateChirp(context.Background(), params)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	respond(writer, http.StatusCreated, createdChirp)
}

func (cfg *apiConfig) handlerCreateUser(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	email := struct {
		Email string `json:"email"`
	}{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&email)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	_, err = mail.ParseAddress(email.Email)
	if err != nil {
		respondError(writer, http.StatusBadRequest, "Invalid email address", err)
		return
	}

	dbUser, err := cfg.queries.CreateUser(context.Background(), email.Email)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	respond(writer, http.StatusCreated, dbUser)
}
