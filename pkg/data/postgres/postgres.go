// Package postgres contains the official implementation of the PostgreSQL integration.
package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/dchenk/mazewire/pkg/env"
	"github.com/lib/pq"
)

const (
	errDupObject = "42710" // https://www.postgresql.org/docs/10/static/errcodes-appendix.html
)

// DB implements the data.DB interface for PostgreSQL databases.
type DB struct {
	db *sql.DB
}

func (d *DB) Init(envVars map[string]string) error {
	dbConn := envVars[env.VarDbConnection]
	dbName := envVars[env.VarDbName]
	dbParams := strings.TrimSpace(envVars[env.VarDbParams])
	if dbParams == "" {
		// This is the only non-required variable that PostgreSQL uses.
		return fmt.Errorf("postgres: must set non-empty DB_PARAMS setting")
	}
	if dbParams[0] != '?' {
		dbParams = "?" + dbParams
	}
	if dbConn[len(dbConn)-1] != '/' {
		dbConn += "/"
	}
	var err error
	d.db, err = sql.Open("postgres", dbConn+dbName+dbParams)
	return err
}

func (d *DB) Ping() error {
	return d.db.Ping()
}

// ErrIsDupKey says if the error reports a database error indicating that an insert would
// cause a duplicate key error.
func (*DB) ErrIsDupKey(e error) bool {
	if e == nil {
		return false
	}
	me, ok := e.(*pq.Error)
	if !ok {
		return false
	}
	return me.Code == errDupObject
}
