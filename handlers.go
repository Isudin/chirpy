package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

func (cfg *apiConfig) handlerResetHits(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
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
