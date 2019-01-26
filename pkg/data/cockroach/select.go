package cockroach

import (
	"context"
	"database/sql"
)

// selStar returns a *sql.Rows and possibly an error.
func (d *DB) selStar(table string, cond string, args ...interface{}) (*sql.Rows, error) {
	if cond == "" {
		return d.db.QueryContext(context.Background(), "SELECT * FROM "+table)
	}
	return d.db.QueryContext(context.Background(), "SELECT * FROM "+table+" WHERE "+cond, args...)
}

// selCols returns a *sql.Rows and possibly an error. The "cols" argument contains a list of all the
// columns to retrieve.
func (d *DB) selCols(table string, cols string, cond string, args ...interface{}) (*sql.Rows, error) {
	q := "SELECT " + cols + " FROM " + table
	if cond == "" {
		return d.db.QueryContext(context.Background(), q)
	}
	return d.db.QueryContext(context.Background(), q+" WHERE "+cond, args...)
}
