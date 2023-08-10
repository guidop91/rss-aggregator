package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
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
