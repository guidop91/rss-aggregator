package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/guidop91/rss-aggregator/internal/database"
)

type createFollowFeedParameters struct {
	FeedId uuid.UUID `json:"feed_id"`
}

func (apiCfg *apiConfig) handleCreateFollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	params := &createFollowFeedParameters{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't decode request body: %v", err))
		return
	}

	feedFollow, dbErr := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		FeedID:    params.FeedId,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
	})

	if dbErr != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't write to database: %v", dbErr))
		return
	}

	respondWithJSON(w, 200, databaseFeedFollowToFeedFollow(feedFollow))
}

func (apiCfg *apiConfig) handleGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, dbErr := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if dbErr != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get feed follows: %v", dbErr))
		return
	}

	parsedFeedFollows := []FeedFollow{}
	for _, feedFollow := range feedFollows {
		parsedFeedFollows = append(parsedFeedFollows, databaseFeedFollowToFeedFollow(feedFollow))
	}

	respondWithJSON(w, 200, parsedFeedFollows)
}

func (apiCfg *apiConfig) handleDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	// get id from query params
	feedFollowId := chi.URLParam(r, "feed_id")
	if feedFollowId == "" {
		respondWithError(w, 400, "Required id query param was not provided")
		return
	}

	// get user feed follows
	feedFollows, dbErr := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if dbErr != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get feed follows: %v", dbErr))
		return
	}

	var emptyFeedFollow FeedFollow
	feedFollowToDelete := emptyFeedFollow
	for _, feedFollow := range feedFollows {
		if feedFollow.ID.String() == feedFollowId {
			feedFollowToDelete = databaseFeedFollowToFeedFollow(feedFollow)
			break
		}
	}

	if feedFollowToDelete == emptyFeedFollow {
		respondWithError(w, 404, "Feed follow to delete was not found")
		return
	}

	deletedFeedFollow, deleteErr := apiCfg.DB.DeleteFeedFollow(r.Context(), feedFollowToDelete.ID)
	if deleteErr != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't delete feed follow entry with id: %v", feedFollowToDelete.ID))
		return
	}

	respondWithJSON(w, 200, deletedFeedFollow)
}
