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
	Id           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

func mapUser(dbUser database.User) User {
	return User{
		Id:          dbUser.ID,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
		Email:       dbUser.Email,
		IsChirpyRed: dbUser.IsChirpyRed,
	}
}

type login struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
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

	refreshToken, token := createTokens(dbUser.ID, writer, cfg.jwtSecret)
	if refreshToken == "" || token == "" {
		return
	}

	pars := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	}
	_, err = cfg.queries.CreateRefreshToken(context.Background(), pars)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	user := mapUser(dbUser)
	user.Token = token
	user.RefreshToken = refreshToken
	respond(writer, http.StatusOK, user)
}

func createTokens(userId uuid.UUID, writer http.ResponseWriter, secret string) (refreshToken, token string) {
	token, err := auth.MakeJWT(userId, secret, time.Hour)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	refreshToken, err = auth.MakeRefreshToken()
	if err != nil {
		respondError(writer, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	return refreshToken, token
}

func (cfg *apiConfig) handlerUpdateUser(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondError(writer, http.StatusUnauthorized, "Could not read the token", err)
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondError(writer, http.StatusUnauthorized, "Error validating token", err)
		return
	}

	var login login
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&login)
	if err != nil {
		respondError(writer, http.StatusBadRequest, "Could not read request body", err)
		return
	}

	password, err := auth.HashPassword(login.Password)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	pars := database.UpdateUserCredentialsParams{
		Email:          login.Email,
		HashedPassword: password,
		ID:             userId,
	}
	user, err := cfg.queries.UpdateUserCredentials(context.Background(), pars)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, "Error updating credentials", err)
		return
	}

	respond(writer, http.StatusOK, mapUser(user))
}
