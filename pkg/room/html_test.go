package room

import (
	"bytes"
	"strconv"
	"testing"
)

func TestElementOpenTag(t *testing.T) {

	tcs := []struct {
		tagName, idAttr string
		classes         []string
		result          string
	}{
		{"div", "", []string{}, `<div>`},
		{"div", "stuff", []string{"my-class"}, `<div id="stuff" class="my-class">`},
		{"h2", "", []string{"c1", "c2", "c3"}, `<h2 class="c1 c2 c3">`},
		{"section", "a-b-c", []string{"my-class"}, `<section id="a-b-c" class="my-class">`},
	}

	var b bytes.Buffer
	for i, tc := range tcs {
		t.Run("case_"+strconv.Itoa(i), func(t *testing.T) {
			b.Reset()
			ElementOpenTag(tc.tagName, tc.idAttr, tc.classes, &b)
			got := b.String()
			if got != tc.result {
				t.Errorf("got wrong output %q", got)
			}
		})
	}

}

func TestWriteIdAttr(t *testing.T) {
	cases := []struct {
		id, output string
	}{
		{"", ""},
		{" ", ""},
		{"\t\n \t", ""},
		{"a", " id=\"a\""},
		{"ab-c", " id=\"ab-c\""},
	}
	var b bytes.Buffer
	for i := range cases {
		t.Run("case_"+strconv.Itoa(i), func(t *testing.T) {
			b.Reset()
			WriteIdAttr(cases[i].id, &b)
			got := b.String()
			if got != cases[i].output {
				t.Errorf("got wrong output %q", got)
			}
		})
	}
}

func TestWriteClassAttr(t *testing.T) {
	cases := []struct {
		classes []string
		output  string
	}{
		{[]string{}, ""},
		{[]string{"a"}, " class=\"a\""},
		{[]string{"ab", "c"}, " class=\"ab c\""},
	}
	var b bytes.Buffer
	for i := range cases {
		t.Run("case_"+strconv.Itoa(i), func(t *testing.T) {
			b.Reset()
			WriteClassAttr(cases[i].classes, &b)
			got := b.String()
			if got != cases[i].output {
				t.Errorf("got wrong output %q", got)
			}
		})
	}
}

func TestMakeClassesList(t *testing.T) {

	tcs := []struct {
		elemType, elemSubType string
		elemId                int64
		userClasses           []string
		result                []string
	}{
		{elemType: "section", elemSubType: "standard", elemId: 52, userClasses: []string{"cool", "classes"},
			result: []string{"room-section", "section-standard", "section-52", "cool", "classes"}},
		{elemType: "row", elemSubType: "halfhalf", elemId: 3928349, userClasses: []string{"a_b-c"},
			result: []string{"room-row", "row-halfhalf", "row-3928349", "a_b-c"}},
	}

	for i, tc := range tcs {
		t.Run("case_"+strconv.Itoa(i), func(t *testing.T) {
			got := makeClassesList(tc.elemType, tc.elemSubType, tc.elemId, tc.userClasses)
			if len(got) != len(tc.result) {
				t.Fatalf("got wrong output length: %v", got)
			}
			for j := range got {
				if got[j] != tc.result[j] {
					t.Fatalf("got wrong output: %v", got)
				}
			}
		})
	}

}
