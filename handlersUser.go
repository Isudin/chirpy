package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/mail"
	"time"

	"github.com/Isudin/chirpy/internal/auth"
	"github.com/Isudin/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token,omitempty"`
}

func mapUser(dbUser database.User) User {
	return User{
		Id:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
}

type login struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	ExpiresIn int    `json:"expires_in_seconds"`
}

func (cfg *apiConfig) handlerCreateUser(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	var login login

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&login)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	_, err = mail.ParseAddress(login.Email)
	if err != nil {
		respondError(writer, http.StatusBadRequest, "Invalid email address", err)
		return
	}

	hash, err := auth.HashPassword(login.Password)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	params := database.CreateUserParams{
		Email:          login.Email,
		HashedPassword: hash,
	}

	dbUser, err := cfg.queries.CreateUser(context.Background(), params)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	respond(writer, http.StatusCreated, mapUser(dbUser))
}

func (cfg *apiConfig) handlerLogin(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	var login login
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&login)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	dbUser, err := cfg.queries.GetUserByEmail(context.Background(), login.Email)
	if err != nil {
		if dbUser.ID == uuid.Nil {
			respondError(writer, http.StatusNotFound, "User not found", err)
			return
		}

		respondError(writer, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	err = auth.CheckPasswordHash(login.Password, dbUser.HashedPassword)
	if err != nil {
		respondError(writer, http.StatusUnauthorized, "Wrong passoword", err)
		return
	}

	expiresIn := time.Hour
	if login.ExpiresIn > 0 {
		expiresIn = time.Duration(login.ExpiresIn * int(time.Second))
	}
	token, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, expiresIn)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	user := mapUser(dbUser)
	user.Token = token
	respond(writer, http.StatusOK, user)
}
