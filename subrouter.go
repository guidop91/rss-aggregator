package main

import "github.com/go-chi/chi"

func (apiCfg apiConfig) assignRouteHandlers(subRouter *chi.Mux) {
	subRouter.Get("/err", handleError)
	subRouter.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handleGetFeedFollows))
	subRouter.Get("/feeds", apiCfg.handleGetFeeds)
	subRouter.Get("/healthz", handleReadiness)
	subRouter.Get("/users", apiCfg.middlewareAuth(apiCfg.handleGetUser))
	subRouter.Get("/posts", apiCfg.middlewareAuth(apiCfg.getPostsForUser))

	subRouter.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handleCreateFollowFeed))
	subRouter.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handleCreateFeed))
	subRouter.Post("/users", apiCfg.handleCreateUser)

	subRouter.Delete("/feed_follows/{feed_id}", apiCfg.middlewareAuth(apiCfg.handleDeleteFeedFollow))
}
