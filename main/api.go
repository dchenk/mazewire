package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dchenk/mazewire/pkg/data"
	"github.com/dchenk/mazewire/pkg/env"
	"github.com/dchenk/mazewire/pkg/log"
	"github.com/dchenk/mazewire/pkg/util"
)

const (
	errProcessingMsg  = "An error occurred when processing your request."
	errIncompleteForm = "Not all required fields are filled out. Please complete the form."
	errNoMembership   = "You do not yet have membership on this site."
	errInvalidLogin   = "The login details you entered are incorrect."
	errDecodingJSON   = "error decoding JSON req"
	errDecodingMsgp   = "error decoding msgp req"
)

func handleAPI(w http.ResponseWriter, r *http.Request, s *data.Site, u *data.User, slugs []string) {
	// The request path may have up to two parts: The first is the API endpoint, and the second is
	// the (optional) API sub-endpoint.
	if len(slugs) < 1 {
		writeBadAPIReq(w)
		return
	}
	ep, ok := apiEndpoints[slugs[0]]
	if !ok {
		writeBadAPIReq(w)
		return
	}

	// At this point, we verified that the request is likely going to a valid endpoint.
	// For the development environment, let the browser make requests to the local server.
	if !env.Prod() {
		w.Header().Add("Access-Control-Allow-Origin", "http://localhost:8082")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			return
		}
	}

	var handler APIHandler

	if len(slugs) == 1 {
		handler = ep.Points(r.Method)
	} else {
		if ep.SubPaths == nil {
			writeBadAPIReq(w)
			return
		}
		// Above we already checked that there are at least two parts in the slug, so here we know that there
		// there are at least three parts in the slug.
		ep, ok = ep.SubPaths[slugs[1]]
		if !ok {
			writeBadAPIReq(w)
			return
		}
		handler = ep.Points(r.Method)
	}

	if handler == nil {
		writeBadAPIReq(w)
		return
	}

	if r.Method == http.MethodGet {
		params, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			log.Err(r, "could not parse query for GET request", err)
			writeApiReqErr(w, http.StatusInternalServerError, errProcessingMsg)
			return
		}
		if queryData := params.Get("data"); queryData != "" {
			decoded, err := hex.DecodeString(queryData)
			if err != nil {
				log.Err(r, fmt.Sprintf("could not hex-decode data in GET request; RawQuery was %q", r.URL.RawQuery), err)
				writeApiReqErr(w, http.StatusInternalServerError, errProcessingMsg)
				return
			}
			decoder, ok := handler.(msgp.Decoder)
			if !ok {
				log.Err(r, "got content-type MessagePack but handler is not a decoder", unexpectedContentType)
				writeApiReqErr(w, http.StatusBadRequest, errProcessingMsg)
				return
			}
			if err := msgp.Decode(bytes.NewBuffer(decoded), decoder); err != nil {
				log.Err(r, errDecodingMsgp, err)
				writeApiReqErr(w, http.StatusInternalServerError, errProcessingMsg)
				return
			}
		}
	} else {
		var contentType string
		if cth := r.Header["Content-Type"]; len(cth) > 0 {
			contentType = cth[0]
		}

		switch contentType {
		case util.ContentTypeMessagePack:
			decoder, ok := handler.(msgp.Decoder)
			if !ok {
				log.Err(r, "got content-type MessagePack but handler is not a decoder", unexpectedContentType)
				writeApiReqErr(w, http.StatusBadRequest, errProcessingMsg)
				return
			}
			if err := msgp.Decode(r.Body, decoder); err != nil {
				log.Err(r, errDecodingMsgp, err)
				writeApiReqErr(w, http.StatusInternalServerError, errProcessingMsg)
				return
			}
		case util.ContentTypeJSON:
			if err := json.NewDecoder(r.Body).Decode(handler); err != nil {
				log.Err(r, errDecodingJSON, err)
				writeApiReqErr(w, http.StatusInternalServerError, errProcessingMsg)
				return
			}
		case util.ContentTypeFormURL:
			decoder, ok := handler.(formURLDecoder)
			if !ok {
				log.Err(r, fmt.Sprintf("got content-type form-urlencoded but handler of type %T is not a form decoder", handler),
					unexpectedContentType)
				writeBadAPIReq(w)
				return
			}
			if err := decoder.decodeFormURL(r); err != nil {
				log.Err(r, "could not decode form URL-encoded payload", err)
				writeApiReqErr(w, http.StatusInternalServerError, errProcessingMsg)
				return
			}
		default:
			// Somehow the Content-Type header is set to something unexpected. Assume the data is supposed to
			// be MessagePack encoded. If the content type is not MessagePack, the decoder will give an error.
			log.Err(r, fmt.Sprintf("got unexpected content type %q", contentType), badContentType)
			decoder, ok := handler.(msgp.Decoder)
			if !ok {
				log.Err(r, "got unexpected content-type but handler is not a decoder", unexpectedContentType)
				writeApiReqErr(w, http.StatusBadRequest, errProcessingMsg)
				return
			}
			if err := msgp.Decode(r.Body, decoder); err != nil {
				log.Err(r, errDecodingMsgp, err)
				writeApiReqErr(w, http.StatusInternalServerError, errProcessingMsg)
				return
			}
		}
	}

	if !handler.authorized(r, s, u) {
		if u.Id == 0 {
			resp := errMustLogin()
			writeApiReqErr(w, resp.Status, resp.err)
			return
		}
		resp := errLowPrivileges()
		writeApiReqErr(w, resp.Status, resp.err)
		return
	}

	resp := handler.handle(r, s, u)

	// Handle logging in and logging out cookies and redirection.
	if slugs[0] == "auth" {
		usingForm := false
		if len(slugs) > 2 && slugs[2] == "form" {
			usingForm = true
		}
		if r.Method == http.MethodPost {
			if resp.err == "" { // Ok to login (the resp.Warnings field is not used here).
				w.Header().Add("Set-Cookie", createLoginCookie(r, u.Id, u.Pass, s.Id))
				if usingForm {
					http.Redirect(w, r, loginRedirect(r, s, u), http.StatusSeeOther)
					return
				}
				// Below we write the error message to the user.
			} else if resp.err == errNoMembership {
				msg := "?msg=" + url.QueryEscape(errNoMembership)
				http.Redirect(w, r, registerRedirect(r, s, u)+msg, http.StatusSeeOther)
				return
			} else if usingForm {
				// Using a login form, the log in failed, so redirect the user to the page they came from.
				http.Redirect(w, r, r.URL.Path+"?msg="+url.QueryEscape(resp.err), http.StatusSeeOther)
				return
			}
			// Else, there is an error. Below we write the error message to the user.
		} else if r.Method == http.MethodDelete {
			w.Header().Add("Set-Cookie", createLogoutCookie())
			http.Redirect(w, r, logoutRedirect(r, s, u), http.StatusSeeOther)
			return
		}
	}

	if resp.err != "" {
		writeApiReqErr(w, resp.Status, resp.err)
		return
	}

	// The request succeeded, so write out the response in MessagePack format.
	w.Header().Set("Content-Type", util.ContentTypeMessagePack)
	if err := msgp.Encode(w, resp); err != nil {
		log.Err(r, "error encoding or writing API response", err)
	}
}

