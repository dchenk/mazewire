package room

import (
	"bytes"

	"github.com/golang/protobuf/proto"
)

// buildHTML writes to buf the entire compiled and built section.
func (s *Section) buildHTML(ds DataStore, mc ModuleCustomizer, buf *bytes.Buffer, css *PageCSS) error {

	if s.Common.Type == StaticElementType {
		buf.Write(s.Html)
		return nil
	}

	if s.Dyn != 0 {

		dynData, err := ds.Get(s.Dyn)
		if err != nil {
			return err
		}

		var dynSettings Section
		err = proto.Unmarshal(dynData.Data, &dynSettings)
		if err != nil {
			return err
		}

		s.buildClassesAndStyles(&dynSettings, css)

	}

	s.openTag(buf)

	for i := range s.Rows {
		if err := s.Rows[i].buildHTML(ds, mc, buf, css); err != nil {
			return err
		}
	}

	s.closeTag(buf)

	return nil

}

func (s *Section) buildClassesAndStyles(settings *Section, css *PageCSS) {

	// Section elements have the "room-section" class.
	s.Common.Classes = makeClassesList("section", settings.Common.Type, s.Dyn, settings.Common.Classes)

	if len(settings.Common.Options) > 0 {
		css.WriteByte('.')
		css.WriteString(s.Common.Classes[2])
		css.WriteByte('{')
		moreClasses := CompileStyles(settings.Common.Options, css)
		css.WriteByte('}')

		s.Common.Classes = append(s.Common.Classes, moreClasses...)
	}

}

func (s *Section) openTag(buf *bytes.Buffer) {
	ElementOpenTag(SectionTagName, s.Common.IdAttr, s.Common.Classes, buf)
}

func (s *Section) closeTag(buf *bytes.Buffer) {
	buf.WriteString("</")
	buf.WriteString(SectionTagName)
	buf.WriteByte('>')
}

func (r *Row) buildHTML(ds DataStore, mc ModuleCustomizer, buf *bytes.Buffer, css *PageCSS) error {

	if r.Common.Type == StaticElementType {
		buf.Write(r.Html)
		return nil
	}

	if r.Dyn != 0 {

		dynData, err := ds.Get(r.Dyn)
		if err != nil {
			return err
		}

		var dynSettings Row
		err = proto.Unmarshal(dynData.Data, &dynSettings)
		if err != nil {
			return err
		}

		r.buildClassesAndStyles(&dynSettings, css)

	}

	r.openTag(buf)

	for _, col := range r.Columns {
		buf.WriteString(ColumnOpen)
		for _, mod := range col.Modules {
			if err := mod.buildHTML(ds, mc, buf, css); err != nil {
				return err
			}
		}
		buf.WriteString(ColumnClose)
	}

	r.closeTag(buf)

	return nil

}

func (r *Row) buildClassesAndStyles(settings *Row, css *PageCSS) {

	// Row elements have the "room-row" class.
	r.Common.Classes = makeClassesList("row", settings.Common.Type, r.Dyn, settings.Common.Classes)

	if len(r.Common.Styles) > 0 {
		css.WriteByte('.')
		css.WriteString(r.Common.Classes[2])
		css.WriteByte('{')
		moreClasses := CompileStyles(settings.Common.Options, css)
		css.WriteByte('}')

		r.Common.Classes = append(r.Common.Classes, moreClasses...)
	}

}

func (r *Row) openTag(buf *bytes.Buffer) {
	ElementOpenTag(RowTagName, r.Common.IdAttr, r.Common.Classes, buf)
}

func (r *Row) closeTag(buf *bytes.Buffer) {
	buf.WriteString("</")
	buf.WriteString(RowTagName)
	buf.WriteByte('>')
}

// buildHTML uses the mc ModuleCustomizer to build the module. The mc passed in must not be nil.
// If the module's is static (as indicated by the Type), then only the Data is written out to page.
func (m *Module) buildHTML(ds DataStore, mc ModuleCustomizer, page *bytes.Buffer, css *PageCSS) error {
	if m.Type == StaticElementType {
		page.Write(m.Data)
		return nil
	}
	builder := mc.GetModuleBuilder(m)
	if builder == nil {
		return ErrMissingModuleBuilder{m}
	}
	return builder.BuildHTML(ds, mc, m, page, css)
}
