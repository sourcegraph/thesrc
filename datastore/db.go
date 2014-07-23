package datastore

import (
	"log"
	"os"
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
		setDBCredentialsFromRDSEnv()

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

// transact calls fn in a DB transaction. If dbh is a transaction, then it just
// calls the function. Otherwise, it begins a transaction, rolling back on
// failure and committing on success.
func transact(dbh modl.SqlExecutor, fn func(dbh modl.SqlExecutor) error) error {
	var sharedTx bool
	tx, sharedTx := dbh.(*modl.Transaction)
	if !sharedTx {
		var err error
		tx, err = dbh.(*modl.DbMap).Begin()
		if err != nil {
			return err
		}
		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()
	}

	if err := fn(tx); err != nil {
		return err
	}

	if !sharedTx {
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

// setDBCredentialsFromRDSEnv copies RDS env vars (RDS_*) to PostgreSQL env vars
// (PG*) for use when deploying to AWS.
func setDBCredentialsFromRDSEnv() {
	m := map[string]string{
		"PGUSER":     "RDS_USERNAME",
		"PGPASSWORD": "RDS_PASSWORD",
		"PGDATABASE": "RDS_DB_NAME",
		"PGHOST":     "RDS_HOSTNAME",
		"PGPORT":     "RDS_PORT",
	}
	for pgName, rdsName := range m {
		if err := os.Setenv(pgName, os.Getenv(rdsName)); err != nil {
			log.Fatal(err)
		}
	}
}
