package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/guidop91/rss-aggregator/internal/database"
)

func dbConnect() *database.Queries {
	postgresConnString := os.Getenv("DB_URL")
	variableMissing(postgresConnString)

	db, errDb := sql.Open("postgres", postgresConnString)
	if errDb != nil {
		log.Fatal("Failure to connect to DB", errDb)
	}

	return database.New(db)
}
