// Package mysql contains the official implementation of the MySQL integration.
package mysql

import (
	"database/sql"

	"github.com/dchenk/mazewire/pkg/env"
	"github.com/go-sql-driver/mysql"
)

const (
	errDupEntry = 1062 // https://dev.mysql.com/doc/refman/8.0/en/error-messages-server.html#error_er_dup_entry
)

// DB implements the data.DB interface for MySQL databases.
type DB struct {
	db *sql.DB
}

func (d *DB) Init(envVars map[string]string) error {
	dbConn := envVars[env.VarDbConnection]
	dbName := envVars[env.VarDbName]
	/*
		dbParams := strings.TrimSpace(envVars[env.VarDbParams])
		if dbParams == "" {
			// This is the only non-required variable that MySQL uses.
			return fmt.Errorf("mysql: must set non-empty DB_PARAMS setting")
		}
		if dbParams[0] != '?' {
			dbParams = "?" + dbParams
		}
	*/
	if dbConn[len(dbConn)-1] != '/' {
		dbConn += "/"
	}
	var err error
	d.db, err = sql.Open("mysql", dbConn+dbName /*+dbParams*/)
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
	me, ok := e.(*mysql.MySQLError)
	if !ok {
		return false
	}
	return me.Number == errDupEntry
}
