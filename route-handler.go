package main

import (
	"net/http"
)

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, struct{}{})
}

func handleError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 400, "Server error")
}
