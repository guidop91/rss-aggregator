package main

import (
	"log"
	"net/http"
	"os"
	"time"

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
	dbInstance := dbConnect()
	apiCfg := apiConfig{
		DB: dbInstance,
	}

	go startScraping(dbInstance, 10, time.Minute)

	// Create router
	router := chi.NewRouter()

	// Assign CORS options middleware
	router.Use(getCorsOptions())

	// Create subrouter with route handlers
	subRouter := chi.NewRouter()
	apiCfg.assignRouteHandlers(subRouter)

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
