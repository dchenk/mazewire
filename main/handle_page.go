package main

import (
	"bytes"
	"database/sql"
	"html"
	"net/http"
	"strconv"
	"strings"

	"github.com/dchenk/mazewire/pkg/data"
	"github.com/dchenk/mazewire/pkg/env"
	"github.com/dchenk/mazewire/pkg/log"
	"github.com/dchenk/mazewire/pkg/room"
	"github.com/dchenk/mazewire/pkg/users"
	"github.com/golang/protobuf/proto"
)

// handlePage provides the entire response for the requested page, post, or other type of content.
func handlePage(r *http.Request, slugs []string, s *data.Site, u *data.User) (int, *contentBuffers) {

	if slugs[0] == "admin" {
		return http.StatusOK, siteAdminPage(r, s, u)
	}

	var (
		reqSlug       string
		reqParentSlug string // the page's parent slug in the request
		siteThemeChan = make(chan room.Tree, 1)
	)

	go siteMainCompiledTheme(s, siteThemeChan)

	if len(slugs) == 1 {
		reqSlug = slugs[0]
	} else {
		reqSlug = slugs[1]
		reqParentSlug = slugs[0]
	}

	content, parentSlug, err := data.Conn.ContentBySiteSlug(s.Id, reqSlug)
	if err != nil {
		if err != sql.ErrNoRows { // A page/post with that slug doesn't exist.
			log.Err(r, "error looking up page", err)
		}
		if reqSlug == "login" { // If the user wants to log in but a login page doesn't exist, give a generic one.
			return http.StatusSeeOther, genericLoginPage(s, u) // TODO: this is wrong, you should not return StatusSeeOther
		}
		return http.StatusNotFound, sitePageNotFound(s, r.Host+r.RequestURI)
	}

	// Check if content exists at the requested slug and has the same parent as in the request.
	if parentSlug != reqParentSlug {
		return http.StatusNotFound, sitePageNotFound(s, r.Host+r.RequestURI)
	}

	siteTheme := <-siteThemeChan
	_ = siteTheme

	// If what is requested has type "page" then we already have the site theme, but if a type "post" is requested then
	// we need to also get the post theme.
	if content.Type == "post" { // TODO: are there other types to check for (maybe later)?
		postTheme, err := sitePostTheme(s)
		if err != nil {
			// TODO
		}
		// TODO
		_ = postTheme
	}

	var version PageVersion
	if _, err = version.UnmarshalMsg(content.Body); err != nil {
		log.Err(r, "could not unmarshall version blob data", err)
		// TODO
	}
	version.Data = content.Body // In case it's needed somewhere.

	var css room.PageCSS
	css.WriteString("<style>")

	// First, write the already compiled CSS.
	css.WriteString(version.CSS)

	var cb contentBuffers

	// The CSS from the body is written after the compiled and user's custom CSS but before the page-wide CSS in the head.
	setupBody(r, s, content, &version, u, &cb.body, &css)

	// The CSS from dynamic elements is written last.
	setupHead(r, s, content, &version, u, &cb.head, &css)

	// Write the CSS to the page if there is anything.
	if css.Len() > len("<style>") {
		cb.head.Write(css.Bytes())
		cb.head.WriteString("</style>")
	}

	return http.StatusOK, &cb

}

func setupHead(_ *http.Request, s *data.Site, content *data.Content, version *PageVersion, u *data.User, head *bytes.Buffer, _ *room.PageCSS) {
	head.WriteString("<title>")
	if content.MetaTitle == "" {
		head.WriteString(html.EscapeString(content.Title))
		head.WriteString(" &ndash; ")
		head.WriteString(html.EscapeString(s.Name))
	} else {
		head.WriteString(content.MetaTitle)
	}
	head.WriteString(`</title><link rel="shortcut icon" href="`)

	if s.Favicon == "" {
		head.WriteString(userContentSrc)
		head.WriteString("favicon.png")
	} else {
		head.WriteString(s.Favicon)
	}
	head.WriteString(`"><meta name="description" content="`)
	head.WriteString(html.EscapeString(content.MetaDesc))
	head.WriteString(`">`)
}

