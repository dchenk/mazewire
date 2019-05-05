package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/dchenk/mazewire/pkg/data"
	"github.com/dchenk/mazewire/pkg/log"
	"github.com/dchenk/mazewire/pkg/roles"
	"github.com/dchenk/mazewire/pkg/room"
	"github.com/dchenk/mazewire/pkg/util"
)

// var siteApiActions = map[string]func(*http.Request, *data.Site, *data.User, *APIResponse){
// 	"create":             siteCreate,
// 	"change_home":        siteChangeHome,
// 	"set-author-display": siteAuthorDisplay,
// 	"logout-redir":       siteLogoutRedir,
// }

type SiteCreate struct {
	Domain string
	Name   string
}

// authorized says if the current user is authorized to create a website.
// TODO: the user must be logged in (at least subscriber) and have at least admin role
// on the current site. If they are logged in but don't have admin role, they should be
// redirected to the main site to create the site there.
func (*SiteCreate) authorized(_ *http.Request, s *data.Site, u *data.User) bool {
	return roles.RoleAtLeast(u.Role, roles.Role_SUBSCRIBER) || s.Id == 1
}

// createWebsite creates a website and sets it up with the right tables and pages/posts.
func (sc *SiteCreate) handle(r *http.Request, s *data.Site, u *data.User) *APIResponse {

	// TODO: check user privileges beyond just checking that the person is a subscriber

	var err error
	sc.Domain, err = util.ExtractDomain(sc.Domain)
	if err != nil {
		return APIResponseErr("It looks like the domain name you provided is not valid.")
	}

	sc.Name = strings.TrimSpace(sc.Name)

	if util.IsAnyStringBlank(sc.Domain, sc.Name) {
		return APIResponseErr("Please provide both a domain name and a website name.")
	}

	// Check if a site with the desired domain already exists.
	_, err = data.Conn.SiteByDomain(sc.Domain)
	if err != nil && err != sql.ErrNoRows {
		return errProcessing()
	}
	if err == nil {
		return APIResponseErr("A site with this domain already exists")
	}

	newID, err := data.Conn.InsertSite(sc.Domain, sc.Name)
	if err != nil {
		log.Err(r, "error inserting site", err)
		return errProcessing()
	}

	var resp APIResponse

	// Create the home page. -- TODO: do all this together with one insert
	// TODO: Create a login page.
	// TODO: Create a sample blog post.
	if _, err := data.Conn.ContentInsert(newID, "/", u.Id, "page", 0, "Home"); err != nil {
		log.Err(r, "could not create home page", err)
		resp.warn("We could not create a homepage for you, but you can easily do that yourself.")
	}

	// respond with the ID of the new site
	resp.Body = &RespSiteCreate{NewId: newID}
	return &resp

}

type RespSiteCreate struct {
	NewId int64 `msgp:"new_id"`
}

// SiteChangeHome: POST
type SiteChangeHome struct {
	Site        int64  `msgp:"site"`     // the site ID; defaults to current host
	HomeNewID   uint32 `msgp:"new_id"`   // the ID of the new home page
	OldHomeSlug string `msgp:"old_slug"` // the new slug of the old home page
}

func (req *SiteChangeHome) authorized(r *http.Request, s *data.Site, u *data.User) bool {
	return roles.RoleAtLeast(u.Role, roles.Role_ADMIN)
}

func (req *SiteChangeHome) handle(r *http.Request, s *data.Site, u *data.User) *APIResponse {

	if req.Site == 0 {
		req.Site = s.Id // Site defaults to the current host
	} else {
		role, err := roles.SiteRole(u.Id, req.Site)
		if err != nil {
			log.Err(r, "could not get logged in site role for user", err)
			return errProcessing()
		}
		if !roles.RoleAtLeast(role, roles.Role_ADMIN) {
			return errLowPrivileges()
		}
	}

	// TODO
	return APIResponseErr("this endpoint is not set up yet")

	return &APIResponse{Body: &RespSiteChangeHome{true}}
}

type RespSiteChangeHome struct {
	Ok bool `msgp:"ok"`
}

// SiteDelete: DELETE site
type SiteDelete struct {
	SiteID int64 `msgp:"site_id"`
}