// An APIEndpoint endpoint represents a set of handlers of API requests at a particular path.
type APIEndpoint struct {
	// SubPathHandler, given the HTTP request method, returns the corresponding APIHandler.
	Points func(string) APIHandler

	// SubPaths contains APIEndpoint configurations for sub-paths of the endpoint.
	SubPaths map[string]APIEndpoint
}

// APIHandler represents a handler of an API endpoint. The concrete type of an APIHandler is able to decode
// the request body and be populated with the contents.
type APIHandler interface {
	// authorized says if the current user is authorized to make the request. This function is always called
	// after the request is decoded.
	authorized(*http.Request, *data.Site, *data.User) bool

	// handle handles the request and returns a response.
	handle(*http.Request, *data.Site, *data.User) *APIResponse
}

// apiEndpoints lists the handling functions to all the API endpoints.
// At the root level in this map and at each nested sub-map of APIEndpoint elements there must be
// a non-nil Points function defined.
// The SubPathHandler may be nil for any endpoint.
var apiEndpoints = map[string]APIEndpoint{
	"user": {
		Points: func(method string) APIHandler {
			switch method {
			case http.MethodGet: // TODO: this should only DEFAULT to getting the current user, but there should be an optional parameter to get another user
				return reqUserCheckCurrent{}
			case http.MethodPost:
				return new(ReqUserCreate)
			case http.MethodPatch:
				return new(ReqUserEdit)
			default:
				return nil
			}
		},
		SubPaths: map[string]APIEndpoint{
			"sites": {
				Points: func(method string) APIHandler {
					switch method {
					case http.MethodGet:
						return userSitesList{}
					default:
						return nil
					}
				},
			},
		},
	},
	"auth": {
		Points: func(method string) APIHandler {
			switch method {
			case http.MethodPost:
				return new(AuthLogin)
			case http.MethodDelete:
				return authLogout{}
			default:
				return nil
			}
		},
		SubPaths: map[string]APIEndpoint{
			"form": {
				Points: func(method string) APIHandler {
					switch method {
					case http.MethodPost:
						return new(authLoginForm)
					default:
						return nil
					}
				},
			},
		},
	},
	"pagepost": {
		Points: func(method string) APIHandler {
			switch method {
			case http.MethodGet:
				return new(ReqPagepostList)
			case http.MethodPost:
				return new(ReqPagepostCreate)
			case http.MethodPatch:
				return new(ReqPageEditPublish)
			case http.MethodDelete:
				return new(SiteDelete)
			default:
				return nil
			}
		},
		SubPaths: map[string]APIEndpoint{
			"element": {
				Points: func(method string) APIHandler {
					switch method {
					case http.MethodPost:
						return new(ReqPagepostMakeDynElem)
					default:
						return nil
					}
				},
			},
		},
	},
	"page-edit": {
		Points: func(method string) APIHandler {
			switch method {
			case http.MethodGet:
				return new(ReqPageEditContent)
			case http.MethodPost:
				return new(ReqPageEditSave)
			case http.MethodPatch:
				return new(ReqPageEditPublish)
			default:
				return nil
			}
		},
	},
	"site": {
		Points: func(method string) APIHandler {
			switch method {
			case http.MethodPost:
				return new(SiteCreate)
			case http.MethodDelete:
				return new(SiteDelete)
			default:
				return nil
			}
		},
		SubPaths: map[string]APIEndpoint{
			"theme": {
				Points: func(method string) APIHandler {
					switch method {
					case http.MethodGet:
						return new(SiteGetTheme)
					default:
						return nil
					}
				},
			},
			"option": {
				Points: func(method string) APIHandler {
					switch method {
					default:
						return nil
					}
				},
			},
		},
	},
	"admin": {
		Points: func(method string) APIHandler {
			switch method {
			case http.MethodGet:
				return reqAdminGetSrcs{}
			case http.MethodPost:
				return new(ReqAdminUpdateSrcs) // need method set of pointer
			default:
				return nil
			}
		},
	},
}

