package log

import (
	"context"
	"fmt"
	stdlog "log"
	"net"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/logging"
	"github.com/dchenk/mazewire/pkg/env"
)

var logger *logging.Logger

func init() {
	// TODO: this somehow must work with logging services other than StackDriver. Or should it?
	logClient, err := logging.NewClient(context.Background(), "projects/"+env.GCP_PROJECT)
	if err != nil {
		stdlog.Fatal("Could not create logging client; ", err)
	}
	defer logClient.Close()

	if env.Prod() {
		if err = logClient.Ping(context.Background()); err != nil {
			stdlog.Fatal("Could not ping logging service; ", err)
		}
	}

	ip := outboundIP()
	if ip == "" {
		ip = "t" + strconv.FormatInt(time.Now().Unix(), 10)
	}

	logger = logClient.Logger("server_" + ip)
}

// Info logs a message at the info severity level.
// The message is printed out to the console when in development mode.
func Info(r *http.Request, msg string) {
	if !env.Prod() {
		fmt.Println("INFO: ", msg)
		return
	}
	entry := logging.Entry{Severity: logging.Info, Payload: msg}
	if r != nil {
		entry.HTTPRequest = &logging.HTTPRequest{Request: r, RemoteIP: r.RemoteAddr}
	}
	logger.Log(entry)
}

// Err logs an error at the error severity level.
// The message with the error is printed out to the console when in development mode.
func Err(r *http.Request, msg string, e error) {
	msg = msg + "; " + e.Error()
	if !env.Prod() {
		stdlog.Println("ERR: ", msg)
		return
	}
	entry := logging.Entry{Severity: logging.Error, Payload: msg}
	if r != nil {
		entry.HTTPRequest = &logging.HTTPRequest{Request: r, RemoteIP: r.RemoteAddr}
	}
	logger.Log(entry)
}

// Critical logs an error at the critical severity level.
// The message with the error is printed out to the console when in development mode.
// This function flushes the logging buffer before returning.
func Critical(r *http.Request, msg string, e error) {
	msg = msg + "; " + e.Error()
	if !env.Prod() {
		stdlog.Println("CRITICAL: ", msg)
		return
	}
	entry := logging.Entry{Severity: logging.Critical, Payload: msg}
	if r != nil {
		entry.HTTPRequest = &logging.HTTPRequest{Request: r, RemoteIP: r.RemoteAddr}
	}
	logger.Log(entry)
	if err := logger.Flush(); err != nil {
		stdlog.Printf("log: could not flush the log after writing a critical error; %v", err)
	}
}

// outboundIP returns the IP address of the machine, or an empty string if the address could not be retrieved.
func outboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()
	ip := make([]rune, 0, 16)
	for _, char := range conn.LocalAddr().String() {
		if char == ':' {
			break
		}
		if char == '.' {
			ip = append(ip, '_')
		} else {
			ip = append(ip, char)
		}
	}
	return string(ip)
}
