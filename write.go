package main

import (
	"net/http"
)

func WriteByteResponse(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Write([]byte(message))
}
func WriteJsonResponse(w http.ResponseWriter, code int, message []byte) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(message)
}
