package util

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestMustGetEnv(t *testing.T) {
	// The PATH variable should exist. If the variable is not found, the function panics.
	if len(MustGetEnv("PATH")) <= 0 {
		t.Fail()
	}
}

func TestIsAnyStringBlank(t *testing.T) {
	blank := [][]string{{""}, {" ", "stuff"}, {"	"}, {" 	"}, {`
 `, "stuff"}}
	notBlank := [][]string{{"stuff", "- "}, {"stuff"}}
	for i := range blank {
		if !IsAnyStringBlank(blank[i]...) {
			t.Errorf("(index %d): bad classification of blank string", i)
		}
	}
	for i := range notBlank {
		if IsAnyStringBlank(notBlank[i]...) {
			t.Errorf("(index %d): bad classification of notBlank string", i)
		}
	}
}

func TestValidEmail(t *testing.T) {

	good := []string{
		`email@example.com`,
		`firstname.lastname@example.com`,
		`firstname+lastname@example.com`,
		`firstname-lastname@example.com`,
		`email@subdomain.example.com`,
		`email@subdomain.sub.example.com`,
		`email@subdomain.sub.example.co.uk`,
		`em.ail@subdomain.sub.example.co.uk`,
		`email@12afas.fff`,
		`1234567890@example.com`,
		`email@example-one.com`,
		`_______@example.com`,
		`email@example.name`,
		`email@example.museum`,
		`email@example.co.jp`,
		`ema{il@example.co.jp`,
		`ema{}il@example.co.jp`,
		`much.a'faf@example.com`,
	}

	bad := []string{
		``,
		` a`,
		`a@b.c`,
		`a@b.co`,
		`adsf@asdf`,
		`sdfa.asdf`,
		`asdf.asdf.asdfas`,
		`.adsf@ags.fdf`,
		`adsf@.asd.fsf`,
		`very.unusual@unusual.com@example.com`,
		`actually@invalid.-com`,
		`actually@inv\alid.-com`,
		`actually@inv\*lid.-com`,
		`actually@inv\rlid.-com`,
		`actually@inv//\/rlid.-com`,
		`act//\/ually@invrlid.-com`,
		`actually@inv\1lid.-com`,
		`actually@inv\_lid.-com`,
		`actually@inval id.-com`,
		`act ually@invalid.-com`,
		`asdf"sdf@asdf.asdf`,
	}

	for _, a := range good {
		if !ValidEmail(a) {
			t.Errorf("wrong classification of good %q email", a)
		}
	}

	for _, a := range bad {
		if ValidEmail(a) {
			t.Errorf("wrong classification of bad %q email", a)
		}
	}

}

func TestValidUsername(t *testing.T) {

	good := []string{
		"asd",
		"g2m",
		"k__",
		"y$$",
		"s$_",
		"u_$",
		"h5$",
		"h_2",
		"alsdlfalsdf",
		"sdf324",
		"g__234dsf$$$",
	}

	bad := []string{
		"",
		"a",
		"ag",
		"2adsf",
		"_asdfadsf",
		"$asdfadsf",
		"____",
		"adsf&asdfs",
		"asdf*sdf",
		"hasdf#sdfs",
		"toolooooonnnnggglooooonggglooonnngggloooooonnnngggloooonnggg", // 60 characters
	}

	for i := range good {
		if !ValidUsername(good[i]) {
			t.Errorf("Good username %q marked as bad\n" + good[i])
		}
	}
	for i := range bad {
		if ValidUsername(bad[i]) {
			t.Errorf("Bad username %q marked as good\n" + bad[i])
		}
	}

}

func TestValidPathSlug(t *testing.T) {

	good := []string{
		"a-v",
		"asdf",
		"_sdf",
		"1g",
		"f-sdf-a.adf",
		"j$$.$$u",
		"54+~a",
		"8*6",
		"r~i",
		"a__$~+g",
		"12.adsf",
		"a++++____---~~~~****.....w",
		"__",
		"SLKDJFLKSDJFLSD", // This is valid, but upon save all capital letters should be set to lowercase.
	}

	bad := []string{
		"",  // cannot be blank
		" ", // space only
		"	", // tab character only
		"g",          // too short
		"_",          // too short
		"-",          // has no word character
		"--",         // starts and ends without a word character
		"-asdf",      // starts without a word character
		"ag-",        // ends without a word character
		"ag–",        // ends without a word character (n-dash)
		"asldf.",     // ends without a word character
		"*sf",        // starts without a word character
		"ag–ava",     // contains n-dash
		"ah^adf",     // contains ^
		"aa&aa",      // contains &
		"aa?aa",      // contains ?
		"aa@aa",      // contains @
		"sdf%5df",    // contains %
		"sdf#5df",    // contains #
		"asdf`dfsdf", // contains `
		"as'sldf",    // contains '
		"longer-than-80-characters-adfasdfasfaaasasdfasfasdfafgasdfasdfasdfaasdfasdfasdfas",
	}

	for _, s := range good {
		if !ValidPathSlug(s) {
			t.Errorf("Good slug %q reported as bad\n", s)
		}
	}

	for _, s := range bad {
		if ValidPathSlug(s) {
			t.Errorf("Bad slug %q reported as good\n", s)
		}
	}

}

