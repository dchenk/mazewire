package cockroach

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/dchenk/mazewire/pkg/data"
	"github.com/dchenk/mazewire/pkg/data/util"
)

// userCols is the columns retrieved when a single *User or a []User is queried. Other columns can be retrieved separately.
const userCols = "id,username,fname,lname,email"

// userWhere retrieves the basic columns in the User rows specified by cond (or all the rows if cond is blank).
// The cond string can contain more than just a WHERE clause, but also LIMIT or ORDER BY.
// If there is an error making the sending the query, the slice returned is nil. Otherwise the slice contains any of the
// rows already scanned.
// To get the user's password, use UserPass.
func (d *DB) usersWhere(cond string, args ...interface{}) ([]data.User, error) {
	rows, err := d.selCols(data.UsersTable, userCols, cond, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	us := make([]data.User, 0, 4)
	for rows.Next() {
		var u data.User
		if err = rows.Scan(&u.Id, &u.Uname, &u.Fname, &u.Lname, &u.Email); err != nil {
			return us, err
		}
		us = append(us, u)
	}
	return us, rows.Err()
}

func (d *DB) UserById(id int64) (*data.User, error) {
	return firstUser(d.usersWhere(util.IdEq(id)))
}

func (d *DB) UserByUsername(uname string) (*data.User, error) {
	return firstUser(d.usersWhere("username=$1", uname))
}

func (d *DB) UserByEmail(email string) (*data.User, error) {
	return firstUser(d.usersWhere("email=$1", email))
}

// UserSiteInfoByID gets the basic info for a user along with the user's role for the site.
// This function always returns a non-nil *User (which has just zero values for all fields if an error occurs).
// The error will be sql.ErrNoRows if a user is not found.
func (d *DB) UserSiteInfoByID(siteID, userID int64) (*data.User, error) {
	u := &data.User{Id: userID}
	q := "SELECT a.username,a.email,a.pass,a.fname,a.lname,IFNULL(b.v,'') FROM " + data.UsersTable + " AS a LEFT JOIN " + data.UserMetaTable +
		" AS b ON (a.id=b.user_id AND b.k=" + db.SingleQuote(SiteRoleKey(siteID)) + ") WHERE a.id=" + strconv.FormatInt(userID, 10)
	rows, err := db.Query(q)
	if err != nil {
		u.Id = 0
		return u, err
	}
	defer rows.Close()
	if rows.Next() {
		return u, rows.Scan(&u.Uname, &u.Email, &u.Pass, &u.Fname, &u.Lname, &u.Role)
	}
	u.Id = 0
	return u, sql.ErrNoRows
}

// UserSiteInfoByUsername gets the basic info for a user along with the user's role for the site.
// This function always returns a non-nil *User (which has just zero values for all fields if an error occurs).
// The error will be sql.ErrNoRows if a user is not found.
func (d *DB) UserSiteInfoByUsername(s *data.Site, uname string) (*data.User, error) {
	u := &data.User{Uname: strings.ToLower(uname)}
	q := "SELECT a.id,a.email,a.pass,a.fname,a.lname,IFNULL(b.v,'') FROM " + data.UsersTable + " AS a LEFT JOIN " + data.UserMetaTable +
		" AS b ON (a.id=b.user_id AND b.k=$1) WHERE a.username=$2"
	rows, err := d.db.Query(q, SiteRoleKey(s.Id), u.Uname)
	if err != nil {
		u.Uname = ""
		return u, err
	}
	defer rows.Close()
	if rows.Next() {
		return u, rows.Scan(&u.Id, &u.Email, &u.Pass, &u.Fname, &u.Lname, &u.Role)
	}
	u.Uname = ""
	return u, sql.ErrNoRows
}

// UserSiteInfoByEmail gets the basic info for a user along with the user's role for the site.
// This function always returns a non-nil *User (which has just zero values for all fields if an error occurs).
// The error will be sql.ErrNoRows if a user is not found.
func (d *DB) UserSiteInfoByEmail(s *data.Site, email string) (*data.User, error) {
	u := &data.User{Email: email}
	q := "SELECT a.id,a.username,a.pass,a.fname,a.lname,IFNULL(b.v,'') FROM " + data.UsersTable + " AS a LEFT JOIN " + data.UserMetaTable +
		" AS b ON (a.id=b.user_id AND b.k=$1) WHERE a.email=$2"
	rows, err := d.db.Query(q, SiteRoleKey(s.Id), email)
	if err != nil {
		u.Email = ""
		return u, err
	}
	defer rows.Close()
	if rows.Next() {
		return u, rows.Scan(&u.Id, &u.Uname, &u.Pass, &u.Fname, &u.Lname, &u.Role)
	}
	u.Email = ""
	return u, sql.ErrNoRows
}

// UserInsert creates a new user record with the given details in u and the password.
func (d *DB) UserInsert(username string, email string, passHash []byte, fname string, lname string) (userID int64, err error) {
	err = d.db.QueryRow("INSERT INTO "+data.UsersTable+" (username, email, pass, fname, lname) VALUES ($1,$2,$3,$4,$5) RETURNING id",
		username, email, passHash, fname, lname).Scan(&userID)
	return
}

// UserCountByUnameEmail says how many users there are with the given username and how many there are
// with the given email address.
func (d *DB) UserCountByUnameEmail(uname, email string) (unameCount int64, emailCount int64, err error) {
	err = d.db.QueryRow("SELECT (SELECT COUNT(*) FROM "+data.UsersTable+" WHERE username=$1),"+
		"(SELECT COUNT(*) FROM "+data.UsersTable+" WHERE email=$2)", uname, email).Scan(&unameCount, &emailCount)
	return
}

func (d *DB) UserPassword(userID int64) (hashedPass []byte, err error) {
	err = d.db.QueryRow("SELECT pass from " + data.UsersTable + " WHERE " + util.IdEq(userID)).Scan(&hashedPass)
	return
}

// firstUser returns the first *User from the given slice, or an error if err != nil or if the slice is empty.
func firstUser(us []data.User, err error) (*data.User, error) {
	if err != nil {
		return nil, err
	}
	if len(us) == 0 {
		return nil, sql.ErrNoRows
	}
	return &us[0], nil
}
