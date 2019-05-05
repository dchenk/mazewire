package util

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	regEmail     = regexp.MustCompile(`\A[\w]+[\w+.!#$%'&{}|~-]+@[\w]+[\w+.-]+\.[a-zA-Z]{2,}\z`)
	regPathSlug  = regexp.MustCompile(`\A\w[\w+.$~*-]{0,78}\w\z`)
	regUsername  = regexp.MustCompile(`\A[a-z0-9][a-z0-9_]{2,48}\z`)
	regDomain    = regexp.MustCompile(`\A[a-z0-9][-a-z0-9.]{1,61}[a-z]\z`)
	regTimestamp = regexp.MustCompile(`\A[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}\z`)
)

const (
	ContentTypeJSON          = "application/json"
	ContentTypeHTML          = "text/html; charset=UTF-8"
	ContentTypeProtobuf      = "application/protobuf"
	ContentTypeTextPlain     = "text/plain"
	ContentTypeTextPlainUTF8 = "text/plain; charset=UTF-8"
	ContentTypeFormURL       = "application/x-www-form-urlencoded"
)

// MustGetEnv gets the environment variable or panics if it's not set.
func MustGetEnv(k string) string {
	v, ok := os.LookupEnv(k)
	if !ok {
		log.Panicf("environment variable %q not set", k)
	}
	return v
}

// IsAnyStringBlank says if any of the given strings is blank when trimmed.
func IsAnyStringBlank(strs ...string) bool {
	for _, s := range strs {
		if strings.TrimSpace(s) == "" {
			return true
		}
	}
	return false
}

// ValidEmail says if the given string has the valid form of an email address.
func ValidEmail(s string) bool {
	return len(s) <= 254 && regEmail.MatchString(s)
}

// ValidUsername says if the given string has the valid form of a username.
func ValidUsername(s string) bool {
	return regUsername.MatchString(s)
}

// ValidPathSlug says if string s conforms to the requirements of the format of a clean path slug.
// A slug is between 2 and 80 characters long, begins and ends with a word character, and has only word
// characters or one of +.$~*- in the middle.
func ValidPathSlug(s string) bool {
	return regPathSlug.MatchString(s)
}

// SanitizePhone removes all non-digits from s to format it roughly according to the North American
// Numbering Plan (NANP). If there is a leading 1 and more than 10 digits total, the 1 is removed.
func SanitizePhone(s string) (string, error) {
	out := make([]byte, 0, 11)
	for i := 0; i < len(s); i++ {
		if strings.IndexByte("0123456789", s[i]) > -1 {
			out = append(out, s[i])
		}
	}

	// Remove the leading 1 in the number if it's not in the area code.
	if len(out) == 11 && out[0] == '1' {
		out = out[1:]
	}

	outStr := string(out)

	if len(outStr) != 10 {
		return outStr, fmt.Errorf("util: string %q has not exactly 10 digits after cleanup", outStr)
	}

	return outStr, nil
}

// SplitRequestPath splits p by "/" and returns a slice of 1<=len<=3 representing the slugs of a request path.
func SplitRequestPath(p string) []string {
	p = p[1:]

	a := make([]string, 0, 3)

	for i := 0; i < 3; i++ {
		m := strings.IndexByte(p, '/')
		if m < 0 {
			m = len(p)
		}
		ts := p[:m]
		if ts == "" {
			break
		}
		a = append(a, ts)
		if m == len(p) {
			break
		}
		p = p[m+1:]
	}

	if len(a) == 0 {
		return []string{"/"}
	}

	return a
}

// TimeStampFormat is the format of a timestamp.
const TimeStampFormat = "2006-01-02 15:04:05"

// TimestampNow gives the current timestamp in the TimeStampFormat format.
func TimestampNow() string {
	return time.Now().Format(TimeStampFormat)
}

// ValidTimestamp says whether s is a valid timestamp.
func ValidTimestamp(s string) bool {
	return regTimestamp.MatchString(s)
}

// ParseTimestamp parses the timestamp in string t using the TimeStampFormat format.
func ParseTimestamp(t string) (time.Time, error) {
	return time.Parse(TimeStampFormat, t)
}

// ExtractDomain extracts a domain name (without a trailing dot) from str.
func ExtractDomain(str string) (string, error) {
	str = strings.ToLower(strings.TrimSpace(str))
	str = strings.TrimPrefix(str, "https://")
	str = strings.TrimPrefix(str, "http://")
	str = strings.Split(str, "/")[0]
	// There must be a dot in the domain name, and the string should pass the simple regex test.
	if strings.IndexByte(str, '.') == -1 || !regDomain.MatchString(str) {
		return "", ErrNoDomain
	}
	return str, nil
}

// ErrNoDomain indicates that a valid domain name was not given.
var ErrNoDomain = errors.New("util: no valid domain was provided")

// SliceContainsString says whether string s is present in slice a.
func SliceContainsString(a []string, s string) bool {
	for i := range a {
		if a[i] == s {
			return true
		}
	}
	return false
}