func TestSanitizePhone(t *testing.T) {

	nonNilErr := errors.New("")

	cases := []struct {
		in, out string
		err     error
	}{
		{"", "", nonNilErr},
		{"          ", "", nonNilErr},
		{"123", "123", nonNilErr},
		{"122233344445", "122233344445", nonNilErr},
		{"22233344445", "22233344445", nonNilErr},
		{"5032223333", "5032223333", nil},
		{"15032222222", "5032222222", nil},
		{"+1 503 222 3333", "5032223333", nil},
		{"1(503) 222-4444", "5032224444", nil},
		{"(503) 444-5555", "5034445555", nil},
		{"a 222 333 4444", "2223334444", nil},
	}

	for i, tc := range cases {

		got, err := SanitizePhone(tc.in)
		if err != nil && tc.err == nil {
			t.Errorf("(index %d) did not expect error; %v", i, err)
		} else if err == nil && tc.err != nil {
			t.Errorf("(index %d) expected but did not get error", i)
		}
		if got != tc.out {
			t.Errorf("(index %d) bad output; got %q", i, got)
		}

	}

}

func BenchmarkSanitizePhone(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SanitizePhone("123")
		SanitizePhone("5032223333")
		SanitizePhone("22233344445")
		SanitizePhone("+1 503 222 3333")
		SanitizePhone("(503) 444-5555")
	}
}

func SanitizePhone2(s string) (string, error) {

	out := make([]byte, 0, 11)
	for i := 0; i < len(s); i++ {
		if bytes.IndexByte(numbers, s[i]) > -1 {
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

var numbers = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

func BenchmarkSanitizePhone2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SanitizePhone2("123")
		SanitizePhone2("5032223333")
		SanitizePhone2("22233344445")
		SanitizePhone2("+1 503 222 3333")
		SanitizePhone2("(503) 444-5555")
	}
}

var reqPaths = map[string][]string{
	"/":                          {"/"},
	"//":                         {"/"},
	"/aaa//bbb":                  {"aaa"},
	"/aaa":                       {"aaa"},
	"/aaa/":                      {"aaa"},
	"/aaa/bbb":                   {"aaa", "bbb"},
	"/aaa/bbb/":                  {"aaa", "bbb"},
	"/aaa/bbb/ccc":               {"aaa", "bbb", "ccc"},
	"/aaa/bbb/ccc/":              {"aaa", "bbb", "ccc"},
	"/aaa/bbb/ccc/ddd":           {"aaa", "bbb", "ccc"},
	"/aaaa/bbbbb/ccccc/dddd/eee": {"aaaa", "bbbbb", "ccccc"},
}

func TestSplitRequestPath(t *testing.T) {

	stringSlicesEqual := func(a, b []string) bool {
		if len(a) != len(b) {
			return false
		}
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}

	funcs := []func(string) []string{SplitRequestPath, splitRequestPathOld, splitRequestPath2}

	for i, f := range funcs {
		for path, slice := range reqPaths {
			got := f(path)
			if !stringSlicesEqual(slice, got) {
				t.Errorf("func %d: Parsed request path bad for path %q; got %v", i, path, got)
			}
		}
	}

}

func splitRequestPathOld(p string) []string {
	if p == "/" {
		return homeSlice
	}

	split := strings.Split(p[1:], "/") // Ignore the leading slash.

	l := len(split)

	// We only care about the first three parts.
	if l > 3 {
		split = split[:3]
	}

	if l > 2 && split[2] == "" {
		split = split[:2]
	}

	if l > 1 && split[1] == "" {
		split = split[:1]
	}

	if len(split) == 1 && split[0] == "" {
		split = homeSlice
	}

	return split
}

func splitRequestPath2(p string) []string {

	split := strings.Split(p[1:], "/") // Ignore the leading slash.

	// We only care about the first three parts.
	if len(split) > 3 {
		split = split[:3]
	}

	if len(split) > 2 && split[2] == "" {
		split = split[:2]
	}

	if len(split) > 1 && split[1] == "" {
		split = split[:1]
	}

	if split[0] == "" {
		return homeSlice
	}

	return split

}

var Split []string // used in benchmarks

func BenchmarkSplitRequestPath(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for p := range reqPaths {
			Split = SplitRequestPath(p)
		}
	}
}

