package main

import (
	"fmt"
	"net/http"

	"github.com/guidop91/rss-aggregator/internal/auth"
	"github.com/guidop91/rss-aggregator/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, authErr := auth.GetApiKey(r.Header)
		if authErr != nil {
			respondWithError(w, 403, authErr.Error())
			return
		}

		user, dbErr := cfg.DB.GetUser(r.Context(), apiKey)
		if dbErr != nil {
			respondWithError(w, 403, fmt.Sprintf("Couldn't get user: %v", dbErr))
			return
		}

		handler(w, r, user)
	}
}
