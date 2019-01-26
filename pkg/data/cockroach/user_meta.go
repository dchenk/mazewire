package cockroach

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/dchenk/mazewire/pkg/data"
	"github.com/dchenk/mazewire/pkg/data/util"
)

// userMetaWhere retrieves UserMeta records and never returns a sql.ErrNoRows error.
func (d *DB) userMetaWhere(cond string, args ...interface{}) ([]data.UserMeta, error) {
	rows, err := d.selStar(data.UserMetaTable, cond, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ums := make([]data.UserMeta, 0, 2)
	for rows.Next() {
		var um data.UserMeta
		if err = rows.Scan(&um.UserId, &um.K, &um.V, &um.Updated); err != nil {
			return ums, err
		}
		ums = append(ums, um)
	}
	return ums, rows.Err()
}

// UserMetaById returns a all the *UserMeta selected by user ID.
// This function never returns a sql.ErrNoRows error.
func (d *DB) UserMetaById(userID int64) ([]data.UserMeta, error) {
	return d.userMetaWhere(util.UserEq(userID))
}

// UserMetaByIdMapped returns a map of the K => V pairs of the user meta selected by user ID.
// This function always returns a non-nil map, though the map may be empty if an error occurs.
func (d *DB) UserMetaByIdMapped(userID int64) (map[string][]byte, error) {
	mapped := make(map[string][]byte)
	ums, err := d.UserMetaById(userID)
	if err != nil {
		return mapped, err
	}
	for i := range ums {
		mapped[ums[i].K] = ums[i].V
	}
	return mapped, nil
}

// UserMetaByIdKey returns a single *UserMeta selected by the user's ID and the K.
func (d *DB) UserMetaByIdKey(userID int64, k string) (*data.UserMeta, error) {
	return firstUserMeta(d.userMetaWhere("user_id=$1 AND k=$2", userID, k))
}

// UserMetaByIdLikeKey returns the *UserMeta rows selected by the user's ID and K values matched by LIKE.
func (d *DB) UserMetaByIdLikeKey(userID int64, k string) ([]data.UserMeta, error) {
	return d.userMetaWhere("user_id=$1 AND k LIKE $2", userID, k)
}

// UserMetaV selects just like UserMetaByIdKey but retrieves only the V of the meta datum.
func (d *DB) UserMetaV(userID int64, k string) (v []byte, err error) {
	rows, err := d.selCols(data.UserMetaTable, "v", "user_id=$1 AND k=$2", userID, k)
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&v)
		return
	}
	err = sql.ErrNoRows
	return
}

// UserMetaUpdate updates a UserMeta (creating a new record in the database if necessary) and returns the number
// of rows affected. The primary key is set by both the user ID and k. The Updated time cannot be set directly but
// is automatically updated by the database.
func (d *DB) UserMetaUpdate(userID int64, k string, v []byte) (int64, error) {
	res, err := d.db.Exec("UPSERT INTO "+data.UserMetaTable+" (user_id,k,v) VALUES ($1,$2,$3)", userID, k, v)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// UserMetasUpdate updates UserMeta records, creating new records in the database if necessary and returns the number
// of rows affected.
//
// The primary key is set by both the user ID and the value of K in each element. The Updated time cannot be set
// directly but is automatically updated by the database.
func (d *DB) UserMetasUpdate(ums []data.UserMeta) (int64, error) {
	if len(ums) == 0 {
		return 0, nil
	}
	var q strings.Builder
	q.WriteString("UPSERT INTO " + data.UserMetaTable + " (user_id,k,v) VALUES ")
	for i := range ums {
		if i > 0 {
			q.WriteByte(',')
		}
		um := &ums[i]
		q.WriteByte('(')
		q.WriteString(strconv.FormatInt(um.UserId, 10))
		q.WriteString(",$")
		q.WriteString(strconv.Itoa(2*i + 1))
		q.WriteString(",$")
		q.WriteString(strconv.Itoa(2*i + 2))
		q.WriteByte(')')
	}
	res, err := d.db.Exec(q.String(), userMetaArgs(ums)...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (d *DB) UserMetaDelete(userID int64, k string) (int64, error) {
	res, err := d.db.Exec("DELETE FROM "+data.UserMetaTable+" WHERE user_id=$1 AND k=$2", userID, k)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// userMetaArgs gives a slice of all the arguments from the K and V fields to be set by ums flattened out.
func userMetaArgs(ums []data.UserMeta) []interface{} {
	args := make([]interface{}, 2*len(ums))
	for i := range ums {
		args[2*i] = ums[i].K
		args[2*i+1] = ums[i].V
	}
	return args
}

// firstUserMeta returns the first *UserMeta from the given slice.
func firstUserMeta(ums []data.UserMeta, err error) (*data.UserMeta, error) {
	if err != nil {
		return nil, err
	}
	if len(ums) == 0 {
		return nil, sql.ErrNoRows
	}
	return &ums[0], nil
}
