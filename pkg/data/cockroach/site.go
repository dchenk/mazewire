package cockroach

import (
	"database/sql"

	"github.com/dchenk/mazewire/pkg/data"
)

func sitesWhere(cond string, args ...interface{}) ([]data.Site, error) {
	rows, err := selCols(data.SitesTable, siteCols, cond, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	sites := make([]data.Site, 0, 4)
	for rows.Next() {
		var s data.Site
		if err = rows.Scan(&s.Id, &s.Domain, &s.Name, &s.Logo, &s.Favicon, &s.Tls); err != nil {
			return sites, err
		}
		sites = append(sites, s)
	}
	return sites, rows.Err()
}

// SitesByIDs retrieves sites by their ID.
func SitesByIDs(ids []int64) ([]data.Site, error) {
	return sitesWhere("id IN (" + data.JoinIDs(ids) + ")")
}

// SiteByDomain retrieves a site by its domain.
func SiteByDomain(domain string) (*data.Site, error) {
	rows, err := db.Query("SELECT " + siteCols + " FROM " + data.SitesTable + " WHERE domain=" + db.SingleQuote(domain))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	s := new(data.Site)
	if rows.Next() {
		err = rows.Scan(&s.Id, &s.Domain, &s.Name, &s.Logo, &s.Favicon, &s.Tls)
	} else {
		return nil, sql.ErrNoRows
	}
	return s, err
}

// InsertSite inserts a record into the sites table and create all needed tables for the site. The int64 returned
// is the ID of the new site (the inserted row). The domain passed in must have already been validated as a real
// domain name.
func (d *DB) InsertSite(domain, name string) (int64, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return 0, err
	}

	var siteID int64
	err = tx.QueryRow("INSERT INTO "+data.SitesTable+" (domain,name) VALUES ($1,$2) RETURNING id", domain, name).Scan(&siteID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tableName := data.BlobsTable(siteID)
	sequenceName := tableName + "_id"

	_, err = tx.Exec("CREATE SEQUENCE " + sequenceName)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	_, err = tx.Exec(d.blobsTableSchema(tableName, sequenceName))
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	return siteID, err
}
