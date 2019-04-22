package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/dchenk/mazewire/pkg/data"
	"github.com/dchenk/mazewire/pkg/filters"
	filter_payloads "github.com/dchenk/mazewire/pkg/filters/payloads"
	"github.com/dchenk/mazewire/pkg/log"
	"github.com/dchenk/mazewire/pkg/plugins"
	"github.com/dchenk/mazewire/pkg/room"
	"github.com/dchenk/mazewire/pkg/users"
	"github.com/golang/protobuf/proto"
)

const (
	// pageBodyRole is the value in the "role" column that Blob records holding page tree structures have.
	pageBodyRole = "body"

	// pageBodyProdRole is the value in the "role" column that Blob records holding compiled page tree structures have.
	pageBodyProdRole = "body_prod"

	// pageBodyProdDynRole is the value in the "role" column that Blob records holding compiled page tree structures
	// with dynamic elements have.
	pageBodyProdDynRole = "body_prod_dyn"
)

// ReqPageEditContent: GET page-edit
type ReqPageEditContent struct {
	Site int64 `msgp:"site"` // the site ID; defaults to current host
	Page int64 `msgp:"page"` // the page ID
}

func (req *ReqPageEditContent) authorized(r *http.Request, s *data.Site, u *data.User) bool {
	return users.RoleAtLeast(u.Role, users.Role_AUTHOR)
}

// Get the content associated with a particular page, including all sections, rows, and modules of the page.
func (req *ReqPageEditContent) handle(r *http.Request, s *data.Site, u *data.User) *APIResponse {

	role := u.Role
	if req.Site == 0 {
		req.Site = s.Id // Site defaults to the current host.
	} else {
		var err error
		role, err = users.SiteRole(u.Id, req.Site)
		if err != nil {
			log.Err(r, "could not get logged in site role for user", err)
			return errProcessing()
		}
		if !users.RoleAtLeast(role, users.Role_AUTHOR) {
			return errLowPrivileges()
		}
	}

	content, err := data.ContentByID(req.Page)
	if err != nil {
		log.Err(r, "could not get content details", err)
		return errProcessing()
	}

	// If the user has just author role on the site, the user must be the author of the page.
	if role == users.Role_AUTHOR && content.Author != u.Id {
		return errLowPrivileges()
	}

	// Get all of the versions of the page.
	pageVersions, err := s.BlobsByRoleK(pageBodyRole, req.Page)
	if err != nil {
		log.Err(r, "error getting page versions", err)
		return errProcessing()
	}

	// Initialize fields that must not be nil.
	body := RespPageEditContent{
		Content:      content,
		PageVersions: make(PageVersions, len(pageVersions)),
	}

	// Prepare part of the response with minimal data from treeVersions.
	for i := range pageVersions {
		v := &pageVersions[i]

		// Decode the data contents.
		if _, err = body.PageVersions[i].UnmarshalMsg(v.V); err != nil {
			log.Err(r, "could not unmarshal PageVersion", err)
			return errProcessing()
		}

		// The Id and Updated fields do not ge unmarshalled.
		body.PageVersions[i].Id = v.Id
		body.PageVersions[i].Timestamp = v.Updated
	}

	// Sort the versions by timestamp to identify the latest one.
	sort.Sort(body.PageVersions)

	if len(body.PageVersions) > 0 {

		// Unmarshall the room.Tree from the latest version's Data field.
		if _, err = body.Tree.UnmarshalMsg(body.PageVersions[0].Data); err != nil {
			log.Err(r, "error unmarshalling page tree", err)
			return errProcessing()
		}

		if dynDataIDs := body.Tree.DynamicDataIDs(); len(dynDataIDs) > 0 {
			dynData, err := s.BlobsIdIn(dynDataIDs)
			if err != nil {
				log.Err(r, "error getting dynamic data for tree", err)
				return errProcessing()
			}
			body.DynData = make(map[string]data.Blob, len(dynData))
			for i := range dynData {
				body.DynData[strconv.FormatInt(dynData[i].Id, 10)] = dynData[i]
			}
		}

	}

	return &APIResponse{Body: &body}
}

