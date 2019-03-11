package users

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/dchenk/mazewire/pkg/data"
)

// RoleAtLeast checks if the user's role (actualRole) is at least minRole.
func RoleAtLeast(actualRole, minRole Role) bool {
	switch actualRole {
	case Role_SUPER:
		return true
	case Role_OWNER:
		return minRole == Role_OWNER || minRole == Role_ADMIN || minRole == Role_EDITOR ||
			minRole == Role_AUTHOR || minRole == Role_SUBSCRIBER || minRole == Role_NONE
	case Role_ADMIN:
		return minRole == Role_ADMIN || minRole == Role_EDITOR || minRole == Role_AUTHOR ||
			minRole == Role_SUBSCRIBER || minRole == Role_NONE
	case Role_EDITOR:
		return minRole == Role_EDITOR || minRole == Role_AUTHOR || minRole == Role_SUBSCRIBER ||
			minRole == Role_NONE
	case Role_AUTHOR:
		return minRole == Role_AUTHOR || minRole == Role_SUBSCRIBER || minRole == Role_NONE
	case Role_SUBSCRIBER:
		return minRole == Role_SUBSCRIBER || minRole == Role_NONE
	case Role_NONE:
		return minRole == Role_NONE
	default:
		return false
	}
}

// SiteRole gives the user's role on the site with the given ID.
// The returned role is NONE if the user has no role on the site or if an error occurs.
//
// If the user's role property is not found for the site, this function checks if the user is a
// super user: if the user is super, the role is Role_SUPER but otherwise the returned role is NONE
// and no error is returned. This function never returns a sql.ErrNoRows error.
func SiteRole(userID, siteID int64) (Role, error) {
	result, err := data.Conn.UserMetaByIdKey(userID, SiteRoleKey(siteID))
	if err != nil {
		if err == sql.ErrNoRows {
			// The user may be a super admin.
			if IsSuper(userID) {
				return Role_SUPER, nil
			}

			// The user's role was not found.
			err = nil
		}
		return Role_NONE, err
	}

	// Ensures that the role retrieved is a real role.
	val, err := strconv.ParseInt(string(result.V), 10, 32)
	role, _ := ValidRoleInt64(val)
	return role, err
}

// SiteRoleKey returns the string that should be used to associate a user with a particular role
// on a site in the user meta table.
func SiteRoleKey(siteID int64) string {
	return "role" + strconv.FormatInt(siteID, 10)
}

// SetSiteRole gives the user the newRole for the site. If the user has no role on that site, the
// user is given one.
func SetSiteRole(userId int64, siteId int64, newRole int64) error {
	if _, ok := ValidRoleInt64(newRole); !ok {
		return fmt.Errorf("users: cannot set invalid role %q for user with ID %d",
			newRole, userId)
	}
	_, err := data.Conn.UserMetaUpdate(userId, SiteRoleKey(siteId),
		[]byte(strconv.FormatInt(newRole, 10)))
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

// RoleIsValid says if the role is of a kind that a user may have.
func ValidRoleStr(role string) (Role, bool) {
	val, ok := Role_value[role]
	return Role(val), ok
}

func ValidRoleInt64(role int64) (Role, bool) {
	switch role {
	case 11:
		return Role_SUPER, true
	case 9:
		return Role_OWNER, true
	case 7:
		return Role_ADMIN, true
	case 5:
		return Role_EDITOR, true
	case 3:
		return Role_AUTHOR, true
	case 1:
		return Role_SUBSCRIBER, true
	case 0:
		return Role_NONE, true
	default:
		return Role_NONE, false
	}
}
