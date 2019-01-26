package env

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
)

func TestCheckForRequiredVars(t *testing.T) {
	cases := []struct {
		vars    map[string]string
		missing []string
	}{
		{
			vars:    nil,
			missing: requiredVars,
		},
		{
			vars:    map[string]string{"ABC": "DEF", "DB_NAME": "value", "TOKEN_KEY": "value2"},
			missing: []string{VarDbConnection},
		},
		{
			vars: func() map[string]string {
				m := make(map[string]string, len(standardVars))
				for _, stdVar := range standardVars {
					m[stdVar] = "\t \n"
				}
				return m
			}(),
			missing: requiredVars,
		},
	}
	for i, tc := range cases {
		t.Run("case_"+strconv.Itoa(i), func(t *testing.T) {
			got := checkForRequiredVars(tc.vars)
			if len(got) != len(tc.missing) {
				t.Logf("got missing list of length %d but expected length %d", len(got), len(tc.missing))
				t.Fatalf("\tgot: %v; expected: %v", got, tc.missing)
			}
			for _, g := range got {
				found := false
				for _, ms := range tc.missing {
					if g == ms {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("marked %v incorrectly as missing", g)
				}
			}
		})
	}
}

func TestParseFile(t *testing.T) {
	cases := []struct {
		data string
		vars map[string]string
	}{
		{
			data: `ABC_DEF=GHI
DB_NAME=something_here
AB=cd ef gh
# Comment.

ANOTHER="the value"`,
			vars: map[string]string{
				"ABC_DEF":     "GHI",
				"DB_NAME":     "something_here",
				"AB":          "cd ef gh",
				"ANOTHER":     "\"the value\"",
				"PLUGINS_DIR": "plugins", // Default value.
			},
		},
		{
			data: "DB_CONNECTION=conn_value\nDB_NAME=database_name\nPLUGINS_DIR=system/plugins",
			vars: map[string]string{
				"DB_CONNECTION": "conn_value",
				"DB_NAME":       "database_name",
				"PLUGINS_DIR":   "system/plugins",
			},
		},
		{data: "", vars: map[string]string{"PLUGINS_DIR": "plugins"}},
		{data: " \n", vars: map[string]string{"PLUGINS_DIR": "plugins"}},
	}
	for i, tc := range cases {
		t.Run("case_"+strconv.Itoa(i), func(t *testing.T) {
			err := testParseFileTemp(tc.data, tc.vars)
			if err != nil {
				t.Error(err.Error())
			}
		})
	}
}

func testParseFileTemp(data string, vars map[string]string) error {
	temp, err := ioutil.TempFile("", "")
	if err != nil {
		return fmt.Errorf("could not create temp file; %v", err)
	}
	name := temp.Name()
	defer os.Remove(name)
	if _, err := temp.WriteString(data); err != nil {
		temp.Close()
		return fmt.Errorf("could not write vars data to temp file; %v", err)
	}
	if err = temp.Close(); err != nil {
		return fmt.Errorf("could not close newly written file; %v", err)
	}
	got, err := parseFile(name)
	if err != nil {
		return fmt.Errorf("could not parse file; %v", err)
	}
	return mapsEqual(vars, got)
}

func mapsEqual(m1, m2 map[string]string) error {
	if len(m1) != len(m2) {
		return fmt.Errorf("maps have unequal length: %d and %d;\nfirst: %v;\nsecond: %v",
			len(m1), len(m2), m1, m2)
	}
	for k, v := range m1 {
		v2 := m2[k]
		if v2 != v {
			return fmt.Errorf("first map has value %q but second has %q for key %q", v, v2, k)
		}
	}
	return nil
}
