package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"
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
		log.Printf("error deleting users: %v", err)
		writer.WriteHeader(500)
		writer.Write([]byte("Something went wrong"))
	}

	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("OK"))
}

func handlerValidateChirp(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	chirp := struct {
		Body string `json:"body"`
	}{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&chirp)
	if err != nil {
		log.Printf("Error decoding parameters: %v", err)
		response := ResponseError{Error: "Something went wrong", statusCode: 500}
		marshalResponse(writer, &response)
		return
	}

	if len(chirp.Body) > 140 {
		response := ResponseError{Error: "Chirp is too long", statusCode: 400}
		marshalResponse(writer, &response)
		return
	}

	validString := validateProfanity(chirp.Body)
	response := ResponseValid{CleanedBody: validString, statusCode: 200}
	marshalResponse(writer, &response)
}

func (cfg *apiConfig) handlerCreateUser(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	email := struct {
		Email string `json:"email"`
	}{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&email)
	if err != nil {
		log.Printf("Error decoding parameters: %v", err)
		response := ResponseError{Error: "Something went wrong", statusCode: 500}
		marshalResponse(writer, &response)
		return
	}

	_, err = mail.ParseAddress(email.Email)
	if err != nil {
		log.Printf("error parsing email address: %v", err)
		response := ResponseError{Error: "Invalid email address", statusCode: 400}
		marshalResponse(writer, &response)
		return
	}

	dbUser, err := cfg.queries.CreateUser(context.Background(), email.Email)
	if err != nil {
		log.Printf("error creating new user: %v", err)
		response := ResponseError{Error: "Something went wrong", statusCode: 500}
		marshalResponse(writer, &response)
		return
	}

	user := mapUser(dbUser)
	user.statusCode = 201
	marshalResponse(writer, &user)
}
