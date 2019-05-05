// Package data supports interacting with the main database and defines a few important data types.
package data

import (
	"fmt"

	"github.com/dchenk/mazewire/pkg/env"
)

// Conn represents the connection to the main database.
//
// This variable is initialized when the Init function in this package is called, and its Close
// method is called before the program exits.
var Conn DB

// Init sets up the database connection to be used throughout the rest of the program's lifetime.
// The function must be called only once.
func Init() error {
	err := Conn.Init(env.Vars())
	if err != nil {
		return fmt.Errorf("data: could not connect to database; %v", err)
	}

	if err = Conn.Ping(); err != nil {
		return fmt.Errorf("data: could not ping database; %v", err)
	}

	return nil
}

// A DB is an abstract handle on a database connection implemented for a particular DBMS.
type DB interface {
	Init(env map[string]string) error
	Ping() error
	BeginTx() (Transaction, error)

	Ops

	// Close closes the underlying representation of the database connection.
	Close() error

	// ErrIsDupKey says if the error reports a database error indicating that an insert would
	// cause a duplicate key error.
	ErrIsDupKey(e error) bool
}

// A Transaction is an initialized database transaction. Each transaction must end with a call to
// either Commit or Rollback.
//
// A Transaction guarantees a serializable isolation level.
type Transaction interface {
	Commit() error
	Rollback() error
	Ops
}

// Ops includes all of the operations that may be performed on a database.
type Ops interface {
	SiteManager
	BlobManager
	ContentManager
	UserManager
	UserMetaManager
	OptionManager
}
