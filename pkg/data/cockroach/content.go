package cockroach

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/dchenk/mazewire/pkg/data"
	"github.com/dchenk/mazewire/pkg/data/util"
)

// contentWhere retrieves all of the columns in the Content rows specified by cond (or all the rows if cond is blank).
// The cond string can contain more than just a WHERE clause, but also LIMIT or ORDER BY.
// If there is an error making the sending the query, the slice returned is nil. Otherwise the slice contains any of the
// rows already scanned.
func (d *DB) contentsWhere(cond string, args ...interface{}) ([]data.Content, error) {
	rows, err := d.selStar(data.ContentTable, cond, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cns := make([]data.Content, 0, 4)
	for rows.Next() {
		var c data.Content
		if err = rows.Scan(&c.Id, &c.Site, &c.Slug, &c.Author, &c.Type, &c.Parent, &c.Title, &c.MetaTitle, &c.MetaDesc, &c.Body, &c.Status, &c.Updated); err != nil {
			return cns, err
		}
		cns = append(cns, c)
	}
	return cns, rows.Err()
}

// ContentByID returns a single Content by its ID.
func (d *DB) ContentByID(id int64) (*data.Content, error) {
	return firstContent(d.contentsWhere(util.IdEq(id)))
}

// ContentsByIDs returns the Content rows selected by their ID.
func (d *DB) ContentsByIDs(ids []int64) ([]data.Content, error) {
	return d.contentsWhere("id IN (" + util.IntsList(ids) + ")")
}

// ContentBySiteSlug returns a single Content by the site ID and slug given; the string returned is the parent slug, which is
// blank if the content has no parent.
func (d *DB) ContentBySiteSlug(siteID int64, slug string) (*data.Content, string, error) {
	c := data.Content{Slug: slug}
	var parentSlug string
	r := d.db.QueryRow(
		"SELECT a.id,a.author,a.type,a.parent,a.title,a.meta_title,a.meta_desc,a.body,a.status,a.updated,IFNULL(b.slug,'') FROM " + data.ContentTable +
			" AS a LEFT JOIN " + data.ContentTable + " AS b ON a.parent=b.id WHERE a.site=" + strconv.FormatInt(siteID, 10) + " AND a.slug=" + util.SingleQuote(slug))
	err := r.Scan(&c.Id, &c.Author, &c.Type, &c.Parent, &c.Title, &c.MetaTitle, &c.MetaDesc, &c.Body, &c.Status, &c.Updated, &parentSlug)
	return &c, parentSlug, err
}

func (d *DB) ContentByAuthor(authorID int64) ([]data.Content, error) {
	return nil, nil // TODO
}

// ContentsList retrieves the specified Content items. The parent argument should be either nil to indicate that the
// parent column should not be considered in the query, or a list of the parent IDs to which the returned items should belong.
// The value of offset may never be negative, which is why its type is uint64.
func (d *DB) ContentsList(siteID int64, cType string, parents []int64, statuses []string, authorID int64, offset uint64) ([]data.Content, error) {
	var parentCheck string
	if len(parents) > 0 {
		parentCheck = " AND parent IN (" + util.IntsList(parents) + ")"
	}
	rows, err := d.db.Query("SELECT id,slug,title,author,parent,status FROM "+data.ContentTable+" WHERE site=$1 AND type=$2"+
		parentCheck+authorCheck(authorID)+statusList(statuses)+
		" ORDER BY title LIMIT 20"+getOffsetOrNot(offset), siteID, cType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cs := make([]data.Content, 0, 4)
	for rows.Next() {
		var c data.Content
		if err = rows.Scan(&c.Id, &c.Slug, &c.Title, &c.Author, &c.Parent, &c.Status); err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	return cs, rows.Close()
}

// CountContent says how many content pages/posts there are with the given parameters and how many such have no parent.
// Optionally give a non-zero authorID to also filter by author.
// The given statuses should not have any punctuation at all (should be already sanitized).
func (d *DB) CountContent(siteID int64, pType string, statuses []string, authorID int64) (countTotal uint16, countParentLevel uint32, err error) {
	coreQ := "SELECT COUNT(*) FROM " + data.ContentTable + " WHERE site=$1 AND type=$2" + authorCheck(authorID) + statusList(statuses)
	err = d.db.QueryRow("SELECT ("+coreQ+"),("+coreQ+" AND parent=0)", siteID, pType).Scan(&countTotal, &countParentLevel)
	return
}

// ContentCountSlug returns the number of rows with the given site ID and slug.
func (d *DB) ContentCountSlug(siteID int64, slug string) (count int64, err error) {
	err = d.db.QueryRow("SELECT COUNT(*) FROM "+data.ContentTable+" WHERE site=$1 AND slug=$2", siteID, slug).Scan(&count)
	return
}

func (d *DB) ContentInsert(siteID int64, slug string, author int64, pType string, parent int64, title string) (insertID int64, err error) {
	err = d.db.QueryRow("INSERT INTO "+data.ContentTable+" (site,slug,author,type,parent,title) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id",
		siteID, slug, author, pType, parent, title).Scan(&insertID)
	return
}

// ContentUpdate updates a single content record, the values passed in as column-name -> value pairs.
// The number returned is the number of rows affected by the update.
func (d *DB) ContentUpdate(contentID int64, vals map[string]interface{}) (int64, error) {
	if len(vals) == 0 {
		return 0, errors.New("data: no columns to update given")
	}
	assignments, args := setValuesList(vals)
	res, err := d.db.Exec("UPDATE "+data.ContentTable+" SET "+assignments+" WHERE id=$1", append(args, contentID)...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// firstContent returns the first *Content from the given slice, or an error if err != nil or if the slice is empty.
func firstContent(cs []data.Content, err error) (*data.Content, error) {
	if err != nil {
		return nil, err
	}
	if len(cs) == 0 {
		return nil, sql.ErrNoRows
	}
	return &cs[0], nil
}

// getOffsetOrNot returns the string " OFFSET n" where n is the offset if n > 0, or an empty string otherwise.
func getOffsetOrNot(n uint64) string {
	if n > 0 {
		return " OFFSET " + strconv.FormatUint(n, 10)
	}
	return ""
}

func statusList(s []string) string {
	if len(s) == 0 {
		return ""
	}
	var sl strings.Builder
	sl.WriteString(" AND status IN (")
	for i := range s {
		if i > 0 {
			sl.WriteByte(',')
		}
		sl.WriteString(util.SingleQuote(s[i]))
	}
	sl.WriteByte(')')
	return sl.String()
}

func authorCheck(id int64) string {
	if id > 0 {
		return " AND author=" + strconv.FormatInt(id, 10)
	}
	return ""
}
