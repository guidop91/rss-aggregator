package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/guidop91/rss-aggregator/internal/auth"
	"github.com/guidop91/rss-aggregator/internal/database"
)

type createParameters struct {
	Name string `json:"name"`
}

func (apiCfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	params := &createParameters{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't decode request body: %v", err))
		return
	}

	user, dbErr := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Name:      params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		ID:        uuid.New(),
	})
	if dbErr != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't write to Database: %v", dbErr))
		return
	}

	respondWithJSON(w, 200, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handleGetUser(w http.ResponseWriter, r *http.Request) {
	apiKey, authErr := auth.GetApiKey(r.Header)
	if authErr != nil {
		respondWithError(w, 403, authErr.Error())
		return
	}

	user, dbErr := apiCfg.DB.GetUser(r.Context(), apiKey)
	if dbErr != nil {
		respondWithError(w, 403, fmt.Sprintf("Couldn't get user: %v", dbErr))
		return
	}

	respondWithJSON(w, 200, databaseUserToUser(user))
}
