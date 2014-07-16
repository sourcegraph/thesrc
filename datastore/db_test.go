package datastore

import (
	"log"
	"os"
	"strings"
)

func init() {
	// Make sure we don't run the tests, which clobber data, on the main DB.
	dbname := os.Getenv("PGDATABASE")
	if dbname == "" {
		dbname = "thesrctest"
	}
	if !strings.HasSuffix(dbname, "test") {
		dbname += "test"
	}
	if err := os.Setenv("PGDATABASE", dbname); err != nil {
		log.Fatal(err)
	}

	// Reset DB.
	Connect()
	Drop()
	Create()
}
