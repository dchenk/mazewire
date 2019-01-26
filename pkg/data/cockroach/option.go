package cockroach

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/dchenk/mazewire/pkg/data"
	"github.com/dchenk/mazewire/pkg/data/util"
)

func (d *DB) optionsWhere(cond string, args ...interface{}) ([]data.Option, error) {
	rows, err := d.selStar(data.OptionsTable, cond, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	opts := make([]data.Option, 0, 4)
	for rows.Next() {
		var o data.Option
		if err = rows.Scan(&o.Site, &o.K, &o.V); err != nil {
			return opts, err
		}
		opts = append(opts, o)
	}
	return opts, rows.Err()
}

// OptionByKey returns a site's option selected by its K.
func (d *DB) OptionByKey(site int64, k string) (*data.Option, error) {
	return firstOption(d.optionsWhere("site=$1 AND k=$2", site, k))
}

// OptionsLikeKey returns Option records selected by their K being like k.
func (d *DB) OptionsLikeKey(site int64, k string) ([]data.Option, error) {
	return d.optionsWhere("site=$1 AND k LIKE $2", site, k)
}

// OptionsKeyIn returns the options selected by site ID and the K IN list.
// The strings in the Ks list must be valid UTF-8 strings.
func (d *DB) OptionsKeyIn(site int64, Ks []string) ([]data.Option, error) {
	return d.optionsWhere("site=$1 AND k IN ("+data.JoinQuoted(Ks)+")", site)
}

// OptionsKeyInMapped returns a map of the K-V pairs of a site's options selected by the Ks IN list.
// This function always returns a non-nil map, though the map may be empty.
// An error is not returned if there are no records retrieved.
func (d *DB) OptionsKeyInMapped(site int64, Ks []string) (map[string][]byte, error) {
	opts, err := d.OptionsKeyIn(site, Ks)
	if err != nil {
		return make(map[string][]byte), err
	}
	mapped := make(map[string][]byte, len(opts))
	for i := range opts {
		mapped[opts[i].K] = opts[i].V
	}
	return mapped, nil
}

// OptionsKeyInMappedStr returns a map of the K-V pairs of a site's options selected by the Ks In list,
// with the V values converted to strings.
// This function always returns a non-nil map, though the map may be empty.
// An error is not returned if there are no records retrieved.
func (d *DB) OptionsKeyInMappedStr(site int64, Ks []string) (map[string]string, error) {
	opts, err := d.OptionsKeyIn(site, Ks)
	if err != nil {
		return make(map[string]string), err
	}
	mapped := make(map[string]string, len(opts))
	for i := range opts {
		mapped[opts[i].K] = string(opts[i].V)
	}
	return mapped, nil
}

// OptionV returns the V of a site's option selected by key.
func (d *DB) OptionV(site int64, k string) ([]byte, error) {
	fo, err := firstOption(d.optionsWhere("site=$1 AND k=$1", site, k))
	if fo != nil {
		return fo.V, err
	}
	return nil, err
}

// OptionUpdate updates an Option or creates a new record in the database if necessary.
//
// The primary key is set by both the site ID and the value of K.
func (d *DB) OptionUpdate(site int64, k string, v []byte) (rowsAffected int64, err error) {
	res, err := d.db.Exec("UPSERT INTO "+data.OptionsTable+" (site,k,v) VALUES site=$1,k=$2,v=$3", site, k, v)
	if err != nil {
		return
	}
	rowsAffected, err = res.RowsAffected()
	return // rowsAffected may be 0 or something else (depending on if there is an error)
}

// OptionUpdateStr updates an Option or creates a new record in the database if necessary).
//
// The primary key is set by both the site ID and the value of K.
func (d *DB) OptionUpdateStr(site int64, k string, v string) (rowsAffected int64, err error) {
	return d.OptionUpdate(site, k, []byte(v))
}

// OptionsUpdate updates Option records or creates new records in the database if necessary.
// The strings passed in as keys in the map must be valid UTF-8 strings.
// Returned is the number of rows affected.
//
// The primary key is set by both the site ID and the value of K in each element.
func (d *DB) OptionsUpdate(site int64, opts map[string]string) (int64, error) {
	var q strings.Builder
	q.WriteString("UPSERT INTO " + data.OptionsTable + " (site,k,v) VALUES ")
	siteID := strconv.FormatInt(site, 10)
	didOne := false
	for k, v := range opts {
		if didOne {
			q.WriteByte(',')
		}
		didOne = true
		q.WriteByte('(')
		q.WriteString(siteID)
		q.WriteByte(',')
		q.WriteString(util.SingleQuote(k))
		q.WriteByte(',')
		q.WriteString(util.SingleQuote(v))
		q.WriteByte(')')
	}
	res, err := d.db.Exec(q.String())
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// OptionDelete deletes a site option selected by its K.
func (d *DB) OptionDelete(site int64, k string) (int64, error) {
	res, err := d.db.Exec("DELETE FROM "+data.OptionsTable+" WHERE site=$1 AND k=$2", site, k)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// optionArgs gives a slice of all the arguments to be set by opts flattened out.
//func optionArgs(siteID int64, opts map[string]string) []interface{} {
//	args := make([]interface{}, 0, 3*len(opts))
//	for k := range opts {
//		args = append(args, siteID, k, opts[k])
//	}
//	return args
//}

// firstOption returns the first *Option from the given slice.
// If the opts slice is empty, a sql.ErrNoRows error is returned.
func firstOption(opts []data.Option, err error) (*data.Option, error) {
	if err != nil {
		return nil, err
	}
	if len(opts) == 0 {
		return nil, sql.ErrNoRows
	}
	return &opts[0], nil
}
