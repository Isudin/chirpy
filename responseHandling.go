package main

import (
	"encoding/json"
	"log"
	"net/http"
)

//TODO: Refactor whole file, remove classes and interfaces and just use two methods instead

type Response interface {
	getStatusCode() int
}

// TODO: add omitempty to json tag
type ResponseError struct {
	statusCode int
	Error      string `json:"error"`
}

func (resp *ResponseError) getStatusCode() int {
	return resp.statusCode
}

// TODO: add omitempty to json tag
type ResponseValid struct {
	statusCode  int
	CleanedBody string `json:"cleaned_body"`
}

func (resp *ResponseValid) getStatusCode() int {
	return resp.statusCode
}

// Obsolete
func marshalResponse(writer http.ResponseWriter, resp Response) {
	writer.WriteHeader(resp.getStatusCode())
	resBody, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		return
	}
	writer.Write(resBody)
}
