package main

import (
	"log"
)

func main() {
	db, err := setupDatabase()
	if err != nil {
		log.Fatalf("could not set up database: %v", err)
	}
	defer db.Close()

	r := setupRouter(db)
	r.Run(":8080")
}
