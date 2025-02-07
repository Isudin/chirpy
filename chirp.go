package main

import (
	"time"

	"github.com/Isudin/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	statusCode int
	Id         uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Body       string    `json:"body"`
	UserID     uuid.UUID `json:"user_id"`
}

func (c *Chirp) getStatusCode() int {
	return c.statusCode
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
