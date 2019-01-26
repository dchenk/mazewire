package users

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/dchenk/mazewire/pkg/data"
	"github.com/dchenk/mazewire/pkg/util"
)

// The following are all of the possible roles a user may have for a site.
const (
	RoleOwner      = "owner"
	RoleAdmin      = "admin"
	RoleEditor     = "editor"
	RoleAuthor     = "author"
	RoleSubscriber = "subscriber"
	RoleNone       = ""

	// RoleSuper is special because it is not a role given to anyone for any site.
	RoleSuper = "super"
)

// RoleAtLeast checks if the user's role (actualRole) is at least minRole. The special role "super" may also be passed in
// as actualRole or minRole.
func RoleAtLeast(actualRole, minRole string) bool {
	switch actualRole {
	case RoleSuper:
		return true
	case RoleOwner:
		return minRole == RoleOwner || minRole == RoleAdmin || minRole == RoleEditor || minRole == RoleAuthor || minRole == RoleSubscriber
	case RoleAdmin:
		return minRole == RoleAdmin || minRole == RoleEditor || minRole == RoleAuthor || minRole == RoleSubscriber
	case RoleEditor:
		return minRole == RoleEditor || minRole == RoleAuthor || minRole == RoleSubscriber
	case RoleAuthor:
		return minRole == RoleAuthor || minRole == RoleSubscriber
	case RoleSubscriber:
		return minRole == RoleSubscriber
	default:
		return false
	}
}

// SiteRole gives the user's role on the site with the given ID.
// The returned role is RoleNone if the user has no role on the site or if an error occurs.
//
// If a user is not found, this function checks if the user is Super: if the user is super, role is RoleSuper but
// otherwise the returned role is RoleNone and no error is returned.
// This function never returns a sql.ErrNoRows error.
func SiteRole(userID, siteID int64) (string, error) {
	opt, err := data.Conn.UserMetaByIdKey(userID, SiteRoleKey(siteID))
	if err != nil {
		if err == sql.ErrNoRows {
			// The user may be a super admin.
			if IsSuper(userID) {
				return RoleSuper, nil
			}
			err = nil
		}
		return RoleNone, err
	}

	// The following block does two things: It ensures that the role retrieved is a defined role in the app,
	// and it eliminates the allocation of a string to be garbage collected at the end of the request.
	switch string(opt.V) {
	case RoleOwner:
		return RoleOwner, err
	case RoleAdmin:
		return RoleAdmin, err
	case RoleEditor:
		return RoleEditor, err
	case RoleAuthor:
		return RoleAuthor, err
	case RoleSubscriber:
		return RoleSubscriber, err
	default:
		return RoleNone, err
	}
}

// SiteRoleKey returns the string that should be used to associate a user with a particular role
// on a site in the user meta table.
func SiteRoleKey(siteID int64) string {
	return "role" + strconv.FormatInt(siteID, 10)
}

// SetSiteRole gives the user the newRole for the site. If the user has no role on that site, the user is given one.
func SetSiteRole(userId int64, siteId int64, newRole string) error {
	if !RoleIsValid(newRole) {
		return fmt.Errorf("users: cannot set invalid role %q for user with ID %d", newRole, userId)
	}
	_, err := data.Conn.UserMetaUpdate(userId, SiteRoleKey(siteId), []byte(newRole))
	return err
}

// IsSuper says if the user is has super privileges.
// If an error occurs getting the user's role, the returned value is false.
func IsSuper(userID int64) bool {
	if userID == 0 {
		return false
	}
	role, err := data.Conn.UserMetaV(userID, "superuser")
	return err == nil && string(role) == "y"
}

// siteRolesAllowed are the kinds of roles users can have for sites
var siteRolesAllowed = [6]string{RoleOwner, RoleAdmin, RoleEditor, RoleAuthor, RoleSubscriber, RoleNone}

// RoleIsValid says if the role is of a kind that may be given to a user.
//
// There is no way to give a user the "super" role through the app.
func RoleIsValid(role string) bool {
	return util.SliceContainsString(siteRolesAllowed[:], role)
}
