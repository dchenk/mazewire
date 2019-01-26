package data

type SiteGetter interface {
	// SiteByDomain retrieves a site by its domain.
	SiteByDomain(domain string) (*Site, error)

	// SitesByIDs retrieves sites by their ID.
	SitesByIDs(ids []int64) ([]Site, error)
}

type SiteInserter interface {
	// InsertSite inserts a Site into the sites table and create all needed tables for the site in a transaction.
	// The int64 returned is the ID of the new site (the inserted row).
	// The domain passed in must have already been validated as a real domain name.
	InsertSite(domain, name string) (int64, error)
}

type SiteDeleter interface {
}

type SiteManager interface {
	SiteGetter
	SiteInserter
	SiteDeleter
}

type BlobGetter interface {
	// BlobsByRoleInK selects blobs where the role matches a role in the given list and the K equals the given k.
	// If not matching Blobs are found, an empty slice is returned but no error.
	BlobsByRoleInK(site int64, roles []string, k int64) ([]Blob, error)

	// BlobByRoleLikeLast returns the last Blob that is matched to the role column using the LIKE feature, ordered
	// by the record ID.
	BlobByRoleLikeLast(site int64, kPattern string) (*Blob, error)
}

type BlobInserter interface {
}

type BlobDeleter interface {
}

type BlobManager interface {
	BlobGetter
	BlobInserter
	BlobDeleter
}

type ContentGetter interface {
	// ContentByID returns a single Content by its ID.
	ContentByID(id int64) (*Content, error)

	// ContentsByIDs returns the Content rows selected by their ID.
	ContentsByIDs(IDs []int64) ([]Content, error)

	ContentsIDIn(ids []int64) ([]Content, error)

	ContentBySiteSlug(siteID int64, slug string) (*Content, string, error)

	ContentByAuthor(authorID int64) ([]Content, error)

	// ContentsList retrieves the specified Content items. The parent argument should be either nil to indicate that the
	// parent column should not be considered in the query, or a list of the parent IDs to which the returned items should belong.
	// The value of offset may never be negative, which is why its type is uint64.
	ContentsList(siteID int64, cType string, parents []int64, statuses []string, authorID int64, offset uint64) ([]Content, error)

	// CountContent says how many content pages/posts there are with the given parameters and how many such have no parent.
	// Optionally give a non-zero authorID to also filter by author.
	// The given statuses should not have any punctuation at all (should be already sanitized).
	CountContent(siteID int64, pType string, statuses []string, authorID int64) (countTotal uint16, countParentLevel uint32, err error)

	// ContentCountSlug returns the number of rows with the given site ID and slug.
	ContentCountSlug(siteID int64, slug string) (count int64, err error)
}

type ContentInserter interface {
	ContentUpdate(contentID int64, vals map[string]interface{}) (int64, error)
}

type ContentDeleter interface {
	DeleteContent(IDs []int64) (rowsAffected int, err error)
}

type ContentManager interface {
	ContentGetter
	ContentInserter
	ContentDeleter
}

type UserGetter interface {
	UserById(id int64) (*User, error)
	UserByUsername(username string) (*User, error)
	UserByEmail(email string) (*User, error)
	UsersByIDs(IDs []int64) ([]User, error)
}

type UserInserter interface {
	// UserInsert inserts a new User.
	// If there are hooks to the user insert operation, the failures in the hook handlers does not cancel the transaction.
	UserInsert(username string, email string, passHash []byte, fname string, lname string) (userID int64, err error)
}

type UserDeleter interface {
	UserDelete(ID int64) error
}

type UserManager interface {
	UserGetter
	UserInserter
	UserDeleter
}

type UserPasswordGetter interface {
	// UserPassword returns a user's hashed password from the database.
	UserPassword(userID int64) ([]byte, error)
}

type UserMetaGetter interface {
	UserMetaById(userID int64) ([]UserMeta, error)
	UserMetaByIdMapped(userID int64) (map[string][]byte, error)
	UserMetaByIdKey(userID int64, k string) (*UserMeta, error)
	UserMetaByIdLikeKey(userID int64, k string) ([]UserMeta, error)
	UserMetaV(userID int64, k string) (v []byte, err error)
}

type UserMetaInserter interface {
	UserMetaUpdate(userID int64, k string, v []byte) (int64, error)
	UserMetasUpdate(ums []UserMeta) (int64, error)
}

type UserMetaDeleter interface {
	UserMetaDelete(userID int64, k string) (int64, error)
}

type UserMetaManager interface {
	UserMetaGetter
	UserMetaInserter
	UserMetaDeleter
}

type OptionGetter interface {
	// OptionByKey returns a site's option selected by its K.
	OptionByKey(site int64, k string) (*Option, error)

	// OptionsLikeKey returns options selected by their K being like k, as matched using the LIKE string
	// pattern matching feature in the database system.
	//
	// The syntax valid for the k string ultimately depends on the database implementation in use, but as
	// a general guideline the core application here should only use pattern strings that are supported in
	// the LIKE feature of all the major DBMSs, such as PostgreSQL (and CockroachDB) and MySQL.
	OptionsLikeKey(site int64, k string) ([]Option, error)

	// OptionsKeyIn returns a site's options selected by the Ks IN list.
	// The strings in the Ks slice must be valid UTF-8 strings.
	OptionsKeyIn(site int64, Ks []string) ([]Option, error)

	// OptionsKeyInMapped returns a map of the K-V pairs of a site's options selected by the Ks IN list.
	// This function always returns a non-nil map, though the map may be empty.
	// An error is not returned if there are no records retrieved.
	OptionsKeyInMapped(site int64, Ks []string) (map[string][]byte, error)

	// OptionsKeyInMappedStr returns a map of the K-V pairs of a site's options selected by the Ks In list,
	// with the V values converted to strings.
	// This function always returns a non-nil map, though the map may be empty.
	// An error is not returned if there are no records retrieved.
	OptionsKeyInMappedStr(site int64, Ks []string) (map[string]string, error)

	// OptionV returns the V of a site's option selected by key.
	OptionV(site int64, k string) ([]byte, error)
}

type OptionInserter interface {
	// OptionUpdate updates an Option or creates a new record in the database if necessary.
	//
	// The primary key is set by both the site ID and the value of K.
	OptionUpdate(site int64, k string, v []byte) (rowsAffected int64, err error)

	// OptionUpdateStr updates an Option or creates a new record in the database if necessary.
	OptionUpdateStr(site int64, k string, v string) (rowsAffected int64, err error)

	// OptionsUpdate updates Option records or creates new records in the database if necessary.
	// The strings passed in as keys in the map must be valid UTF-8 strings.
	// Returned is the number of rows affected.
	OptionsUpdate(site int64, opts map[string]string) (int64, error)
}

type OptionDeleter interface {
	// OptionDelete deletes a site option selected by its K.
	OptionDelete(site int64, k string) (rowsAffected int64, err error)
}

type OptionManager interface {
	OptionGetter
	OptionInserter
	OptionDeleter
}
