package main

import (
	"log"
	"path/filepath"

	"github.com/go-pg/migrations/v8"

	"github.com/Torakushi/example-ws-events/internal/persistence"
)

func main() {
	// Database connection
	db, err := persistence.NewConnection("postgres://postgres:postgres@localhost:5432/nested?sslmode=disable")
	defer db.Close()

	collection := migrations.NewCollection()

	_, _, err = collection.Run(db, "init")
	if err != nil {
		log.Fatal(err)
		return
	}

	err = collection.DiscoverSQLMigrations(filepath.Join(".", "migrations"))
	if err != nil {
		log.Fatal(err)
		return
	}

	// Run the migrations
	_, _, err = collection.Run(db, "up")
	if err != nil {
		log.Fatal(err)
		return
	}
}
