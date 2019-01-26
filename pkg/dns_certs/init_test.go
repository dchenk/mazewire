package dns_certs

import (
	"sort"
	"testing"

	"cloud.google.com/go/storage"
)

func TestRemoveBlank(t *testing.T) {

	lists := []struct {
		start, want []string
	}{
		{
			[]string{"abc"},
			[]string{"abc"},
		},
		{
			[]string{"abc", ""},
			[]string{"abc"},
		},
		{
			[]string{"abc", "", "def"},
			[]string{"abc", "def"},
		},
	}

	for i, list := range lists {

		got := removeBlank(list.start)
		for gotI := range got {
			if got[gotI] != list.want[gotI] {
				t.Errorf("(index %d) got bad slice: %v", i, got)
			}
		}

	}

}

func TestObjAttrsList_sort(t *testing.T) {

	oal := &objAttrsList{[]*storage.ObjectAttrs{{Name: "abc"}, {Name: "aaa"}, {Name: "z"}}}
	want := []*storage.ObjectAttrs{{Name: "aaa"}, {Name: "abc"}, {Name: "z"}}

	sort.Sort(oal)
	for i, o := range oal.oa {
		if o.Name != want[i].Name {
			t.Errorf("got bad sorting: %v", oal.oa)
		}
	}

}
