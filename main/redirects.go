package main

import (
	"database/sql"
	"net/http"
	"net/url"

	"github.com/dchenk/mazewire/pkg/data"
	"github.com/dchenk/mazewire/pkg/util"
)

func loginRedirect(r *http.Request, s *data.Site, u *data.User) (redirectPath string) {
	redirectPath = r.FormValue("redir") // The login form can have a field set to the path to redirect to upon login.
	if util.IsAnyStringBlank(redirectPath) {
		// Identify which page the user should be redirected to.
		rp, _ := s.OptionV("login-redirect-path")
		redirectPath = string(rp)
		if redirectPath == "" {
			redirectPath = "/admin" // Default to site admin page.
		}
	}
	return
}

func logoutRedirect(r *http.Request, s *data.Site, u *data.User) (redirectPath string) {
	rp, err := s.OptionV("logout-page-slug")
	redirectPath = string(rp)
	if err != nil {
		log.Err(r, "could not get site logout redirect URL", err)
	}
	if util.IsAnyStringBlank(redirectPath) {
		// Go to homepage. ... or TODO: go to the main login page
		redirectPath = "/?msg=" + url.QueryEscape("You are logged out")
	} else {
		// Go to page specified.
		redirectPath = "/" + redirectPath + "?msg=" + url.QueryEscape("You are logged out")
	}
	return
}

// registerRedirect returns a complete redirect path where the user should go to register under the current site.
// TODO: offer redirection to a page for user to add themselves to the site (already a Mazewire user)
func registerRedirect(r *http.Request, s *data.Site, u *data.User) (redirectPath string) {
	rp, err := s.OptionV("register-page-slug")
	if err != nil && err != sql.ErrNoRows {
		log.Err(r, "could not get option register-page-slug", err)
	}
	redirectPath = string(rp)
	if util.IsAnyStringBlank(redirectPath) {
		// Go to homepage. ... or TODO: go to the main login page
		redirectPath = r.Host + "/register"
	} else {
		// Go to page specified.
		redirectPath = r.Host + "/" + redirectPath
	}
	return
}
