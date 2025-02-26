package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Isudin/chirpy/internal/auth"
	"github.com/google/uuid"
)

type polkaUpgrade struct {
	Event string `json:"event"`
	Data  struct {
		UserId uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlerPolkaWebhook(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	var polka polkaUpgrade
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&polka)
	if err != nil {
		respondError(writer, http.StatusBadRequest, "Could not decode body", err)
		return
	}

	if polka.Event != "user.upgraded" {
		respondError(writer, http.StatusNoContent, "", err)
		return
	}

	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil || apiKey != cfg.polkaKey {
		respondError(writer, http.StatusUnauthorized, "Incorrect token", err)
		return
	}

	err = cfg.queries.UpgradeUser(context.Background(), polka.Data.UserId)
	if err != nil {
		respondError(writer, http.StatusNotFound, "User not found", err)
		return
	}

	respond(writer, http.StatusNoContent, nil)
}
