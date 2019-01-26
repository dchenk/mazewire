package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dchenk/mazewire/pkg/data"
	"github.com/dchenk/mazewire/pkg/log"
	"github.com/dchenk/mazewire/pkg/users"
	"github.com/dchenk/mazewire/pkg/util"
)

// ReqPagepostList: GET pagepost
type ReqPagepostList struct {
	Site    int64  `msgp:"site"`    // the site ID, defaults to current host
	PpType  string `msgp:"pp_type"` // either "pages" or "posts", defaults to pages
	Offset  uint64 `msgp:"offset"`  // the pagination offset; defaults to 0
	Trashed bool   `msgp:"trashed"` // whether to get trashed items, otherwise defaults to getting "published", "draft", and "unsaved" items
}

// authorized checks just if the user is at least an author on the current site. A further check needs to ensure
// that the user can make changes to non-current sites.
func (*ReqPagepostList) authorized(_ *http.Request, _ *data.Site, u *data.User) bool {
	return users.RoleAtLeast(u.Role, users.RoleAuthor)
}

// List either pages or posts (for the admin area), paged by 20 results.
func (req *ReqPagepostList) handle(r *http.Request, s *data.Site, u *data.User) *APIResponse {

	if req.Site == 0 {
		req.Site = s.Id // Site defaults to the current host.
	} else {
		role, err := users.SiteRole(u.Id, req.Site)
		if err != nil {
			log.Err(r, "could not get logged in site role for user", err)
			return errProcessing()
		}
		if !users.RoleAtLeast(role, users.RoleAuthor) {
			return errLowPrivileges()
		}
	}

	respBody := RespPagepostList{
		Offset: req.Offset,
	}

	var statuses []string
	if req.Trashed {
		statuses = []string{"trashed"}
	} else {
		statuses = []string{"published", "draft", "unsaved"}
	}

	// The authorCheck variable is 0 only if the user making this request has higher than author privileges,
	// meaning that all pages/posts will be retrieved.
	var authorCheck int64
	if u.Role == users.RoleAuthor {
		authorCheck = u.Id
	}

	doneChan := make(chan bool, 1) // semaphore

	// Get the first 20 parent-level pages or posts, ordered by title and offset.
	var mainList []data.Content
	go func() {
		var err error // Shadow to allow changing of outer pagesposts.
		mainList, err = data.ContentsList(req.Site, req.PpType, []int64{0}, statuses, authorCheck, req.Offset)
		if err != nil {
			log.Err(r, "error with pagesposts ContentsList query", err)
			doneChan <- false
			return
		}
		doneChan <- true
	}()

	var err error

	// Count matching total and parent-level pages or posts.
	_, respBody.Count, err = data.CountContent(req.Site, req.PpType, statuses, authorCheck)
	if err != nil {
		log.Err(r, "error getting counts of pagesposts; user is "+u.Uname, err)
		<-doneChan // Empty the channel.
		return errProcessing()
	}

	// Wait for row counting goroutine to finish.
	if !<-doneChan {
		return errProcessing()
	}

	// respBody.Items will contain filtered pages/posts with their children.
	respBody.Items = make([]PagepostListing, 0, len(mainList)) // mainList approximates the final length

	// Check which if any parent and child pages/posts need to be retrieved.
	var parentsNeedingChildren []int64
	for _, aPagepost := range mainList {

		respBody.Items = append(respBody.Items, PagepostListing{
			Content:  &aPagepost,
			Children: make([]data.Content, 0, 1),
		})

		// This item can have children only if it does not have a parent.
		// (No items should have the home page as parent.)
		if aPagepost.Parent == 0 {
			parentsNeedingChildren = append(parentsNeedingChildren, aPagepost.Id)
		}

	}

	// Get all child pages/posts.
	if len(parentsNeedingChildren) > 0 {

		rowsChildren, err := data.Conn.ContentsList(req.Site, req.PpType, parentsNeedingChildren, statuses, authorCheck, 0)
		if err != nil {
			log.Err(r, "error getting rowsChildren", err)
			return errProcessing()
		}

		// Nest each child element under the appropriate parent.
		for i := range rowsChildren {
			for parentI := range respBody.Items {
				if respBody.Items[parentI].Id == rowsChildren[i].Parent {
					respBody.Items[parentI].Children = append(respBody.Items[parentI].Children, rowsChildren[i])
					break
				}
			}
		}

	}

	return &APIResponse{Body: &respBody}
}

