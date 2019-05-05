package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dchenk/mazewire/pkg/data"
	"github.com/dchenk/mazewire/pkg/log"
	"github.com/dchenk/mazewire/pkg/util"
)

// getCurrentUser sends a non-nil *data.User to chan c. If a real logged in user is behind the
// request, then the User struct  will be filled in with the basic info of the user and their role
// for the current site.
func getCurrentUser(r *http.Request, c chan *userSite) {

	us := &userSite{u: new(data.User)}
	defer func() { // Make sure the created User struct is sent when the function exits.
		c <- us
	}()

	cookie, err := r.Cookie("mwuser")
	if err == http.ErrNoCookie || len(cookie.Value) < 40 {
		// The base64-encoded cookie value includes a hashed signature, so in all the whole thing
		// must be at least 40 bytes long.
		return
	}

	splitCookie := strings.Split(cookie.Value, ".")

	if len(splitCookie) != 2 {
		log.Err(r, userCookieErr(cookie.Value), errors.New("cookie does not have two parts"))
		return
	}

	// Hash the cookie body to verify that it was not tampered with
	sig := hmac.New(sha256.New, TOKEN_KEY)
	sig.Write([]byte(splitCookie[0]))

	// Compare the cookie's hashed body with the signature that came with it.
	if base64.RawURLEncoding.EncodeToString(sig.Sum(nil)) != splitCookie[1] {
		log.Err(r, userCookieErr(cookie.Value), errors.New("cookie signature is not correct"))
		return
	}

	body, err := base64.RawURLEncoding.DecodeString(splitCookie[0])
	if err != nil {
		log.Err(r, userCookieErr(cookie.Value), err)
		return
	}

	var expiration int64
	expiration, body, err = msgp.ReadInt64Bytes(body)
	if err != nil {
		log.Err(r, userCookieErr(cookie.Value), err)
		return
	}

	if time.Now().Unix() > expiration {
		return
	}

	if len(body) < 32 || !bytes.Equal(body[:md5.Size], userAgentToken(r.UserAgent())) {
		return
	}

	var userID int64
	userID, body, err = msgp.ReadInt64Bytes(body[md5.Size:])
	if err != nil {
		log.Err(r, userCookieErr(cookie.Value), err)
		return
	}

	us.siteID, body, err = msgp.ReadInt64Bytes(body)
	if err != nil {
		log.Err(r, userCookieErr(cookie.Value), err)
		return
	}

	// Get the user's info, including their password.
	userGot, err := data.Conn.UserSiteInfoByID(us.siteID, userID)
	if err != nil {
		log.Err(r, "error getting current user basic info", err)
		return
	}

	if !bytes.Equal(userGot.Pass[10:31], body) {
		return
	}

	// At this point, we know that everything with the cookie is good, and the user is authenticated.
	us.u = userGot
}

func userCookieErr(cookieVal string) string {
	return fmt.Sprintf("user cookie error; cookie value: %q", cookieVal)
}

// ReqUserCreate: POST user
// Users can only create other users on a website if they're currently on that website.
type ReqUserCreate struct {
	// Username.
	Uname string

	// Email address.
	Email string

	// The user's password in plain text.
	Pass string

	// First name.
	Fname string

	// Last name.
	Lname string

	// The role the user is going to have on the site.
	Role string
}