func setupBody(r *http.Request, s *data.Site, content *data.Content, version *PageVersion, u *data.User, body *bytes.Buffer, css *room.PageCSS) {
	pageId := "page-" + strconv.FormatInt(content.Id, 10)

	room.ElementOpenTag("div", pageId, []string{}, body)

	// The page container may set the maximum width of the standard sections directly inside.
	if rs := strings.TrimSpace(version.Styles["section_row_space"]); rs != "" {
		css.WriteString("#" + pageId + " > .room-section{width:" + rs + "px;}") // TODO: this must not be here.
	}

	ds := dataStore{s}

	// Unmarshal the room.Tree from version.Data.
	if version.Static {
		body.Write(version.Data)
	} else {
		var cm room.Module
		if err := proto.Unmarshal(version.Data, &cm); err != nil {
			log.Err(r, "error unmarshalling compiled tree", err)
			body.WriteString(errOccurredEnd)
			return
		}

		// Compile the dynamic page.
		ct := room.Tree{Sections: make([]*room.Section, 0, 4)}
		if err := ct.BuildHTML(&ds, nil, &cm, body, css); err != nil {
			log.Err(r, "error building to view", err)
			body.WriteString(errOccurredEnd)
			return
		}
	}

	body.WriteString("</div>")

	// TODO: get javascript source file URLs to link and compile the links to the files
	var javaScriptSrcs string
	body.WriteString(javaScriptSrcs)
}

// errOccurredEnd is the string to write just before exiting if responding with an error message.
const errOccurredEnd = "An error occurred loading the page.</div>"

func siteAdminPage(r *http.Request, s *data.Site, u *data.User) *contentBuffers {

	var cb contentBuffers
	cb.head.Grow(768)

	cb.head.WriteString(`<title>Admin Dashboard &ndash; `)
	cb.head.WriteString(s.Name)
	cb.head.WriteString(`</title><link rel="shortcut icon" href="`)
	cb.head.WriteString(userContentSrc)
	if s.Favicon == "" {
		cb.head.WriteString("favicon.png")
	} else {
		cb.head.WriteString(s.Favicon)
	}
	cb.head.WriteString(`">`)

	// Check for error or alert messages that should be shown to the user.
	msgs, err := getSiteMessages(s.Id, u.Role)
	if err != nil {
		log.Err(r, "error getting user's site messages", err)
		msgs = append(msgs, data.SiteMessage{Role: u.Role, K: errProcessingMsg, Message: "An error occurred loading this page."})
	}
	if u.Id == 0 {
		msgs = append(msgs, data.SiteMessage{Message: "You are not logged in."})
	}

	if len(msgs) > 0 { // todo make this called from within the app
		cb.head.WriteString("<script>const window.msgs=[")
		for i := range msgs {
			if i > 0 {
				cb.head.WriteByte(',')
			}
			cb.head.WriteString(strconv.Quote(msgs[i].Message))
		}
		cb.head.WriteString("];</script>")
	}

	currentSources := getSrcVersions(r)

	srcRoot := env.Vars()[env.VarAdminSrcRoot] // TODO: clean this up

	cb.head.WriteString(`<link rel="stylesheet" href="https://fonts.googleapis.com/icon?family=Material+Icons"><link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Muli:300,400,400i,600,700">`)

	cb.head.WriteString(`<link rel="stylesheet" href="`)
	cb.head.WriteString(srcRoot)
	cb.head.WriteString(`style.`)
	cb.head.WriteString(currentSources[adminSrcFullNames[0]])
	cb.head.WriteString(`.css">`)

	userImgBytes, err := data.Conn.UserMetaV(u.Id, "profile_img")
	if err != nil && err != sql.ErrNoRows {
		log.Err(r, "could not get user profile image", err)
	}
	userImg := string(userImgBytes)

	// Write out the activeSite object in the head.
	cb.head.WriteString("<script>window.activeSite={id:")
	cb.head.WriteString(strconv.FormatInt(s.Id, 10))
	cb.head.WriteString(",domain:")
	cb.head.WriteString(strconv.Quote(s.Domain))
	cb.head.WriteString(",name:")
	cb.head.WriteString(strconv.Quote(s.Name))
	cb.head.WriteString(",logo:")
	cb.head.WriteString(strconv.Quote(s.Logo))
	cb.head.WriteString(",favicon:")
	cb.head.WriteString(strconv.Quote(s.Favicon))
	cb.head.WriteString(",role:")
	cb.head.WriteString(strconv.Quote(u.Role))
	cb.head.WriteString(`};`)

	// Write out the userInfo object in the head.
	cb.head.WriteString("window.userInfo={id:")
	cb.head.WriteString(strconv.FormatInt(u.Id, 10))
	cb.head.WriteString(",fname:")
	cb.head.WriteString(strconv.Quote(u.Fname))
	cb.head.WriteString(",lname:")
	cb.head.WriteString(strconv.Quote(u.Lname))
	cb.head.WriteString(",img:")
	cb.head.WriteString(strconv.Quote(userImg))
	cb.head.WriteString(`};</script>`)

	cb.body.WriteString(`<div id="app"></div>`)

	// Order matters in this loop.
	// Start at 1 to not write the CSS link again.
	for i := 1; i < len(adminSrcFullNames); i++ {
		cb.body.WriteString(`<script src="`)
		cb.body.WriteString(srcRoot)
		cb.body.WriteString(adminSrcNames[i])
		cb.body.WriteByte('.')
		cb.body.WriteString(currentSources[adminSrcFullNames[i]])
		cb.body.WriteString(`.js"></script>`)
	}

	return &cb

}