type RespPagepostList struct {
	Items  []PagepostListing `msgp:"items"`
	Count  uint32            `msgp:"count"`
	Offset uint64            `msgp:"offset"`
}

// ReqPagepostCreate: POST pagepost
type ReqPagepostCreate struct {
	Site   int64  `msgp:"site"`    // the site ID; defaults to ID of current host
	Slug   string `msgp:"slug"`    // the URL formatted slug of the page or post itself, without slashes, maximum 80 characters
	PpType string `msgp:"pp_type"` // either "page" or "post"; defaults to "page"
	Parent int64  `msgp:"parent"`  // optionally set this new page under a parent page; defaults to 0
	Title  string `msgp:"title"`   // title of the page or post, maximum 255 characters
}

// authorized checks just if the user has at least author role on the current site.
func (req *ReqPagepostCreate) authorized(_ *http.Request, _ *data.Site, u *data.User) bool {
	return users.RoleAtLeast(u.Role, users.RoleAuthor)
}

// handle creates either a page or a post.
func (req *ReqPagepostCreate) handle(r *http.Request, s *data.Site, u *data.User) *APIResponse {

	if req.Site == 0 {
		req.Site = s.Id // Site defaults to the current host.
	} else {
		role, err := data.Conn.SiteRole(req.Site)
		if err != nil {
			log.Err(r, "could not get logged in site role for user", err)
			return errProcessing()
		}
		if !users.RoleAtLeast(role, users.RoleAuthor) {
			return errLowPrivileges()
		}
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Slug = strings.TrimSpace(req.Slug)

	if req.Title == "" {
		return APIResponseErr("You must provide a title.")
	}

	// The user may try to create the home page. This is checked for below.
	if req.Slug != "/" && !util.ValidPathSlug(req.Slug) {
		e := "You must provide a valid URL path slug (between 2 and 80 characters). You cannot have special characters at the ends of the slug, " +
			"and you cannot have these characters at all: '?', '/', ':', '&', '%', '!', '#', '%', '^', '(', ')', '='"
		return APIResponseErr(e)
	}

	if slugIsDisallowed(req.Slug) {
		return APIResponseErr("This slug means something special and is not allowed. Please pick something else.")
	}

	if req.PpType != "page" && req.PpType != "post" {
		log.Err(r, fmt.Sprintf("got type %q", req.PpType), errors.New("unexpected page/post type"))
		return errProcessing()
	}

	// Check if a page or post with that slug exists.
	count, err := data.ContentCountSlug(req.Site, req.Slug)
	if err != nil {
		return errProcessing()
	}
	if count > 0 {
		return APIResponseErr("A page or article with that slug already exists on this site.")
	}

	newID, err := data.ContentInsert(req.Site, req.Slug, u.Id, req.PpType, req.Parent, req.Title)
	if err != nil {
		log.Err(r, "error inserting pagepost", err)
		return errProcessing()
	}
	if newID == 0 { // Sanity check just in case.
		log.Err(r, "got 0 for newID for page but no error", errors.New("bad response from data package"))
		return errProcessing()
	}

	return &APIResponse{Body: RespPagepostCreate{newID}}
}

type RespPagepostCreate struct {
	NewId int64 `msgp:"new_id"`
}

// slugIsDisallowed says if the given slug is in the disallowed slugs list.
func slugIsDisallowed(slug string) bool {
	for i := range disallowedContentSlugs {
		if disallowedContentSlugs[i] == slug {
			return true
		}
	}
	return false
}

var disallowedContentSlugs = []string{"api", "content", "admin", ".well-known"} // Check for homepage "/" separately.

// ReqPagepostMakeDynElem: POST pagepost/element
type ReqPagepostMakeDynElem struct {
	// TODO
}

// authorized says just if the user has at least author role on the current site.
func (req *ReqPagepostMakeDynElem) authorized(r *http.Request, s *data.Site, u *data.User) bool {
	return users.RoleAtLeast(u.Role, users.RoleAuthor)
}

func (req *ReqPagepostMakeDynElem) handle(r *http.Request, s *data.Site, u *data.User) *APIResponse {
	// TODO:
	return APIResponseErr("This feature is not yet available")
}

// updateContentMeta updates the values for a Content record.
// If any error occurs, the update operation will be canceled by sending the error message
// to the "errs" channel and the function returning immediately.
func updateContentMeta(r *http.Request, u *data.User, s *data.Site, existing *data.Content, values map[string]string, errs chan string, wg *sync.WaitGroup) {

	defer wg.Done()

	args := make(map[string]interface{}) // the new values

	for _, col := range contentRecordCols {
		if val, ok := values[col]; ok {
			includeCol := false // Indicate whether to include the column in the save.
			switch col {
			case "slug":
				if util.ValidPathSlug(val) || slugIsDisallowed(val) {
					includeCol = true
				} else {
					errs <- "The URL slug you provided is invalid."
					return
				}
			case "title":
				if util.IsAnyStringBlank(val) || len(val) > 255 {
					errs <- "The page title must be between 1 and 255 characters long."
					return
				} else {
					includeCol = true
				}
			case "author":
				uIdStr := strconv.FormatInt(u.Id, 10) // string ID of current user
				if uIdStr == val {                    // Check if the new author is the user making the API request.
					// We already checked that the current user has at least author role on the site; just check
					// if the current user is already the author of the page (to avoid extra work for the DB).
					if existing.Author == u.Id {
						includeCol = false
					}
				} else {
					// Check if the user being given authorship has at least "author" role on the site being edited now.
					parsedAuthorID, err := strconv.ParseInt(strings.TrimSpace(val), 10, 64)
					if err != nil {
						log.Err(r, fmt.Sprintf("could not extract new author ID from string %q", val), err)
						errs <- errProcessingMsg
						return
					}
					newAuthor, err := data.UserSiteInfoByID(s.Id, parsedAuthorID)
					if err != nil {
						log.Err(r, "could not get new author info to update page meta", err)
						errs <- errProcessingMsg
						return
					}
					// Check if the person who is supposed to be the new author has at least "author" role on the site.
					if users.RoleAtLeast(newAuthor.Role, users.RoleAuthor) {
						includeCol = true
					} else {
						errs <- "The new author must have at least 'author' role on this site."
						return
					}
				}
			case "meta_title", "meta_desc":
				if len(val) <= 255 {
					includeCol = true
				} else {
					errs <- "The page title and meta description must be no more than 255 characters long."
					return
				}
			case "status":
				// This function can be used to set the page status only to either "draft" or "trashed".
				includeCol = (val == "draft" || val == "trashed") && val != existing.Status
			case "updated": // Here "updated" means what the page publish timestamp will be displayed.
				// Check if val is a valid timestamp and not later than the current time.
				t, err := util.ParseTimestamp(val)
				if !util.ValidTimestamp(val) || err != nil {
					errs <- "There is a problem with the publish timestamp you provided."
					return
				}
				// Check the timestamp, adjusting the server's timestamp up to the maximal timezone offset from the user.
				if t.After(time.Now().Add(time.Hour * 13)) {
					errs <- "You cannot provide a timestamp that is after the current time."
					return
				}
				includeCol = true
			}
			if includeCol {
				args[col] = val
			}
		}
	}

	if len(args) > 0 {
		if _, err := data.Conn.ContentUpdate(existing.Id, args); err != nil {
			log.Err(r, "error updating page meta", err)
			errs <- errProcessingMsg
		}
	}

}

// contentRecordCols lists each of the page meta options that could be submitted with requests to change a page's meta data.
// These are also the column names in the DB for pages and posts.
var contentRecordCols = [7]string{"author", "title", "slug", "meta_title", "meta_desc", "status", "updated"}

// The PagepostListing type represents a page or post listing element that may have child pages or posts.
type PagepostListing struct {
	*data.Content `msgp:"content"`
	Children      []data.Content `msgp:"children"`
}