// handle creates a user on the site s. This API endpoint accepts the Protocol Buffers, JSON, and
// x-www-form-urlencoded content types.
// TODO: If a user exists, assign to site.
func (req *ReqUserCreate) handle(r *http.Request, s *data.Site, u *data.User) *APIResponse {
	// All values must be provided.
	if util.IsAnyStringBlank(req.Uname, req.Fname, req.Lname, req.Email, req.Pass, req.Role) {
		return APIResponseErr(errIncompleteForm)
	}

	if !util.ValidEmail(req.Email) {
		return APIResponseErr("The email you entered seems to be invalid.")
	}

	// Set all capital letters in the chosen username to lowercase.
	req.Uname = strings.ToLower(req.Uname)

	if !util.ValidUsername(req.Uname) {
		return APIResponseErr("The username you selected is not valid. It must begin with a letter and be at least 3 characters long; the special characters allowed are '$' and '_'.")
	}

	// If the user is registering not on the main site, check if registrations are open for the site.
	if s.Id != 1 {
		// TODO: What site is the user being added to? What's the role? Is it allowed?
	}

	if !users.RoleIsValid(req.Role) {
		log.Err(r, fmt.Sprintf("non standard user role desired %q", req.Role), errors.New("bad role requested"))
		return APIResponseErr("An error occurred. Something is wrong with the role you're supposed to be given on the site.")
	}

	// TODO: who is creating the new user? Someone already logged in? How to handle that?
	if u.Id != 0 {
		// the user creating this account needs to have at least the role that they're giving the new user
		if !users.RoleAtLeast(u.Role, req.Role) {
			return APIResponseErr("Your role in this site doesn't allow you to create a user with this role")
		}
	}

	// TODO: Validate the password as sufficiently strong.
	if len(req.Pass) < 6 {
		return APIResponseErr("The password you chose is too short.")
	}

	passHash, err := hashPassword(req.Pass)
	if err != nil {
		log.Err(r, "error hashing user password", err)
		return errProcessing()
	}

	// Check if the email is already taken. Although we are not doing this in a transaction and someone could
	// take the username within the next couple microseconds until the call to create the user, the database contains
	// a uniqueness constraint on the username column, so the worst thing that can happen is this request will fail.
	unameCount, emailCount, err := data.Conn.UserCountByUnameEmail(req.Uname, req.Email)
	if err != nil {
		log.Err(r, "error checking if user exists with Uname: "+req.Uname, err)
		return errProcessing()
	}

	if emailCount > 0 {
		// TODO: offer a link to sign in with the username retrieved... "Try signing in with your Mazewire username <b>"+existsUsername+"</b>"
		// TODO: <br><a href="/existing-accounts" target="_blank">LEARN MORE</a>
		e := "The email address you entered is already taken by a user. Perhaps you have already signed up for an account " +
			"on a website built with Mazewire and need to log in?"
		return APIResponseErr(e)
	}

	if unameCount > 0 {
		e := "The username you chose is already taken by a user. Perhaps you have already signed up for an account " +
			"on a website built with Mazewire and need to log in?"
		return APIResponseErr(e)
	}

	// TODO: do this in a transaction
	newID, err := data.Conn.UserInsert(req.Uname, req.Email, passHash, u.Fname, u.Lname)
	if err != nil || newID == 0 { // make sure no rows were inserted
		log.Err(r, "could not create user", err)
		return errProcessing()
	}

	resp := APIResponse{
		Body: &RespUserCreate{req.Uname, newID},
	}

	// add the user to the site they registered under --- TODO: do this all in a transaction
	// if err := newUser.SetSiteRole(s.Id, reqData.Role); err != nil {
	// 	log.Err(r, "error setting new user site role", err)
	// 	rt := setNewUserSiteRoleRetrier{
	// 		SiteID: s.Id,
	// 		UserID: newID,
	// 		Role:   reqData.Role,
	// 	}
	// 	// Retry the request right away, and if that fails save it for later.
	// 	if err := rt.Retry(); err != nil {
	// 		//if err := saveRetry(r, rt); err != nil {
	// 		//	log.Err(r, "could not save retry", err)
	// 		//	resp.warn("Part of the process to create your website failed. Please contact support for assistance.")
	// 		//} else {
	// 		//	resp.warn("Part of the process to create your website failed. Please try again.")
	// 		//}
	// 		// Save a message (as an option) saying that the blobs table did not get created.
	// 		go saveUserMessage(r, newID, rt.Key(), "Part of the process to create your account did not complete.")
	// 	}
	// 	return nil // TODO -----
	// }

	// TODO: on client side, after receiving OK response, log the user in immediately and redirect to account page only if no user is currently logged in
	return &resp
}

// authorized returns true for user creation requests.
func (*ReqUserCreate) authorized(*http.Request, *data.Site, *data.User) bool {
	return true
}

// decodeFormURL implements formURLDecoder for the UserCreate endpoint.
func (uc *ReqUserCreate) decodeFormURL(r *http.Request) error {
	uc.Uname = r.FormValue("uname")
	uc.Email = r.FormValue("email")
	uc.Fname = r.FormValue("fname")
	uc.Lname = r.FormValue("uname")
	uc.Role = r.FormValue("role")
	uc.Pass = r.FormValue("pass")
	return nil
}

