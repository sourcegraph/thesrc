package datastore

import (
	"log"
	"sync"

	"github.com/jmoiron/modl"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DB is the global database.
var DB = &modl.DbMap{Dialect: modl.PostgresDialect{}}

// DBH is a modl.SqlExecutor interface to DB, the global database. It is better
// to use DBH instead of DB because it prevents you from calling methods that
// could not later be wrapped in a transaction.
var DBH modl.SqlExecutor = DB

var connectOnce sync.Once

// Connect connects to the PostgreSQL database specified by the PG* environment
// variables. It calls log.Fatal if it encounters an error.
func Connect() {
	connectOnce.Do(func() {
		var err error
		DB.Dbx, err = sqlx.Open("postgres", "")
		if err != nil {
			log.Fatal("Error connecting to PostgreSQL database (using PG* environment variables): ", err)
		}
		DB.Db = DB.Dbx.DB
	})
}

var createSQL []string

// Create the database schema. It calls log.Fatal if it encounters an error.
func Create() {
	if err := DB.CreateTablesIfNotExists(); err != nil {
		log.Fatal("Error creating tables: ", err)
	}
	for _, query := range createSQL {
		if _, err := DB.Exec(query); err != nil {
			log.Fatalf("Error running query %q: %s", query, err)
		}
	}
}

// Drop the database schema.
func Drop() {
	// TODO(sqs): raise errors?
	DB.DropTables()
}
