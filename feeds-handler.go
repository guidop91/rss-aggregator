package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/guidop91/rss-aggregator/internal/database"
)

type createFeedParameters struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (apiCfg *apiConfig) handleCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	params := &createFeedParameters{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't decode request body: %v", err))
		return
	}

	feed, dbErr := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})
	if dbErr != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't write to Database: %v", dbErr))
		return
	}

	feedFollow, dbFeedFollowErr := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})
	if dbFeedFollowErr != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't write to Database: %v", dbFeedFollowErr))
		return
	}

	feedResponse := make(map[string]interface{})
	feedResponse["feed"] = databaseFeedToFeed(feed)
	feedResponse["feedFollow"] = databaseFeedFollowToFeedFollow(feedFollow)

	respondWithJSON(w, 200, feedResponse)
}

func (apiCfg *apiConfig) handleGetFeeds(w http.ResponseWriter, r *http.Request) {
	feedList, dbErr := apiCfg.DB.GetFeeds(r.Context())
	if dbErr != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't read from Database: %v", dbErr))
		return
	}

	parsedFeedList := []Feed{}
	for _, feedItem := range feedList {
		parsedFeedList = append(parsedFeedList, databaseFeedToFeed(feedItem))
	}

	respondWithJSON(w, 200, parsedFeedList)
}
