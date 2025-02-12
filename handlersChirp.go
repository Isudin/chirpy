package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Isudin/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
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
	// chirp := struct {
	// 	Body   string    `json:"body"`
	// 	UserID uuid.UUID `json:"user_id"`
	// }{}
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
		respondError(writer, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	respond(writer, http.StatusOK, mapChirp(chirp))
}
