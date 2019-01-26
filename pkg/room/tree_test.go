package room

import (
	"bytes"
	"strconv"
	"testing"
)

func TestTree_mergeStaticSections(t *testing.T) {

	cases := []struct {
		t1, t2 Tree
	}{
		{
			t1: Tree{Sections: []*Section{
				{Common: &Common{Type: StaticElementType}, Html: []byte("ab")},
				{Common: &Common{Type: StaticElementType}, Html: []byte("cd")},
			}},
			t2: Tree{Sections: []*Section{
				{Common: &Common{Type: StaticElementType}, Html: []byte("abcd")},
			}},
		},
		{
			t1: Tree{Sections: []*Section{
				{Common: &Common{Type: StaticElementType}, Html: []byte("ab")},
				{Common: &Common{Type: StaticElementType}, Html: []byte("cd")},
				{Common: &Common{Type: "something"}, Html: []byte("kuv")},
			}},
			t2: Tree{Sections: []*Section{
				{Common: &Common{Type: StaticElementType}, Html: []byte("abcd")},
				{Common: &Common{Type: "something"}, Html: []byte("kuv")},
			}},
		},
		{
			t1: Tree{Sections: []*Section{
				{Common: &Common{Type: "something"}, Html: []byte("ab")},
				{Common: &Common{Type: "something"}, Html: []byte("cd")},
			}},
			t2: Tree{Sections: []*Section{
				{Common: &Common{Type: "something"}, Html: []byte("ab")},
				{Common: &Common{Type: "something"}, Html: []byte("cd")},
			}},
		},
		{
			t1: Tree{Sections: []*Section{
				{Common: &Common{Type: "something"}, Html: []byte("ji")},
				{Common: &Common{Type: StaticElementType}, Html: []byte("ab")},
				{Common: &Common{Type: "something"}, Html: []byte("kuv")},
				{Common: &Common{Type: StaticElementType}, Html: []byte("cd")},
			}},
			t2: Tree{Sections: []*Section{
				{Common: &Common{Type: "something"}, Html: []byte("ji")},
				{Common: &Common{Type: StaticElementType}, Html: []byte("ab")},
				{Common: &Common{Type: "something"}, Html: []byte("kuv")},
				{Common: &Common{Type: StaticElementType}, Html: []byte("cd")},
			}},
		},
		{
			t1: Tree{Sections: []*Section{
				{Common: &Common{Type: "something"}, Html: []byte("ji")},
				{Common: &Common{Type: StaticElementType}, Html: []byte("ab")},
				{Common: &Common{Type: "something"}, Html: []byte("kuv")},
				{Common: &Common{Type: StaticElementType}, Html: []byte("cd")},
				{Common: &Common{Type: StaticElementType}, Html: []byte{'e'}},
				{Common: &Common{Type: StaticElementType}, Html: []byte{'y'}},
			}},
			t2: Tree{Sections: []*Section{
				{Common: &Common{Type: "something"}, Html: []byte("ji")},
				{Common: &Common{Type: StaticElementType}, Html: []byte("ab")},
				{Common: &Common{Type: "something"}, Html: []byte("kuv")},
				{Common: &Common{Type: StaticElementType}, Html: []byte("cdey")},
			}},
		},
	}

	for i := range cases {
		t.Run("case_"+strconv.Itoa(i), func(t *testing.T) {
			t1 := &cases[i].t1
			t2 := &cases[i].t2
			t1.mergeStaticSections()
			if len(t1.Sections) != len(t2.Sections) {
				t.Errorf("unequal number of sections; %d and %d", len(t1.Sections), len(t2.Sections))
			}
			for ii := range t1.Sections {
				if !bytes.Equal(t1.Sections[ii].Html, t2.Sections[ii].Html) {
					t.Errorf("unequal bytes of sections at %d; got %q", i, t1.Sections[ii].Html)
				}
			}
		})
	}

}

func TestTree_DynamicDataIDs(t *testing.T) {
	cases := []struct {
		t   Tree
		IDs []int64
	}{
		{sampleCompiledTree1, []int64{1}},
	}
	for i := range cases {
		t.Run("case_"+strconv.Itoa(i), func(t *testing.T) {
			tc := &cases[i]
			got := tc.t.DynamicDataIDs()
			if len(got) != len(tc.IDs) {
				t.Fatalf("unequal lengths; needed %d IDs but got %d", len(tc.IDs), len(got))
			}
		})
	}
}
