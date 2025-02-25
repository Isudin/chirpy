package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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
		fmt.Println(err)
		respondError(writer, http.StatusBadRequest, "Could not decode body", err)
		return
	}

	if polka.Event != "user.upgraded" {
		fmt.Println("Event: " + polka.Event)
		respondError(writer, http.StatusNoContent, "", err)
		return
	}

	err = cfg.queries.UpgradeUser(context.Background(), polka.Data.UserId)
	if err != nil {
		fmt.Println(err)
		respondError(writer, http.StatusNotFound, "User not found", err)
		return
	}
	fmt.Println("No error")

	respond(writer, http.StatusNoContent, nil)
}