// RespPageEditContent is the response body for GET page-edit requests.
type RespPageEditContent struct {
	*data.Content `msgp:"content"`
	Tree          room.Tree            `msgp:"tree"`     // the last saved (current) content tree; represents type PageTree
	DynData       map[string]data.Blob `msgp:"dyn_data"` // data for dynamic elements within the tree
	PageVersions  `msgp:"versions"`    // the tree structure versions and other details
	UserCSS       string               `msgp:"user_css"` // copied from Pagepost so that client does not decode []byte to string
}

// ReqPageEditSave: POST page-edit
type ReqPageEditSave struct {
	Site         int64    `msgp:"site"` // the site ID; defaults to current host
	Page         int64    `msgp:"page"` // the page ID
	Tree         msgp.Raw `msgp:"tree"` // The entire new page structure (role = "body"); represents type room.Tree
	data.Content `msgp:"content"`
	PageVersion  `msgp:"version"` // The Id, Data, and Timestamp fields are not used.
}

// authorized checks just if the user is at least an author on the current site.
func (req *ReqPageEditSave) authorized(r *http.Request, s *data.Site, u *data.User) bool {
	return users.RoleAtLeast(u.Role, users.Role_AUTHOR)
}

// handle saving changes to a page.
func (req *ReqPageEditSave) handle(r *http.Request, s *data.Site, u *data.User) *APIResponse {

	role := u.Role
	if req.Site == 0 {
		req.Site = s.Id // Site defaults to the current host.
	} else {
		var err error
		role, err = users.SiteRole(u.Id, req.Site)
		if err != nil {
			log.Err(r, "could not get logged in site role for user", err)
			return errProcessing()
		}
		if !users.RoleAtLeast(role, users.Role_AUTHOR) {
			return errLowPrivileges()
		}
	}

	content, err := data.ContentByID(req.Page)
	if err != nil {
		log.Err(r, "could not get content details", err)
		return errProcessing()
	}

	// If the user has just "author" role on the site, the user must be the author of the page.
	if role == users.Role_AUTHOR && content.Author != u.Id {
		return APIResponseErr("You must be the author of this page to edit it.")
	}

	// Ensure that this API endpoint is only used for pages, not posts or anything else.
	if content.Type != "page" {
		log.Err(r, "the page-edit endpoint is being accessed to save a non-page",
			fmt.Errorf("got content type %q", content.Type))
		return errProcessing()
	}

	var wg sync.WaitGroup
	wg.Add(2)
	errs := make(chan string, 2) // Provide enough buffer for both functions to each write one error message.

	// If an error occurs updating the page meta, the update operation will be canceled by sending the error
	// message to the "errs" channel and the function returning immediately.
	go updateContentMeta(r, u, s, content, req.Meta, errs, &wg)

	// Get the last saved version of the body tree.
	var latestPageBody *data.Blob
	go func() {
		defer wg.Done()
		var err error
		latestPageBody, err = s.BlobByRoleKLast(pageBodyRole, req.Page)
		if err != nil && err != sql.ErrNoRows {
			log.Err(r, "error getting body tree for page", err)
			errs <- errProcessingMsg
		}
	}()

	wg.Wait()
	select {
	case e := <-errs:
		// Report the first error to the user.
		return APIResponseErr(e)
	default:
		// Good to continue
	}

	var latestVersion PageVersion

	// latestPageBody will be nil if no version has been saved yet.
	if latestPageBody != nil {
		if _, err = latestVersion.UnmarshalMsg(latestPageBody.V); err != nil {
			log.Err(r, "could not unmarshal latest page version body", err)
			return errProcessing()
		}
	}

	// Make sure that the MessagePack-encoded tree structure is validly encoded.
	{
		var t room.Tree
		if o, err := t.UnmarshalMsg(req.Tree); err != nil || len(o) != 0 {
			if err == nil {
				err = errors.New("remaining bytes exist after unmarshalling tree")
			}
			log.Err(r, "got badly formatted tree from client", err)
			return errProcessing()
		}
	}
	req.PageVersion.Data = req.Tree

	versionData, err := req.PageVersion.MarshalMsg(nil)
	if err != nil {
		log.Err(r, "could not marshall page version data", err)
		return errProcessing()
	}

	var respBody RespPageEditSave
	respBody.InsertTime = time.Now()

	respBody.InsertId, err = s.BlobInsertFull(pageBodyRole, req.Page, versionData, respBody.InsertTime)
	if err != nil {
		log.Err(r, "error inserting new page version", err)
		return errProcessing()
	}

	// TODO: also delete older versions list the ones that were deleted (or send an entire new list of PageVersions)

	return &APIResponse{Body: respBody}
}

