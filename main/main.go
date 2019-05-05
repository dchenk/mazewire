package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"database/sql"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-plugin"
	"golang.org/x/oauth2/google"

	"github.com/dchenk/mazewire/pkg/data"
	"github.com/dchenk/mazewire/pkg/default_cert"
	"github.com/dchenk/mazewire/pkg/env"
	"github.com/dchenk/mazewire/pkg/log"
	"github.com/dchenk/mazewire/pkg/util"
)

var (
	// userContentSrc  = "https://storage.googleapis.com/" + env.ContentBucket + "/"
	gcpDefaultCreds *google.Credentials
)

func main() {
	log.Info(nil, "=== STARTED SERVER ===")

	var err error
	gcpDefaultCreds, err = google.FindDefaultCredentials(
		context.Background(), "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		log.Critical(nil, "could not get GCP credentials", err)
		return
	}

	if err = data.Init(); err != nil {
		log.Critical(nil, "could not initialize DB connection", err)
		return
	}
	defer data.Conn.Close()

	defer plugin.CleanupClients()

	http.HandleFunc("/", handler)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	addrHTTPS, addrHTTP := ":https", ":http"
	if !env.Prod() {
		addrHTTPS, addrHTTP = ":10443", ":8080"
	}

	go func() {
		serverHTTP := &http.Server{
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
			IdleTimeout:  time.Minute,
		}
		httpLn, err := net.Listen("tcp", addrHTTP)
		if err != nil {
			log.Critical(nil, "could not start HTTP listener", err)
			os.Exit(1)
		}
		err = serverHTTP.Serve(httpLn)
		if err != http.ErrServerClosed {
			log.Critical(nil, "http server stopped", err)
		}
		os.Exit(1)
	}()

	serverHTTPS := &http.Server{
		TLSConfig: &tls.Config{
			GetCertificate: func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
				// TODO: return a certificate based on the ServerName
				return cert.DefaultCert()
			},
		},
		ReadTimeout:  time.Second * 12,
		WriteTimeout: time.Second * 12,
		IdleTimeout:  time.Minute * 2,
	}

	httpsLn, err := net.Listen("tcp", addrHTTPS)
	if err != nil {
		log.Critical(nil, "could not start HTTPS listener", err)
		return
	}

	if env.Prod() {
		go serveHealthCheck()
	}

	err = serverHTTPS.ServeTLS(httpsLn, "", "")
	if err != http.ErrServerClosed {
		log.Critical(nil, "https server stopped", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {

	userChan := make(chan *userSite, 1)
	go getCurrentUser(r, userChan)

	s, err := data.Conn.SiteByDomain(strings.ToLower(r.Host))
	if err != nil {
		if err != sql.ErrNoRows {
			log.Err(r, "error looking up host "+r.Host, err)
			http.Error(w, "An error occurred.", http.StatusNotFound)
			return
		}
		http.NotFound(w, r)
		return
	}

	if s.Tls == 2 && r.TLS == nil {
		http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
		return
	}

	slugs := util.SplitRequestPath(r.URL.Path)

	userSite := <-userChan // Block until the user is retrieved.

	// At this point, if someone is logged in then we we got a user from the cookie but we need to verify
	// that the cookie was made for the current host.
	if userSite.u.Id != 0 {
		if userSite.siteID != s.Id {
			// This person may be trying to hack in to the site using a cookie from another site.
			log.Err(r, "bad cookie", errors.New("bad site ID in cookie"))
			// On the front end, this error message may help identify that this is where the error occurred.
			http.Error(w, "An error occurred.", http.StatusNotFound)
			return
		}
	}

	if slugs[0] == "api" {
		handleAPI(w, r, s, userSite.u, slugs[1:])
		return
	}

	w.Header().Set("Content-Type", util.ContentTypeHTML)

	// Respond with a site page or the site admin page, always with the status 200 (even if showing a 404).
	code, page := handlePage(r, slugs, s, userSite.u)
	w.WriteHeader(code)
	// TODO: allow customizing the lang on a per-page basis
	if err := writeDocHTML(w, page, "en"); err != nil {
		// log it
	}
}

// userSite conveniently wraps a pointer to a User and a site ID and is used as the value sent down the channel
// by getCurrentUser.
type userSite struct {
	u      *data.User
	siteID int64
}

// writeApiReqErr writes out an API request error message. The HTTP status set is always not 200. If status code
// given is 0 or 200, the status written is 500.
func writeApiReqErr(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", util.ContentTypeTextPlainUTF8)
	if status > 0 && status != 200 {
		w.WriteHeader(status)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write([]byte(message))
}

func writeBadAPIReq(w http.ResponseWriter) {
	writeApiReqErr(w, http.StatusBadRequest, "Bad Request")
}

// writeDocHTML writes a response to an HTTP request given the HTML head and body buffers.
// The lang argument specifies the "lang" attribute set on the opening <html> tag.
func writeDocHTML(w http.ResponseWriter, content *contentBuffers, lang string) error {
	_, err := w.Write(htmlOpen1)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(lang))
	if err != nil {
		return err
	}

	_, err = w.Write(htmlOpen2)
	if err != nil {
		return err
	}

	_, err = w.Write(content.head.Bytes())
	if err != nil {
		return err
	}

	_, err = w.Write(htmlHeadBody)
	if err != nil {
		return err
	}

	_, err = w.Write(content.body.Bytes())
	if err != nil {
		return err
	}

	_, err = w.Write(htmlClose)
	return err
}

var (
	htmlOpen1    = []byte(`<!DOCTYPE html><html lang="`)
	htmlOpen2    = []byte(`"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1">`)
	htmlHeadBody = []byte("</head><body>")
	htmlClose    = []byte("</body></html>")
)

type contentBuffers struct {
	head, body bytes.Buffer
}

// The main site has ID = 1.
func mainSite() *data.Site {
	return &data.Site{Id: 1}
}

// serveHealthCheck listens on port 85 for the Google Cloud load balancer health check requests and responds.
// The health check monitor should send a body of the string "a" and expects a response of the string "a".
// This function terminates the program with code 1 if it fails to begin listening.
func serveHealthCheck() {
	ln, err := net.ListenTCP("tcp", &net.TCPAddr{Port: 85})
	if err != nil {
		log.Critical(nil, "could not begin health check listener; %v", err)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Err(nil, "could not accept health check request", err)
			continue
		}
		// The load balancer sends a particular string and expects it echoed back.
		if _, err := io.Copy(conn, conn); err != nil {
			log.Err(nil, "could not write health check response", err)
		}
		conn.Close()
	}
}
