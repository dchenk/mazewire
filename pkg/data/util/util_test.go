package util

import "testing"

func TestIdEq(t *testing.T) {

	pairs := map[int64]string{
		82:    "id=82",
		0:     "id=0",
		95899: "id=95899",
		-234:  "id=-234",
	}

	for k, v := range pairs {
		if str := IdEq(k); str != v {
			t.Errorf("bad int64 id stringifying with %d", k)
		}
	}

}

func TestUserEq(t *testing.T) {

	pairs := map[int64]string{
		82: "user=82",
		0:  "user=0",
	}

	for k, v := range pairs {
		if str := UserEq(k); str != v {
			t.Errorf("bad int64 id stringifying with %d", k)
		}
	}

}

func TestSingleQuote(t *testing.T) {

	cases := []struct {
		before, after string
	}{
		{"abc", "'abc'"},
		{"a'c", "'a''c'"},
		{"'", "''''"},
		{"a''b", "'a''''b'"},
		{"\t", "'\t'"},
		{"\n", "'\n'"},
		{"\n", `'
'`},
		{`a	bc
d`, "'a\tbc\nd'"},
	}

	for i, pair := range cases {
		got := SingleQuote(pair.before)
		if got != pair.after {
			t.Errorf("(index %d) did not escape string correctly; got: %s", i, got)
		}
	}

}

func TestJoinSingleQuoted(t *testing.T) {
	cases := []struct {
		vals   []string
		output string // the different ways to do this, given that map elements are not ordered
	}{
		{vals: []string{"a", "b", "a\tb"}, output: "'a','b','a\tb'"},
		{vals: []string{}, output: ""},
		{vals: []string{"abc"}, output: "'abc'"},
		{vals: []string{"''a", "'b'", "a'"}, output: "'''''a','''b''','a'''"},
	}
	for i, pair := range cases {
		got := JoinQuoted(pair.vals)
		if got != pair.output {
			t.Errorf("(index %d) did not construct list correctly; got %v", i, got)
		}
	}
}
