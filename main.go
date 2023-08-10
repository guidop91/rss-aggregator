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
	godotenv.Load()

	// Load required env variables
	portString := os.Getenv("PORT")
	variableMissing(portString)

	apiCfg := apiConfig{
		DB: dbConnect(),
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	subRouter := chi.NewRouter()
	subRouter.Get("/healthz", handleReadiness)
	subRouter.Get("/err", handleError)
	subRouter.Post("/users", apiCfg.handleCreateUser)

	router.Mount("/v1", subRouter)

	log.Printf("Server running on port %v\n", portString)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("Fail to run server", err)
	}
}