type RespUserCreate struct {
	Username string `msgp:"uname" json:"uname"`
	Id       int64  `msgp:"id" json:"id"`
}

// reqUserCheckCurrent: GET: user
type reqUserCheckCurrent struct{}

func (reqUserCheckCurrent) authorized(*http.Request, *data.Site, *data.User) bool {
	return true
}

// userCheckCurrent responds with the username of the currently logged in user. If nobody is logged in, the response body is "Nobody".
func (reqUserCheckCurrent) handle(_ *http.Request, _ *data.Site, u *data.User) *APIResponse {
	var uname string
	if u.Id == 0 {
		uname = "Nobody"
	} else {
		uname = u.Uname
	}
	return &APIResponse{Body: msgp.Raw(msgp.AppendString(nil, uname))}
}

// ReqUserEdit: PATCH user
type ReqUserEdit struct {
	// TODO
	UserID int64 `msgp:"user_id"` // the user to edit
}

func (ue *ReqUserEdit) authorized(_ *http.Request, _ *data.Site, u *data.User) bool {
	// TODO: or, authorize if current user is at least admin and the user to edit is a subscriber, in which case the current
	// user can only edit options of the user.
	return u.Id == ue.UserID
}

// userEdit handles all changes to a user, including changes to a user's meta data.
func (*ReqUserEdit) handle(r *http.Request, _ *data.Site, _ *data.User) *APIResponse {
	// TODO
	return APIResponseErr("this API endpoint is not ready")
}

// delete a user and all of the usermeta associated with the user
// API params:
// - user (username)
// - pass (password)
func ReqUserDelete(r *http.Request, _ *data.Site, _ *data.User) *APIResponse {
	// TODO
	return APIResponseErr("this API endpoint is not ready")
}

// userSitesList: GET user/sites
// Retrieve the sites with which the user is associated
type userSitesList struct{}

func (userSitesList) authorized(_ *http.Request, _ *data.Site, u *data.User) bool {
	return u.Id != 0
}

func (userSitesList) handle(r *http.Request, _ *data.Site, u *data.User) *APIResponse {
	metas, err := data.Conn.UserMetaByIdLikeKey(u.Id, "role%")
	if err != nil {
		log.Err(r, "could not get a user's associated sites", err)
		return errProcessing()
	}
	// Extract the site IDs.
	ids := make([]int64, 0, len(metas))
	for i := range metas {
		k := metas[i].K
		id, err := strconv.ParseInt(k[len("role"):], 10, 64)
		if err != nil {
			log.Err(r, fmt.Sprintf("could not parse site ID from the role key %q", k), err)
			return errProcessing()
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		// This should never happen because the user is logged in.
		log.Err(r, "got zero ids of associated sites", errors.New("unexpected empty list of associated sites"))
		return errProcessing()
	}
	// Retrieve info for the sites.
	sites, err := data.Conn.SitesByIDs(ids)
	if err != nil {
		log.Err(r, "could not get sites by IDs", err)
		return errProcessing()
	}
	// Convert the sites to UserSitesListItem items.
	list := make(RespUserSitesList, len(sites))
	for i := range sites {
		list[i].ID = sites[i].Id
		list[i].Domain = sites[i].Domain
		list[i].Name = sites[i].Name
		list[i].Logo = sites[i].Logo
		for _, m := range metas {
			if m.K == users.SiteRoleKey(sites[i].Id) {
				list[i].Role = string(metas[i].V)
				break
			}
		}
	}
	return &APIResponse{Body: list}
}

type RespUserSitesList []UserSitesListItem

type UserSitesListItem struct {
	ID      int64
	Domain  string
	Name    string
	Logo    string
	Favicon string
	Role    string
}

// userAgentToken returns the (salted) user agent token to use in a cookie.
// The returned slice is a 32-byte long hash made out of the concatenation
// of the user agent of the client and the global TOKEN_SALT.
func userAgentToken(userAgent string) []byte {
	h := md5.Sum([]byte(userAgent + TOKEN_SALT))
	return h[:]
}

// hashPassword uses bcrypt to hash the password string and possibly return an error.
func hashPassword(p string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(p), 14)
}
