package main

import (
	"net/http"

	"github.com/dchenk/mazewire/pkg/data"
	"github.com/dchenk/mazewire/pkg/log"
	"github.com/dchenk/mazewire/pkg/users"
	"github.com/dchenk/msgp/msgp"
)

// adminSrcNames lists the short name (used in the URL) of each of the admin dashboard static source files,
// and the adminSrcFullNames array is constructed out of this array. The order of the names here determines
// the order in which the files are loaded in to the document.
var adminSrcNames = [4]string{"style", "manifest", "vendor", "app"}
var adminSrcFullNames = func() [4]string {
	var full [4]string
	for i := range adminSrcNames {
		full[i] = "admin-" + adminSrcNames[i]
	}
	return full
}()

// reqAdminGetSrcs: GET admin
type reqAdminGetSrcs struct{}

// authorized returns true for requests to view the latest admin source file versions.
func (reqAdminGetSrcs) authorized(r *http.Request, _ *data.Site, _ *data.User) bool {
	return true
}

// adminGetSrcs returns a map listing the current admin JS and CSS file hashes.
func (reqAdminGetSrcs) handle(r *http.Request, _ *data.Site, _ *data.User) *APIResponse {
	return &APIResponse{
		Body: msgp.Raw(msgp.AppendMapStrStr(make([]byte, 0, 128), getSrcVersions(r))),
	}
}

// ReqAdminUpdateSrcs: POST admin
type ReqAdminUpdateSrcs struct {
	Srcs map[string]string `msgp:"srcs"`
}

// authorized says if the user is a super user.
func (ReqAdminUpdateSrcs) authorized(_ *http.Request, _ *data.Site, u *data.User) bool {
	return users.IsSuper(u.Id)
}

// updateSrcVersions updates the names of the latest JS and CSS source files for the admin area app.
// Properties in the provided map (each is optional, and only provided ones are updated):
//  - admin-css
//  - admin-manifest (JS)
//  - admin-vendor (JS)
//  - admin-app (JS)
func (req ReqAdminUpdateSrcs) handle(r *http.Request, s *data.Site, u *data.User) *APIResponse {
	if len(req.Srcs) == 0 {
		return APIResponseErr("Missing the srcs param")
	}

outerLoop:
	for k := range req.Srcs {
		for _, name := range adminSrcFullNames {
			if k == name {
				break outerLoop
			}
		}
		// Ensure that we're not saving garbage.
		delete(req.Srcs, k)
	}

	if len(req.Srcs) == 0 {
		return APIResponseErr("No changes were submitted")
	}

	changed, err := data.Conn.OptionsUpdate(1, req.Srcs)
	if err != nil {
		log.Err(r, "error updating src versions", err)
		return errProcessing()
	}

	// The response is just the number of source file versions updated.
	return &APIResponse{
		Body: msgp.Raw(msgp.AppendInt64(make([]byte, 0, 9), changed)),
	}
}

// getSrcVersions retrieves the latest hashes of the admin static source files.
// This function always returns a map with the keys
// though the map may be empty if an error occurs getting the data.
func getSrcVersions(r *http.Request) map[string]string {
	opts, err := data.Conn.OptionsKeyInMappedStr(1, adminSrcFullNames[:])
	if err != nil || len(opts) == 0 {
		if err != nil {
			log.Err(r, "error getting source versions", err)
		}
		// TODO: fallback to default (updated with each deployment automatically)
		return opts
	}
	return opts
}
