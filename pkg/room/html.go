package room

import (
	"bytes"
	"strconv"
	"strings"
)

const (
	// ColumnOpen is the tag that opens a column.
	ColumnOpen = `<div class="room-col">`

	// ColumnClose is the tag that closes a column.
	ColumnClose = "</div>"
)

// ElementOpenTag writes to buf an opening HTML tag.
func ElementOpenTag(tagName string, idAttr string, classes []string, buf *bytes.Buffer) {
	buf.WriteByte('<')
	buf.WriteString(tagName)
	WriteIdAttr(idAttr, buf)
	WriteClassAttr(classes, buf)
	buf.WriteByte('>')
}

// WriteIdAttr checks if id is not an empty string, and if it's not empty this function writes to the
// buffer the string ` id="id-given"` where "id-given" is the id argument.
func WriteIdAttr(id string, b *bytes.Buffer) {
	id = strings.TrimSpace(id)
	if id != "" {
		b.WriteString(" id=")
		b.WriteString(strconv.Quote(id))
	}
}

// WriteClassAttr checks if classes is not empty, and if it's not empty this function writes to the
// buffer the string ` class="classes given"` where "classes given" is the classes argument's values
// joined by a space.
func WriteClassAttr(classes []string, b *bytes.Buffer) {
	if len(classes) > 0 {
		b.WriteString(" class=")
		b.WriteString(strconv.Quote(strings.Join(classes, " ")))
	}
}

// makeClassesList returns the list of classes the element should have in the open tag.
// A class of the format elemType+"-"+elemDynId is included if elemDynId is not zero.
func makeClassesList(elemType string, elemSubType string, elemDynId int64, userClasses []string) []string {
	list := make([]string, 2, 3+len(userClasses))
	list[0] = "room-" + elemType
	list[1] = elemType + "-" + elemSubType
	if elemDynId != 0 {
		list = append(list, elemType+"-"+strconv.FormatInt(elemDynId, 10))
	}
	list = append(list, userClasses...)
	return list
}
