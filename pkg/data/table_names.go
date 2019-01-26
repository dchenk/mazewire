package data

import "strconv"

const (
	SitesTable        = "sites"
	ContentTable      = "content"
	UsersTable        = "users"
	UserMetaTable     = "user_meta"
	OptionsTable      = "options"
	MediaTable        = "media"
	SiteMessagesTable = "site_messages"
	UserMessagesTable = "user_messages"
)

// BlobsTable gives the name of the blobs table for the site with the given ID.
// If siteID is zero, then the name of the system's blobs table is returned.
func BlobsTable(siteID int64) string {
	const table = "blobs"
	if siteID == 0 {
		return table
	}
	return table + strconv.FormatInt(siteID, 10)
}