func (*SiteDelete) authorized(_ *http.Request, _ *data.Site, u *data.User) bool {
	return roles.RoleAtLeast(u.Role, roles.Role_OWNER)
}

func (*SiteDelete) handle(*http.Request, *data.Site, *data.User) *APIResponse { // TODO: finish, require email confirmation
	return APIResponseErr("This feature is not set up yet.")
}

// toggle displaying of page or post author when listing pages or posts in admin area
func siteAuthorDisplay(r *http.Request, s *data.Site, u *data.User) *APIResponse {
	return APIResponseErr("this endpoint is not set up yet")
}

func siteLogoutRedir(r *http.Request, s *data.Site, u *data.User) *APIResponse {
	return APIResponseErr("this endpoint is not set up yet")
}

// SiteGetTheme: GET site/theme
type SiteGetTheme struct {
	Site int64 `msgp:"site"`
}

// authorized checks just if the user has at least author role on the current site.
func (*SiteGetTheme) authorized(_ *http.Request, _ *data.Site, u *data.User) bool {
	return roles.RoleAtLeast(u.Role, roles.Role_AUTHOR)
}

// handle returns the current theme for a site. This response body is a room.Tree.
func (req *SiteGetTheme) handle(r *http.Request, s *data.Site, u *data.User) *APIResponse {
	if req.Site == 0 {
		req.Site = s.Id // Site defaults to the current host.
	} else {
		role, err := roles.SiteRole(u.Id, req.Site)
		if err != nil {
			log.Err(r, "could not get logged in site role for user", err)
			return errProcessing()
		}
		if !roles.RoleAtLeast(role, roles.Role_AUTHOR) {
			return errLowPrivileges()
		}
	}
	theme, err := siteMainTheme(s)
	if err != nil {
		log.Err(r, "could not get site main theme", err)
		return errProcessing()
	}
	return &APIResponse{Body: theme}
}

const siteWorkingThemeOption = "site-theme-working"

// siteMainTheme gets the working (not yet compiled) theme of the site.
// If an error occurs making the theme, the default theme for all sites is produced.
// This function never returns a nil room.Tree slice.
func siteMainTheme(s *data.Site) (room.Tree, error) {
	themeData, err := s.OptionByKey(siteWorkingThemeOption)
	if err != nil {
		return make(room.Tree, 0), err
	}
	var theme room.Tree
	_, err = theme.UnmarshalMsg(themeData.V)
	if err != nil {
		theme = make(room.Tree, 0)
	}
	return theme, err
}

const siteThemeOption = "site-theme"

// siteMainCompiledTheme sends the main compiled site theme of s down c.
// If an error occurs making the theme, the default theme for all sites is produced.
// This function never sends a nil pointer or a pointer to a nil room.CompiledTree slice down c.
func siteMainCompiledTheme(s *data.Site, c chan room.CompiledTree) {
	themeData, err := s.OptionByKey(siteThemeOption)
	if err != nil {
		c <- defaultSiteTheme()
		return
	}
	var theme room.CompiledTree
	_, err = theme.UnmarshalMsg(themeData.V)
	if err != nil {
		c <- defaultSiteTheme()
		return
	}
	c <- theme
}

const postThemeOption = "post-theme"

func sitePostTheme(s *data.Site) (room.CompiledTree, error) {
	themeData, err := s.OptionByKey(postThemeOption)
	if err != nil {
		return defaultPostTheme(), err
	}
	var theme room.CompiledTree
	_, err = theme.UnmarshalMsg(themeData.V)
	if err != nil {
		return defaultPostTheme(), fmt.Errorf("error unmarshalling post theme; %s", err)
	}
	return theme, nil
}

func defaultSiteTheme() room.CompiledTree {
	return room.CompiledTree{
		room.CompiledSection{
			Type: "standard",
			Rows: []room.CompiledRow{
				{
					Type: "fullwidth",
					Modules: [][]room.CompiledModule{
						{
							{
								Type: "body",
							},
						},
					},
				},
			},
		},
	}
}

func defaultPostTheme() room.CompiledTree {
	return room.CompiledTree{
		room.CompiledSection{
			Type: "standard",
			Rows: []room.CompiledRow{
				{
					Type: "fullwidth",
					Modules: [][]room.CompiledModule{
						{
							{
								Type: "post_body",
							},
						},
					},
				},
			},
		},
	}
}
