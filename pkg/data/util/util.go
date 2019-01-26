package util

import (
	"strconv"
	"strings"
)

// IdEq returns the string `id=ID_GIVEN`
func IdEq(id int64) string {
	return "id=" + strconv.FormatInt(id, 10)
}

// UserEq returns the string `user=ID_GIVEN`
func UserEq(id int64) string {
	return "user=" + strconv.FormatInt(id, 10)
}

// IntsList returns a list of the vals separated by a comma.
func IntsList(vals []int64) string {
	var list strings.Builder
	list.Grow(len(vals) * 6)
	for i := range vals {
		if i > 0 {
			list.WriteByte(',')
		}
		list.WriteString(strconv.FormatInt(vals[i], 10))
	}
	return list.String()
}

// SingleQuote quotes string s with single quotes and escapes single quotes within.
func SingleQuote(s string) string {
	return "'" + singleQuoteReplacer.Replace(s) + "'"
}

var singleQuoteReplacer = strings.NewReplacer("'", "''")

// JoinQuoted joins the list of vals single-quoted and separated by commas.
func JoinQuoted(vals []string) string {
	var values strings.Builder
	for i := range vals {
		if i > 0 {
			values.WriteByte(',')
		}
		values.WriteString(SingleQuote(vals[i]))
	}
	return values.String()
}
