package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

// CreateDBHandler opens a handler to the database and returns an
// object that can be passed around to other db functions.
func CreateDBHandler() *sql.DB {
	db, err := sql.Open("postgres", "dbname=dbconcerts")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
