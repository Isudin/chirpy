package main

import (
	"net/http"
)

type RefreshToken struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefresh(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	// token, err := auth.GetBearerToken(req.Header)
	// if err != nil {
	// 	respondError(writer, http.StatusUnauthorized, "Invalid credentialns", err)
	// 	return
	// }

	// dbToken, err := cfg.queries.GetTokenData(context.Background(), token)
	// if err != nil {
	// 	respondError(writer, http.StatusInternalServerError, "Something went wrong", err)
	// 	return
	// }

	// respond(writer, http.StatusOK, RefreshToken{Token: token})
}