// A formURLDecoder is able to take the data of a Request and decode the application/x-www-form-urlencoded data
// into the data structure.
type formURLDecoder interface {
	decodeFormURL(*http.Request) error
}

var unexpectedContentType = errors.New("unexpected content type")

// APIResponse represents a response by a handler of an API endpoint. If an API endpoint handler returns an APIResponse with
// a non-empty err field, all the other fields will be ignored and the HTTP response status will not be 200.
type APIResponse struct {
	Status   int          // the HTTP status code (not encoded in response); this should be left at 0 to indicate status OK
	Body     msgp.Encoder // the body of the response
	Warnings []string     // warnings are added when non-critical parts of a request fail
	err      string       // an error is set when the request fails, in which case only the error message is sent in the response
}

// EncodeMsg implements msgp.Encodable. The body of the response must never be nil.
func (ar *APIResponse) EncodeMsg(en *msgp.Writer) error {
	// The map header indicates the number of elements (either 1 or 2).
	if len(ar.Warnings) == 0 {
		err := en.Append(0x81)
		if err != nil {
			return err
		}
	} else {
		err := en.Append(0x82)
		if err != nil {
			return err
		}
		err = en.Append(0xa4, 'w', 'a', 'r', 'n')
		if err != nil {
			return err
		}
		err = en.WriteArrayHeader(uint32(len(ar.Warnings)))
		if err != nil {
			return err
		}
		for i := range ar.Warnings {
			err = en.WriteString(ar.Warnings[i])
			if err != nil {
				return err
			}
		}
	}

	// The body must not be nil.
	err := en.Append(0xa4, 'b', 'o', 'd', 'y')
	if err != nil {
		return err
	}

	// Write the body.
	return ar.Body.EncodeMsg(en)
}

// warn adds a warning message to the API response, which indicates to the user that some but not all of the
// requested actions were completed successfully. An HTTP status code is not set.
func (ar *APIResponse) warn(m string) {
	ar.Warnings = append(ar.Warnings, m)
}

// APIResponseErr returns an APIResponse with the 500 error and the message given.
func APIResponseErr(msg string) *APIResponse {
	return &APIResponse{Status: http.StatusInternalServerError, err: msg}
}

// errProcessing indicates that an error occurred processing the request.
// The HTTP status code is 500.
func errProcessing() *APIResponse {
	return &APIResponse{Status: http.StatusInternalServerError, err: errProcessingMsg}
}

// errMustLogin indicates that the user must log in.
// The HTTP status code is 403.
func errMustLogin() *APIResponse {
	return &APIResponse{Status: http.StatusForbidden, err: "You must log in to continue."}
}

// errLowPrivileges indicates that the user does not have the privileges to do the request.
// The HTTP status code is 400.
func errLowPrivileges() *APIResponse {
	return &APIResponse{Status: http.StatusBadRequest, err: "Your privileges on this site are not high enough for this action."}
}

var badContentType = errors.New("bad content type")
