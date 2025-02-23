package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Isudin/chirpy/internal/auth"
	"github.com/Isudin/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id,omitempty"`
}

func mapChirp(dbChirp database.Chirp) Chirp {
	return Chirp{
		Id:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
}

func mapChirps(dbChirps []database.Chirp) []Chirp {
	var chirps []Chirp
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, mapChirp(dbChirp))
	}

	return chirps
}

func (cfg *apiConfig) handlerCreateChirp(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	var chirp Chirp
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&chirp)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondError(writer, http.StatusUnauthorized, err.Error(), err)
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondError(writer, http.StatusUnauthorized, "Error validating token", err)
		return
	}

	if len(chirp.Body) > 140 {
		respondError(writer, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	validBody := validateProfanity(chirp.Body)

	params := database.CreateChirpParams{
		Body:   validBody,
		UserID: userId,
	}

	createdChirp, err := cfg.queries.CreateChirp(context.Background(), params)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	respond(writer, http.StatusCreated, mapChirp(createdChirp))
}

func (cfg *apiConfig) handlerListChirps(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	chirps, err := cfg.queries.ListChirps(context.Background())
	if err != nil {
		respondError(writer, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	respond(writer, http.StatusOK, mapChirps(chirps))
}

func (cfg *apiConfig) handlerGetChirpById(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	parsedId, err := uuid.Parse(req.PathValue("id"))
	if err != nil {
		respondError(writer, http.StatusBadRequest, "Invalid chirp id", err)
		return
	}

	chirp, err := cfg.queries.GetChirpById(context.Background(), parsedId)
	if err != nil {
		respondError(writer, http.StatusNotFound, "Chirp not found", err)
		return
	}

	respond(writer, http.StatusOK, mapChirp(chirp))
}

func (cfg *apiConfig) handlerDeleteChirp(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondError(writer, http.StatusUnauthorized, "Could not validate the token", err)
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondError(writer, http.StatusUnauthorized, "Could not validate the token", err)
		return
	}

	chirpId, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respondError(writer, http.StatusForbidden, "Cannot delete the chirp", err)
		return
	}

	chirp, err := cfg.queries.GetChirpById(context.Background(), chirpId)
	if err != nil {
		respondError(writer, http.StatusNotFound, "Chirp not found", err)
		return
	}

	if chirp.UserID != userId {
		respondError(writer, http.StatusForbidden, "Cannot delete the chirp", err)
		return
	}

	pars := database.DeleteChirpParams{
		ID:     chirpId,
		UserID: userId,
	}
	err = cfg.queries.DeleteChirp(context.Background(), pars)
	if err != nil {
		respondError(writer, http.StatusNotFound, "Chirp not found", err)
		return
	}

	respond(writer, http.StatusNoContent, err)
}
