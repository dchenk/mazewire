package cockroach

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dchenk/mazewire/pkg/data"
	"github.com/dchenk/mazewire/pkg/data/util"
	"github.com/dchenk/mazewire/pkg/log"
	"github.com/pkg/errors"
)

// blobsWhere retrieves blobs from the specified site's table.
func (d *DB) blobsWhere(siteID int64, cond string, args ...interface{}) ([]data.Blob, error) {
	rows, err := d.selStar(data.BlobsTable(siteID), cond, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	bbs := make([]data.Blob, 0, 4)
	for rows.Next() {
		var b data.Blob
		if err = rows.Scan(&b.Id, &b.Role, &b.K, &b.V, &b.Updated); err != nil {
			return bbs, err
		}
		bbs = append(bbs, b)
	}
	return bbs, rows.Err()
}

// BlobsByRoleInK selects blobs where the role matches a role in the given list and the K equals the given k.
// If not matching Blobs are found, an empty slice is returned but no error.
func (d *DB) BlobsByRoleInK(site int64, roles []string, k int64) ([]data.Blob, error) {
	rows, err := d.selStar(data.BlobsTable(site), "role IN ("+util.JoinQuoted(roles)+") AND k=$1", k)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	bbs := make([]data.Blob, 0, 2)
	for rows.Next() {
		var b data.Blob
		if err = rows.Scan(&b.Id, &b.Role, &b.K, &b.V, &b.Updated); err != nil {
			return bbs, err
		}
		bbs = append(bbs, b)
	}
	return bbs, rows.Err()
}

// BlobByRoleLikeLast returns the last Blob that is matched to the role column using the LIKE feature, ordered
// by the record ID.
func (d *DB) BlobByRoleLikeLast(site int64, kPattern string) (*data.Blob, error) {
	return firstBlob(d.blobsWhere(site, "role LIKE "+db.SingleQuote(kPattern)+" ORDER BY id DESC LIMIT 1"))
}

// BlobsForTree retrieves the data blobs necessary to build a page tree. The IDs list contains the IDs of the tree elements, the
// static settings of dynamic elements, and the published rows of the dynamic elements.
func (d *DB) BlobsIdIn(site int64, IDs []int64) ([]Blob, error) {
	return blobsWhere(s, "id IN ("+util.IntsList(IDs)+")")
}

// BlobsIdInMapped gets the blobs selected with the map keys being each of the row's ID.
func (s *Site) BlobsIdInMapped(site int64, IDs []int64) (map[int64]Blob, error) {
	bbs, err := s.BlobsIdIn(IDs)
	if err != nil {
		return nil, err
	}
	m := make(map[int64]Blob, len(bbs))
	for i := range bbs {
		m[bbs[i].Id] = bbs[i]
	}
	return m, nil
}

// BlobByID returns a single Blob by its ID.
// This function wil return a sql.ErrNoRows error if no matching row is found.
func (s *Site) BlobByID(site int64, id int64) (*Blob, error) {
	return firstBlob(blobsWhere(s, db.IdEq(id)))
}

// BlobsByRoleK returns multiple Blob records selected by both the role and the k.
// This function never returns a sql.ErrNoRows error.
func (s *Site) BlobsByRoleK(site int64, role string, k int64) ([]Blob, error) {
	return blobsWhere(s, "role=$1 AND k=$2", role, k)
}

// BlobByRoleK returns a single Blob selected by both the role and the k. The rows are retrieved
// This function wil return a sql.ErrNoRows error if no matching row is found.
//func (s *Site) BlobByRoleK(role string, k int64) (*Blob, error) {
//	return firstBlob(s.BlobsByRoleK(role, k))
//}

// BlobByRoleKLast returns the last Blob, ordered by the row ID, selected by both the role and the k.
// This function wil return a sql.ErrNoRows error if no matching row is found.
func BlobByRoleKLast(site int64, role string, k int64) (*Blob, error) {
	return firstBlob(blobsWhere(site, "role=$1 AND k=$2 ORDER BY id LIMIT 1"))
}

// BlobVByRoleK returns the v of a single Blob selected by both the role and the k.
// This function wil return a sql.ErrNoRows error if no matching row is found.
func BlobVByRoleK(site int64, role string, k int64) (v []byte, err error) {

	rows, err := selCols(db.BlobsTable(s.Id), "v", "role=$1 AND k=$2 LIMIT 1", role, k)
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

// BlobInsert inserts a Blob with the given data and returns the ID of the inserted row.
func (s *Site) BlobInsert(site int64, role string, k int64, v []byte) (insertID int64, err error) {
	err = db.QueryRow("INSERT INTO "+db.BlobsTable(s.Id)+" (role,k,v) VALUES ($1,$2,$3) RETURNING id", role, k, v).Scan(&insertID)
	return
}

// BlobInsertFull inserts a Blob with the given data, including a timestamp for the "updated" column, and returns the ID of
// the inserted row.
func (s *Site) BlobInsertFull(site int64, role string, k int64, v []byte, updated time.Time) (insertID int64, err error) {
	err = db.QueryRow("INSERT INTO "+db.BlobsTable(s.Id)+" (role,k,v,updated) VALUES ($1,$2,$3,$4) RETURNING id",
		role, k, v, updated).Scan(&insertID)
	return
}

// BlobUpdate updates a Blob with the given data and returns the number of rows affected.
func (s *Site) BlobUpdate(site int64, id int64, v []byte) (int64, error) {
	res, err := db.Exec("UPDATE "+db.BlobsTable(s.Id)+" SET v=$1 WHERE id=$2", v, id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// BlobInsertMulti inserts multiple Blob elements into the database.
//func (s *Site) BlobInsertMulti(blobs []Blob) error {
//	_, err := db.Exec("INSERT INTO "+db.BlobsTable(s.Id)+" (role,k,v) VALUES ($1,$2,$3)", blobs...) // TODO: this is wrong, should be unrolled
//	return err
//}

// BlobAppendRole appends a role to the role column of a Blob and returns the number of rows affected.
func (d *DB) BlobAppendRole(site int64, id int64, role string) (int64, error) {
	res, err := db.Exec("UPDATE " + db.BlobsTable(s.Id) + " SET role=CONCAT(role," + db.SingleQuote(role) + ") WHERE id=" + strconv.FormatInt(id, 10))
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// BlobDelete deletes a site blob row selected by its ID and returns the number of rows affected.
func (d *DB) BlobDelete(site int64, id int64) (rowsAffected int64, err error) {
	res, err := d.db.Exec("DELETE FROM " + db.BlobsTable(s.Id) + " WHERE id=" + strconv.FormatInt(id, 10))
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// firstBlob returns a pointer to the first Blob from the given slice, or an error if err != nil or if the
// slice is empty.
func firstBlob(bs []data.Blob, err error) (*data.Blob, error) {
	if err != nil {
		return nil, err
	}
	if len(bs) == 0 {
		return nil, sql.ErrNoRows
	}
	return &bs[0], nil
}

// blobsTableSchema returns the schema for a Blobs table for a site. If either tableName or sequenceName is blank,
// the function panics.
func (d *DB) blobsTableSchema(tableName, sequenceName string) string {
	if tableName == "" || sequenceName == "" {
		msg := fmt.Sprintf("cockroach: blank table name (%q) or sequence name (%q)", tableName, sequenceName)
		log.Critical(nil, msg, errors.New("invalid schema name"))
		panic(msg)
	}
	var s strings.Builder
	s.WriteString("CREATE TABLE ")
	s.WriteString(tableName)
	s.WriteString(" (id INT PRIMARY KEY DEFAULT nextval('")
	s.WriteString(sequenceName)
	s.WriteString("'), role STRING NOT NULL, k INT NOT NULL DEFAULT 0, v BYTES NOT NULL, updated TIMESTAMP NOT NULL DEFAULT now(), INDEX role_k (role,k))")
	return s.String()
}
