package main

import (
	"fmt"
	"net/http"

	"github.com/guidop91/rss-aggregator/internal/database"
)

func (apiCfg *apiConfig) getPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, dbErr := apiCfg.DB.GetUserPosts(r.Context(), database.GetUserPostsParams{
		UserID: user.ID,
		Limit:  10,
	})
	if dbErr != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't read from Database: %v", dbErr))
		return
	}

	parsedPostList := []Post{}
	for _, post := range posts {
		parsedPostList = append(parsedPostList, databasePostToPost(post))
	}

	respondWithJSON(w, 200, parsedPostList)
}
