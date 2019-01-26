package modules

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/dchenk/mazewire/pkg/room"
	"github.com/golang/protobuf/proto"
)

func mustProtoMarshal(v proto.Message) []byte {
	data, err := proto.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}

func TestDefaultModuleCustomizer_GetModuleCompiler(t *testing.T) {
	// mustCompile names all of the module types for which this package has a compiler.
	mustCompile := []string{room.TreeType, TypeHTML, TypeImage, TypeNav, TypeText}
	for i := range mustCompile {
		t.Run("case_"+strconv.Itoa(i), func(t *testing.T) {
		})
	}
}

func TestHTML_Compile(t *testing.T) {
	cases := []struct {
		room.Module
		Out []byte
	}{
		{
			Module: room.Module{Type: "html", Data: mustProtoMarshal(&HTML{
				Html: []byte("hello"),
			})},
			Out: []byte("<div>hello</div>"),
		},
		{
			Module: room.Module{Type: "html", Data: mustProtoMarshal(&HTML{
				Tag:  "h3",
				Html: []byte("hello"),
			})},
			Out: []byte("<h3>hello</h3>"),
		},
		{
			Module: room.Module{Type: "html", Data: mustProtoMarshal(&HTML{
				Tag:    "span",
				IdAttr: "thing",
				Html:   []byte("<h1>hi</h1><span class=\"a\">text</span>"),
			})},
			Out: []byte("<span id=\"thing\"><h1>hi</h1><span class=\"a\">text</span></span>"),
		},
	}
	for i := range cases {
		t.Run("case_"+strconv.Itoa(i), func(t *testing.T) {
			m, err := DefaultCustomizer.GetModuleCompiler(&cases[i].Module).Compile(nil, nil, &cases[i].Module, nil)
			if err != nil {
				t.Fatalf("got an error compiling; %v", err)
			}
			if m.Type != room.StaticElementType {
				t.Errorf("got module type %q", m.Type)
			}
			if !bytes.Equal(m.Data, cases[i].Out) {
				t.Errorf("got bad output: %q", m.Data)
			}
		})
	}
}

func TestNavLinks_Compile(t *testing.T) {

	cases := []struct {
		nav  Nav
		want string
	}{
		{
			nav: Nav{
				Common: &room.Common{},
				Links:  []*NavLink{{Text: "The link text", Href: "href1", Target: ""}}},
			want: `<ul><li><a href="href1">The link text</a></li></ul>`,
		},
		{
			nav: Nav{
				Common: &room.Common{},
				Links: []*NavLink{{Text: "The link text", Href: "href1", Target: "_blank"},
					{Text: "the 2nd text", Href: "href2", Target: ""},
					{Text: "the 3rd text", Href: "href3", Target: "blank"},
				}},
			want: `<ul><li><a href="href1" target="_blank">The link text</a></li><li><a href="href2">the 2nd text</a></li><li><a href="href3" target="blank">the 3rd text</a></li></ul>`,
		},
	}

	var b bytes.Buffer
	for i := range cases {
		t.Run("case_"+strconv.Itoa(i), func(t *testing.T) {
			b.Reset()
			err := cases[i].nav.BuildHTML(nil, nil, nil, &b, nil)
			if err != nil {
				t.Fatal("got an error;", err)
			}
			got := b.String()
			if got != cases[i].want {
				t.Errorf("got bad nav link list: %s at index %d", got, i)
			}
		})
	}

}
