package main

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/dchenk/mazewire/pkg/data"
	"github.com/dchenk/mazewire/pkg/env"
	"github.com/dchenk/mazewire/pkg/log"
	"github.com/dchenk/mazewire/pkg/users"
	"github.com/dchenk/mazewire/pkg/util"
	"golang.org/x/crypto/bcrypt"
)

type AuthLogin struct {
	User string `msgp:"user"` // can be either username or email address
	Pass string `msgp:"pass"`
}

// authorized returns true for login requests.
func (*AuthLogin) authorized(*http.Request, *data.Site, *data.User) bool {
	return true
}

// authLogin allows a user to log in via an HTTP request not using a form.
// The response will have the err field either blank or providing an error message for the user.
func (al *AuthLogin) handle(r *http.Request, s *data.Site, u *data.User) *APIResponse {
	return checkLoginCreds(r, s, u, al.User, al.Pass)
}

// RespUserLogin is the message sent with OK responses for login requests.
var RespUserLogin = msgp.Raw([]byte{0xa2, 'o', 'k'})

// authLoginForm: POST auth
// This endpoint handles user log-ins via a form.
// Use the application/x-www-form-urlencoded content-type for the request.
// This type implements formURLDecoder; it is not exported because it doesn't need a
// MessagePack decoder.
type authLoginForm struct {
	User string // can be either username or email address
	Pass string
}

// authorized returns true for login requests by form.
func (*authLoginForm) authorized(*http.Request, *data.Site, *data.User) bool {
	return true
}

// handle allows a user to log in via a form.
// The APIResponse returned gets its Body set here because only the Err field is used to provide a (possibly blank) message.
// User u is updated by pointer if found.
// API params:
//  - user (username or email address)
//  - pass (password)
func (auf *authLoginForm) handle(r *http.Request, s *data.Site, u *data.User) *APIResponse {
	return checkLoginCreds(r, s, u, auf.User, auf.Pass)
}

func (auf *authLoginForm) decodeFormURL(r *http.Request) error {
	auf.User = r.PostFormValue("user")
	auf.Pass = r.PostFormValue("pass")
	return nil
}

// checkLoginCreds either sets the Err field of the returned APIResponse to an error message or leaves that field blank to indicate
// that the user should be logged in.
// The data.User passed in should be not nil in case a real user's data is retrieved and populated into the struct.
// This function always returns a non-nil *APIResponse.
func checkLoginCreds(r *http.Request, s *data.Site, u *data.User, unameEmail string, pass string) *APIResponse {
	if u.Id != 0 {
		return APIResponseErr("You are already logged in as user: " + u.Uname)
	}

	unameEmail = strings.TrimSpace(unameEmail)
	pass = strings.TrimSpace(pass)
	if util.IsAnyStringBlank(unameEmail, pass) {
		return APIResponseErr("You must provide a username and a password")
	}

	var (
		err   error
		tempU *data.User
	)
	if strings.ContainsRune(unameEmail, '@') {
		tempU, err = data.Conn.UserSiteInfoByEmail(s, unameEmail)
	} else {
		tempU, err = data.Conn.UserSiteInfoByUsername(s, unameEmail)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return APIResponseErr(errInvalidLogin)
		}
		log.Err(r, "error checking if email exists", err)
		return errProcessing()
	}

	if tempU.Role == users.RoleNone {
		return APIResponseErr(errNoMembership)
	}

	if err = bcrypt.CompareHashAndPassword(tempU.Pass, []byte(pass)); err != nil {
		if err != bcrypt.ErrMismatchedHashAndPassword {
			log.Err(r, "could not compare bcrypt hashes", err)
		}
		return APIResponseErr(errInvalidLogin)
	}

	// Copy tempU fields into u.
	*u = *tempU

	// Set the body whether or not err is blank. The body must be set to something, and the error will
	// be checked anyway.
	return &APIResponse{Body: RespUserLogin}
}

// createLoginCookie creates a complete cookie for a user to log in.
// The cookie's expiration is set to 4 hours. The hashedUserPass must be the current, validly hashed password
// belonging to the user.
// The first part of the cookie body, terminated by a dot, consists of five parts:
//  - MessagePack-encoded int64 giving the cookie expiration as a Unix timestamp (seconds)
//  - 32-byte hash of the user agent token
//  - MessagePack-encoded int64 giving the user ID
//  - MessagePack-encoded int64 giving the site ID (each cookie belongs to only one site)
//  - 21-byte fragment of the hashed password
func createLoginCookie(r *http.Request, userID int64, hashedUserPass []byte, siteID int64) string {
	body := make([]byte, 0, msgp.Int64Size*3+md5.Size+21)

	body = msgp.AppendInt64(body, time.Now().Add(time.Hour*4).Unix())

	body = append(body, userAgentToken(r.UserAgent())...)

	body = msgp.AppendInt64(body, userID)

	body = msgp.AppendInt64(body, siteID)

	// The password fragment does not including the bcrypt header identifying the
	// hashing algorithm used.
	body = append(body, hashedUserPass[10:31]...)

	bodyEncoded := make([]byte, base64.RawURLEncoding.EncodedLen(len(body)))
	base64.RawURLEncoding.Encode(bodyEncoded, body)

	var cookie strings.Builder
	cookie.Grow(cookieLenBase + len(bodyEncoded))

	cookie.WriteString("wwuser=")
	cookie.Write(bodyEncoded)
	cookie.WriteByte('.')

	// The signature is a signature of the base64-encoded body.
	sig := hmac.New(sha256.New, TOKEN_KEY)
	sig.Write(bodyEncoded)
	cookie.WriteString(base64.RawURLEncoding.EncodeToString(sig.Sum(nil)))

	if env.Prod() {
		cookie.WriteString("; Secure; HttpOnly; Path=/")
	} else {
		cookie.WriteString("; HttpOnly; Path=/")
	}

	return cookie.String()
}

// cookieLenBase is the length of a constant string included in the auth cookie body.
// The 43 represents base64.RawURLEncoding.EncodedLen(sha256.Size) from the cookie body signature.
const cookieLenBase = len("wwuser=") + 1 + len("; Secure; HttpOnly; Path=/") + 43

// authLogout: DELETE auth
type authLogout struct{}

// authorized returns true.
func (authLogout) authorized(*http.Request, *data.Site, *data.User) bool {
	return true
}

// handle logs the user out by resetting the user cookie.
// This function is here just for API handler map consistency and is handled in the main handler.
// This is the single APIHandler handle function that returns a nil *APIResponse.
func (authLogout) handle(_ *http.Request, _ *data.Site, _ *data.User) *APIResponse {
	return nil
}

func createLogoutCookie() string {
	if env.Prod() {
		return "wwuser=none; Secure; HttpOnly; Path=/"
	}
	return "wwuser=none; HttpOnly; Path=/"
}
