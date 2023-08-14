package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/guidop91/rss-aggregator/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	// Load environment secrets from .env file
	godotenv.Load()

	// Load required env variables
	portString := os.Getenv("PORT")
	invariant(portString)

	// Create api config struct
	apiCfg := apiConfig{
		DB: dbConnect(),
	}

	// Create router
	router := chi.NewRouter()

	// Assign CORS options middleware
	router.Use(getCorsOptions())

	// Create subrouter with route handlers
	subRouter := chi.NewRouter()
	subRouter.Get("/healthz", handleReadiness)
	subRouter.Get("/err", handleError)
	subRouter.Get("/users", apiCfg.middlewareAuth(apiCfg.handleGetUser))
	subRouter.Post("/users", apiCfg.handleCreateUser)
	subRouter.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handleCreateFeed))
	subRouter.Get("/feeds", apiCfg.handleGetFeeds)
	subRouter.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handleCreateFollowFeed))
	subRouter.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handleGetFeedFollows))

	// Mount subrouter to main router
	router.Mount("/v1", subRouter)

	// Create HTTP server with router as handler
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server running on port %v\n", portString)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("Fail to run server", err)
	}
}

func getCorsOptions() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	})
}

// 7:57:14
