package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respond(writer http.ResponseWriter, statusCode int, resp interface{}) {
	writer.WriteHeader(statusCode)
	resBody, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		return
	}
	writer.Write(resBody)
}

func respondError(writer http.ResponseWriter, statusCode int, msg string, err error) {
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Error message: %v\n", msg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	respond(writer, statusCode, errorResponse{Error: msg})
}