// TODO: notFoundPage gives the complete response contentBuffers displaying the site's 404 page when no resource exists at the requested path
func sitePageNotFound(s *data.Site, path string) *contentBuffers {
	var cb contentBuffers
	cb.body.WriteString("Page not found at " + path)
	// TODO: build the appropriate page for the site
	return &cb
}

// genericLoginPage returns a form for the user to log in because the site doesn't have this page.
// This is a temporary feature for bootstrapping the system.
func genericLoginPage(s *data.Site, u *data.User) *contentBuffers {
	var cb contentBuffers

	cb.head.WriteString(`<title>Log In</title><link rel="shortcut icon" href="`)
	if s.Favicon == "" {
		cb.head.WriteString(userContentSrc + `favicon.png">`)
	} else {
		cb.head.WriteString(s.Favicon + `">`)
	}

	cb.head.WriteString("<style>#msg {color: #EF5350;}</style>")

	cb.body.WriteString("<h2>Log in to your account</h2>")

	if u.Role != users.RoleNone {
		cb.body.WriteString(`<h4 id="msg">You are already logged in as the user: ` + u.Uname + `</h4>`)
		cb.body.WriteString(`<form action="/api/user/logout" method="POST"><input type="submit" value="Log Out"></form>`)
	} else {
		cb.body.WriteString(`<h4 id="msg"></h4>
		<form action="/api/user/login-form" method="POST">
			<label for="user">Username or Email Address</label><br>
			<input type="text" name="user" id="user"><br>
			<label for="pass">Password</label><br>
			<input type="password" name="pass" id="pass"><br>
			<input type="submit" value="Log In">
		</form>
		<script>
			var msgR = new RegExp(/msg=([^&]*)/);
			const msgMatch = msgR.exec(window.location.search);
			if (msgMatch) {
				if (msgMatch[1] !== "") { document.getElementById("msg").innerHTML = decodeURIComponent(msgMatch[1]); }
			}
		</script>`)
	}

	return &cb
}