type RespPageEditSave struct {
	InsertId   int64     `msgp:"insert_id"`
	InsertTime time.Time `msgp:"insert_dt"`
}

// ReqPageEditPublish: PATCH pagepost
type ReqPageEditPublish struct {
	Site int64 `msgp:"site"` // the site ID; defaults to current host
	Page int64 `msgp:"page"` // the page ID
}

func (req *ReqPageEditPublish) authorized(_ *http.Request, _ *data.Site, u *data.User) bool {
	return users.RoleAtLeast(u.Role, users.Role_AUTHOR)
}

// pageEditPublish compiles and publishes a page.
func (req *ReqPageEditPublish) handle(r *http.Request, s *data.Site, u *data.User) *APIResponse {

	// Get the user's role for the site being published to.
	role := u.Role
	if req.Site == 0 {
		req.Site = s.Id // Site defaults to the current host.
	} else {
		var err error
		role, err = users.SiteRole(u.Id, req.Site)
		if err != nil {
			log.Err(r, "could not get logged in site role for user", err)
			return errProcessing()
		}
		if !users.RoleAtLeast(role, users.Role_AUTHOR) {
			return errLowPrivileges()
		}
	}

	content, err := data.ContentByID(req.Page)
	if err != nil {
		log.Err(r, "could not get content details", err)
		return errProcessing()
	}

	// if the user has just "author" role on the site, the user must be the author of the page being published
	if role == users.Role_AUTHOR && content.Author != u.Id {
		return APIResponseErr("You must be the author of this page to publish it.")
	}

	latestPageBlob, err := s.BlobByRoleKLast(pageBodyRole, req.Page)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Err(r, "attempt to publish page with no version saved", err)
			return APIResponseErr("It looks like you are trying to publish a page where you have not saved anything yet.")
		}
		log.Err(r, "error getting pave version to publish", err)
		return errProcessing()
	}

	var latestVersion PageVersion
	if _, err = latestVersion.UnmarshalMsg(latestPageBlob.V); err != nil {
		log.Err(r, "could not unmarshal latest page version body", err)
		return errProcessing()
	}

	roomTreeModule := room.Module{
		Type: room.TreeType,
		Data: latestVersion.Data,
	}

	// Create a buffer for all of the CSS of the page.
	var pageCSS room.PageCSS

	// pt will be the compiled page.
	pt := make(room.Tree, 0, 8)

	compiledTree, allStatic, err := pt.Compile(&dataStore{s}, nil, &roomTreeModule, &pageCSS)
	if err != nil {
		log.Err(r, "could not compile a page tree", err)
		return errProcessing()
	}

	// The user's custom CSS code is written last.
	userCSSMessage := &filter_payloads.UserCSS{Code: latestVersion.CSS}
	userCSSResp, err := plugins.DoFilter(filters.UserCSS, userCSSMessage)
	if err != nil {
		log.Err(r, "could not filter user CSS", err)
		return errProcessing()
	}
	var ok bool
	userCSSMessage, ok = userCSSResp.(*filter_payloads.UserCSS)
	assertType(r, ok, "UserCSS", userCSSMessage)
	pageCSS.WriteString(userCSSMessage.Code)

	if allStatic {
		latestVersion.Static = true
		latestVersion.Data = compiledTree.Data
	} else {
		latestVersion.Data, err = proto.Marshal(compiledTree)
		if err != nil {
			log.Err(r, "error marshalling compiled tree", err)
			return errProcessing()
		}
	}

	versionData, err := latestVersion.MarshalMsg(nil)
	if err != nil {
		log.Err(r, "could not marshal latest version data", err)
		return errProcessing()
	}

	colUpdates := map[string]interface{}{
		"body":    versionData,
		"status":  "published",
		"updated": time.Now(),
	}

	_, err = data.ContentUpdate(req.Page, colUpdates)
	if err != nil {
		log.Err(r, "error updating pagepost to publish", err)
		return errProcessing()
	}

	// Respond with true.
	resp := APIResponse{
		Body: msgp.Raw(msgp.AppendBool(nil, true)),
	}
	return &resp
}