func BenchmarkSplitRequestPathOld(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for p := range reqPaths {
			Split = splitRequestPathOld(p)
		}
	}
}

func BenchmarkSplitRequestPath2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for p := range reqPaths {
			Split = splitRequestPath2(p)
		}
	}
}

// TestTimestampNow simply makes sure that TimestampNow does not return the zero time.
func TestTimestampNow(t *testing.T) {
	tnp, err := ParseTimestamp(TimestampNow())
	if err != nil {
		t.Errorf("error creating timestamp; %s", err)
		t.FailNow()
	}
	if tnp.Equal(time.Time{}) {
		t.Errorf("got zero for TimestampNow")
	}
}

func TestParseTimestamp(t *testing.T) {

	vals := map[string]time.Time{
		"2018-09-11 11:08:55": time.Date(2018, 9, 11, 11, 8, 55, 0, time.UTC),
		"2029-10-03 23:00:53": time.Date(2029, 10, 03, 23, 0, 53, 0, time.UTC),
	}

	for k, v := range vals {
		tp, err := ParseTimestamp(k)
		if err != nil {
			t.Errorf("could not parse time: %s", err)
			t.FailNow()
		}
		if !tp.Equal(v) {
			t.Errorf("parsed time incorrect; got %s for %q", tp, k)
		}
	}

}

func TestValidTimestamp(t *testing.T) {
	good := []string{"1000-01-01 00:00:00", "9999-12-31 23:59:59", "2019-01-01 00:24:00"}
	bad := []string{"", "1", "abc", "1000-01-01", "00:00:00", "9999-12-31 23::", "2019-01-01 5"}
	for _, ts := range good {
		if !ValidTimestamp(ts) {
			t.Errorf("Failed test of goodTimestamp with good timestamp %q\n", ts)
		}
	}
	for _, ts := range bad {
		if ValidTimestamp(ts) {
			t.Errorf("Failed test of goodTimestamp with bad timestamp %q\n", ts)
		}
	}
}

func TestExtractDomain(t *testing.T) {

	type strErr struct {
		domain string
		err    error
	}

	invalid := strErr{"", ErrNoDomain}
	goodCom := strErr{"good.com", nil}
	subGoodCom := strErr{"sub.good.com", nil}

	cases := map[string]strErr{
		"":          invalid,
		"    ":      invalid,
		"http://":   invalid,
		"http:///":  invalid,
		"https://":  invalid,
		"https:///": invalid,
		"bad":       invalid,
		"good.com":  goodCom,
		"	good.com ": goodCom,
		"http://good.com/":    goodCom,
		"http://good.com":     goodCom,
		"https://good.com/":   goodCom,
		"http://good.com/   ": goodCom,
		"Good.com":            goodCom,
		"sub.good.com":        subGoodCom,
		"sub.good.com/":       subGoodCom,
		"Sub.GOOD.com/":       subGoodCom,
		"http://sub.good.com": subGoodCom,
	}

	for k, v := range cases {
		gotD, gotErr := ExtractDomain(k)
		if gotD != v.domain {
			t.Errorf("extracted (non)domain %q wrong; got %q", k, gotD)
		}
		if v.err != nil && gotErr == nil {
			t.Errorf("returned nil error but had to return error %q for domain %q", v.err.Error(), k)
		}
		if v.err == nil && gotErr != nil {
			t.Errorf("returned error %q but had to return nil error for domain %q", gotErr.Error(), k)
		}
	}

}

func TestSliceContainsString(t *testing.T) {
	has := map[string][]string{
		"a": {"a", "b"},
		"e": {"c", "e", "b"},
		"":  {"stuff", ""},
	}
	hasNot := map[string][]string{
		"y":  {"a", "b"},
		"i":  {"c", "e", "b"},
		"hl": {},
	}
	for str, slice := range has {
		if !SliceContainsString(slice, str) {
			t.Errorf("saying slice does not have string %q", str)
		}
	}
	for str, slice := range hasNot {
		if SliceContainsString(slice, str) {
			t.Errorf("saying slice has string %q", str)
		}
	}
}
