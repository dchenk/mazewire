package modules

// This file contains all of the built-in ModuleMaker implementations.

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/dchenk/go-render-quill"
	"github.com/dchenk/mazewire/pkg/room"
	"github.com/golang/protobuf/proto"
)

const (
	TypeNav   = "nav"
	TypeHTML  = "html"
	TypeText  = "text"
	TypeImage = "image"
)

// defaultModuleMakerSource is the built-in ModuleSource. It only looks at the given Module's type to try
// to find a suitable Module implementation.
//
// We use an empty struct here because this package does not attempt to customize the output of any module
// on a per-request basis. How could it? It's up to your application to create a ModuleSource that will be
// able to take whatever data you give it so that you can compile and build custom modules.
type defaultModuleCustomizer struct{}

// GetModuleCompiler only looks at the Module's Type field to identify if a default ModuleCustomizer exists.
func (ds defaultModuleCustomizer) GetModuleCompiler(m *room.Module) room.ModuleCompiler {
	switch m.Type {
	case room.TreeType:
		return new(room.Tree)
	//case "blurb":
	//	cm.BodyHTML = []byte("<div class=blurb></div>")
	//}
	case TypeHTML:
		return new(HTML)
	case TypeNav:
		return new(Nav)
	case TypeText:
		return new(Text)
	case TypeImage:
		return new(Image)
	default:
		return nil
	}
}

func (ds defaultModuleCustomizer) GetModuleBuilder(compiledModule *room.Module) room.ModuleBuilder {
	switch compiledModule.Type {
	case TypeNav:
		return new(Nav)
	case TypeText:
		return nil // Dynamic Text not yet implemented.
		// return new(Text)
	case TypeImage:
		return new(Image)
	default:
		return nil
	}
}

// DefaultCustomizer is the default room.ModuleCustomizer.
var DefaultCustomizer defaultModuleCustomizer

// Compile implements ModuleCompiler for HTMLModule. This function always returns a static element.
// If the Tag field is empty, the HTML is rendered inside a <div> element.
func (h *HTML) Compile(_ room.DataStore, _ room.ModuleCustomizer, m *room.Module, _ *room.PageCSS) (*room.Module, error) {
	if err := proto.Unmarshal(m.Data, h); err != nil {
		return nil, err
	}
	tag := strings.TrimSpace(h.Tag)
	if tag == "" {
		tag = "div"
	}
	openTag := "<" + tag
	if h.IdAttr != "" {
		openTag += " id=" + strconv.Quote(h.IdAttr)
	}
	openTag += ">"
	closeTag := "</" + tag + ">"
	l1 := len(openTag)
	l2 := l1 + len(h.Html)
	data := make([]byte, l2+len(closeTag))
	copy(data, openTag)
	copy(data[l1:], h.Html)
	copy(data[l2:], closeTag)
	return &room.Module{Type: room.StaticElementType, Data: data}, nil
}

func (i *Image) Compile(_ room.DataStore, _ room.ModuleCustomizer, m *room.Module, css *room.PageCSS) (*room.Module, error) {
	if err := proto.Unmarshal(m.Data, i); err != nil {
		return nil, err
	}
	var b bytes.Buffer
	if i.LinkUrl != "" {
		b.WriteString("<a href=")
		b.WriteString(strconv.Quote(i.LinkUrl))
		b.WriteByte('>')
	}
	b.WriteString("<img src=")
	b.WriteString(strconv.Quote(i.Src))
	if i.Alt != "" {
		b.WriteString(" alt=")
		b.WriteString(strconv.Quote(i.Alt))
	}
	room.WriteIdAttr(i.Common.IdAttr, &b)
	room.WriteClassAttr(i.Common.Classes, &b)
	b.WriteByte('>')
	if i.LinkUrl != "" {
		b.WriteString("</a>")
	}
	return &room.Module{Type: TypeImage, Dyn: m.Dyn, Data: b.Bytes()}, nil
}

// BuildHTML implements room.ModuleBuilder for Image. NOTE: For now, the Nav is compiled to its final static form and BuildHTML
// simply outputs the pre-built HTML. We're doing it this way to prepare for making Nav dynamic.
func (i *Image) BuildHTML(_ room.DataStore, _ room.ModuleCustomizer, m *room.Module, b *bytes.Buffer, _ *room.PageCSS) error {
	b.Write(m.Data) // TODO: when we're ready for it, we'll be doing proto.Unmarshall into the Image.
	return nil
}

// Compile implements ModuleCompiler for Nav.
func (n *Nav) Compile(ds room.DataStore, mc room.ModuleCustomizer, m *room.Module, css *room.PageCSS) (*room.Module, error) {
	if err := proto.Unmarshal(m.Data, n); err != nil {
		return nil, err
	}
	// TODO: distinguish between static and dynamic modules
	var b bytes.Buffer
	err := n.BuildHTML(ds, nil, nil, &b, css)
	if err != nil {
		return nil, err
	}
	return &room.Module{Type: TypeNav, Dyn: m.Dyn, Data: b.Bytes()}, nil
}

// BuildHTML implements ModuleBuilder. (For now this exists because we will add some dynamic functionality.)
func (n *Nav) BuildHTML(_ room.DataStore, _ room.ModuleCustomizer, _ *room.Module, b *bytes.Buffer, _ *room.PageCSS) error {
	room.ElementOpenTag("nav", n.Common.IdAttr, n.Common.Classes, b)
	for _, nl := range n.Links {
		b.WriteString("<li><a href=")
		b.WriteString(strconv.Quote(nl.Href))
		if nl.Target != "" {
			b.WriteString(" target=")
			b.WriteString(strconv.Quote(nl.Target))
		}
		b.WriteByte('>')
		b.WriteString(nl.Text)
		b.WriteString("</a></li>")
	}
	b.WriteString("</ul>")
	return nil
}

// Compile implements ModuleCompiler for TextModule.
func (t *Text) Compile(_ room.DataStore, _ room.ModuleCustomizer, m *room.Module, css *room.PageCSS) (*room.Module, error) {
	err := proto.Unmarshal(m.Data, t)
	if err != nil {
		return nil, err
	}
	cm := &room.Module{Type: TypeText, Dyn: m.Dyn}
	cm.Data, err = quill.Render(t.Ops)
	return cm, err
}
