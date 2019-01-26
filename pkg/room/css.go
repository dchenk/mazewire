package room

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/dchenk/mazewire/pkg/room/special"
	"github.com/golang/protobuf/proto"
)

// PageCSS embeds the bytes.Buffer type and contains the state of the CSS on the page, which is a record of the modules
// that have written their CSS to the page. Before a module writes anything to a PageCSSState, it should first check if
// another module of the same type in the Tree has already written the same CSS.
//
// After the compilation of a Tree, if PageCSS turns out to be not empty then it should be stored and passed in to the
// Tree upon building.
type PageCSS struct {
	bytes.Buffer

	// WrittenModules contains a list of the module types that have already written to the page CSS.
	// When a Module writes its CSS to the page, it should say that it did so already so that if the same type of module
	// appears on the page again it doesn't try to write duplicate CSS.
	writtenModules []string
}

// ProtoMessage returns a proto.Message that can be used to
// The concrete type of the message is special.PageCSS.
func (css *PageCSS) ProtoMessage() proto.Message {
	return &special.PageCSS{
		Code:           css.Buffer.Bytes(),
		WrittenModules: css.writtenModules,
	}
}

// UnmarshalProto reads the Protocol Buffers encoded message into the internal data structure.
// The struct is first reset.
func (css *PageCSS) UnmarshalProto(msg []byte) error {
	css.Reset()
	var m special.PageCSS
	err := proto.Unmarshal(msg, &m)
	if err != nil {
		return err
	}
	css.Write(m.Code)
	css.writtenModules = append(css.writtenModules, m.WrittenModules...)
	return nil
}

func (css *PageCSS) Reset() {
	css.Buffer.Reset()
	css.writtenModules = css.writtenModules[:0]
}

// RecordWriting records that a module with the given type just now wrote its CSS to the page.
func (css *PageCSS) RecordWriting(moduleType string) {
	css.writtenModules = append(css.writtenModules, moduleType)
}

// IsWritten says if a module with the given type has already written its CSS to the page.
func (css *PageCSS) IsWritten(moduleType string) bool {
	for _, t := range css.writtenModules {
		if t == moduleType {
			return true
		}
	}
	return false
}

// CompileStyles writes CSS code out to the css buffer and gives a list of additional classes the element
// being processed should be given. The styles map has CSS property keys and values, though most keys are
// shortened, not standard CSS property names.
func CompileStyles(styles map[string]string, css *PageCSS) (classes []string) {
	for prop, val := range styles {
		switch prop {
		case "bkg_color":
			if styles["bkg_type"] == "color" {
				css.WriteString("background-color:")
				css.WriteString(val)
				css.WriteByte(';')
			}
		case "bkg_img":
			val = strings.TrimSpace(val)
			if val != "" {
				css.WriteString("background-image:url('")
				css.WriteString(val)
				css.WriteString("') no-repeat;")
			}
		case "bkg_size":
			if strings.TrimSpace(styles["bkg_img"]) != "" {
				css.WriteString("background-size:")
				css.WriteString(val)
				css.WriteByte(';')
			}
		case "color":
			css.WriteString("color:")
			css.WriteString(val)
			css.WriteByte(';')
		case "font_fam":
			ss := strings.Split(strings.TrimSpace(val), ",")
			css.WriteString("font-family:")
			for i, f := range ss {
				f = strings.TrimSpace(f)
				if len(f) > 0 {
					if i > 0 {
						css.WriteByte(',')
					}
					if f[len(f)-1] == '"' {
						css.WriteString(f)
					} else {
						css.WriteString(strconv.Quote(f))
					}
				}
			}
			css.WriteByte(';')
		case "pad_bottom", "pad_top", "paragraph_margin_bottom", "row_col_space", "section_row_space":
			classes = append(classes, prop+"-"+val)
		}
		// All other properties are either non-styling options or are used only as indicators for what
		// to do with the other properties.
	}
	return
}
