package main

import (
	"context"
	"net/http"
	"time"

	"github.com/Isudin/chirpy/internal/auth"
)

type Token struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefresh(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondError(writer, http.StatusUnauthorized, "Invalid credentials", err)
		return
	}

	dbToken, err := cfg.queries.GetTokenData(context.Background(), refreshToken)
	if err != nil {
		respondError(writer, http.StatusUnauthorized, "Invalid credentials", err)
		return
	}

	token, err := auth.MakeJWT(dbToken.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	respond(writer, http.StatusOK, Token{Token: token})
}

func (cfg *apiConfig) handlerRevoke(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondError(writer, http.StatusUnauthorized, "Invalid credentials", err)
		return
	}

	err = cfg.queries.RevokeToken(context.Background(), refreshToken)
	if err != nil {
		respondError(writer, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	respond(writer, http.StatusNoContent, nil)
}
