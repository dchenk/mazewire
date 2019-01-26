package cockroach

import "testing"

func TestSetValuesList(t *testing.T) {

	cases := []struct {
		vals        map[string]interface{}
		listOptions []string // the different ways to do this, given that map elements are not ordered
	}{
		{
			vals:        map[string]interface{}{},
			listOptions: []string{""},
		},
		{
			vals:        map[string]interface{}{"a": 1},
			listOptions: []string{"`a`=$1"},
		},
		{
			vals:        map[string]interface{}{"a": 1, "b": 1},
			listOptions: []string{"`a`=$1,`b`=$2", "`b`=$1,`a`=$2"},
		},
	}

	for i, pair := range cases {
		list, args := setValuesList(pair.vals)
		if len(args) != len(pair.vals) {
			t.Errorf("(index %d) did not construct arguments list correctly; got %v", i, args)
		}
		goodCase := false
		for _, opt := range pair.listOptions {
			if list == opt {
				goodCase = true
				break
			}
		}
		if !goodCase {
			t.Errorf("(index %d) did not construct list correctly; got: %s", i, list)
		}
	}

}
