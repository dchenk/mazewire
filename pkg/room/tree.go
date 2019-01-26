package room

import (
	"bytes"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

// StaticElementType is the Type given to those Section, Row, and Module elements whose final HTML form is
// known after compilation.
const StaticElementType = "_static"

// TreeType is the Type given to the Module elements that contain a Tree in the Data.
const TreeType = "room_tree"

// SectionTagName is the HTML tag name given to Section elements.
const SectionTagName = "section"

func (s *Section) compile(ds DataStore, ms ModuleCustomizer, css *PageCSS) (*Section, error) {

	if s.Common.Options == nil {
		s.Common.Options = make(map[string]string)
	}

	compiled := &Section{
		Common: &Common{
			Type:   s.Common.Type,
			IdAttr: strings.TrimSpace(s.Common.IdAttr),
		},
		Dyn:  s.Dyn,
		Rows: make([]*Row, len(s.Rows)),
	}

	allStatic := true // Assume at first that all rows are going to be static HTML (nothing is dynamic).

	for rowI := range s.Rows {
		newRow, err := s.Rows[rowI].compile(ds, ms, css)
		if err != nil {
			return compiled, err
		}
		compiled.Rows[rowI] = newRow

		// If any row inside this section is dynamic or has a dynamic module, then this entire section is
		// not fully static.
		if newRow.Common.Type != StaticElementType {
			allStatic = false
		}
	}

	if s.Dyn == 0 {
		compiled.buildClassesAndStyles(s, css)

		// allStatic will still be false here if some rows or modules are dynamic.
		if allStatic {
			compiled.Common.Type = StaticElementType
			var b bytes.Buffer
			compiled.openTag(&b)
			for i := range compiled.Rows {
				b.Write(compiled.Rows[i].Html)
			}
			compiled.closeTag(&b)
			compiled.Html = b.Bytes()
		}
	}

	return compiled, nil

}

// RowTagName is the HTML tag name given to Row elements.
const RowTagName = "div"

func (r *Row) compile(ds DataStore, mc ModuleCustomizer, css *PageCSS) (*Row, error) {

	if r.Common.Options == nil {
		r.Common.Options = make(map[string]string)
	}

	cr := &Row{
		Common: &Common{
			Type:   r.Common.Type,
			IdAttr: strings.TrimSpace(r.Common.IdAttr),
		},
		Dyn:     r.Dyn,
		Columns: make([]*Column, len(r.Columns)),
	}

	allStatic := true // Assume at first that all modules are going to be static HTML (nothing is dynamic).

	for i, col := range r.Columns {

		cr.Columns[i].Modules = make([]*Module, len(col.Modules))

		var err error

		for mi, mod := range col.Modules {

			cr.Columns[i].Modules[mi], err = mod.compile(ds, mc, css)
			if err != nil {
				return cr, err
			}

			// If any module inside this row is dynamic, then this entire row is not fully static.
			if cr.Columns[i].Modules[mi].Type != StaticElementType {
				allStatic = false
			}

		}

	}

	if r.Dyn == 0 {
		cr.buildClassesAndStyles(r, css)

		// allStatic will still be false here if some modules are dynamic.
		if allStatic {
			cr.Common.Type = StaticElementType
			var b bytes.Buffer
			cr.openTag(&b)
			for _, col := range cr.Columns {
				b.WriteString(ColumnOpen)
				for _, mod := range col.Modules {
					b.Write(mod.Data)
				}
				b.WriteString(ColumnClose)
			}
			cr.closeTag(&b)
			cr.Html = b.Bytes()
		}
	}

	return cr, nil

}

// compile compiles the Module using the ModuleCompiler that mc returns. The mc argument must not be nil.
func (m *Module) compile(ds DataStore, mc ModuleCustomizer, css *PageCSS) (*Module, error) {
	compiler := mc.GetModuleCompiler(m)
	if compiler == nil {
		return nil, ErrMissingModuleCompiler{m}
	}
	return compiler.Compile(ds, mc, m, css)
}

// Compile takes the contents of the Tree and compiles all sections, rows, and modules as much as possible,
// the remaining elements left to being built dynamically at the time a page is requested.
//
// The m.Data field is unmarshalled and used to construct the Tree.
func (t *Tree) Compile(ds DataStore, mc ModuleCustomizer, m *Module, css *PageCSS) (*Module, error) {
	err := proto.Unmarshal(m.Data, t)
	if err != nil {
		return nil, err
	}

	if mc == nil {
		return nil, errors.New("room: Compile was given a nil ModuleCustomizer")
	}

	if err = t.cacheDynamicData(ds); err != nil {
		return nil, err
	}

	// compiled will become the data within the *Module returned.
	// If the entire tree turns turns out to be static, then the Data field will contain all the built HTML.
	compiled := Tree{
		Sections: make([]*Section, len(t.Sections)),
	}

	allStatic := true // Indicate if contents are all static HTML.

	// Compile the tree one section at a time.
	for i, sec := range t.Sections {
		compiled.Sections[i], err = sec.compile(ds, mc, css)
		if err != nil {
			return nil, err
		}
		if compiled.Sections[i].Common.Type != StaticElementType {
			allStatic = false
		}
	}

	if allStatic {
		cm := Module{Type: StaticElementType}
		var totalSize int
		for _, sec := range compiled.Sections {
			totalSize += len(sec.Html)
		}
		cm.Data = make([]byte, 0, totalSize)
		for _, sec := range compiled.Sections {
			cm.Data = append(cm.Data, sec.Html...)
		}
		return &cm, nil
	}

	// Check for adjacent sections that are static HTML and combine them.
	compiled.mergeStaticSections()

	cm := Module{Type: TreeType}

	cm.Data, err = proto.Marshal(&compiled)
	return &cm, err
}

// BuildHTML implements ModuleBuilder for Tree. This function writes the built HTML to the page buffer.
// The Data field of the compiled Module is unmarshalled and used to construct the Tree.
func (t *Tree) BuildHTML(ds DataStore, mc ModuleCustomizer, compiled *Module, page *bytes.Buffer, css *PageCSS) error {
	err := proto.Unmarshal(compiled.Data, t)
	if err != nil {
		return err
	}
	if mc == nil {
		return errors.New("room: BuildHTML was given a nil ModuleCustomizer")
	}
	if err = t.cacheDynamicData(ds); err != nil {
		return err
	}
	for _, section := range t.Sections {
		if err = section.buildHTML(ds, mc, page, css); err != nil {
			return err
		}
	}
	return nil
}

// mergeStaticSections loops through the sections of the tree and merges the data of adjacent static
// sections.
func (t *Tree) mergeStaticSections() {
	var sec, secNext *Section
	for i := 0; i < len(t.Sections)-1; i++ {
		sec = t.Sections[i]
		secNext = t.Sections[i+1]
		if sec.Common.Type == StaticElementType && secNext.Common.Type == StaticElementType {
			sec.Html = append(sec.Html, secNext.Html...)
			t.Sections = append(t.Sections[:i+1], t.Sections[i+2:]...)
			i--
		}
	}
}

// DynamicDataIDs returns the IDs of the Datum entities that contain the dynamic settings for this Tree.
func (t *Tree) DynamicDataIDs() []int64 {
	IDs := make([]int64, 0, len(t.Sections)*2) // rough length estimate
	for _, sec := range t.Sections {
		if sec.Dyn != 0 && !containsInt64(IDs, sec.Dyn) {
			IDs = append(IDs, sec.Dyn)
		}
		for _, row := range sec.Rows {
			if row.Dyn != 0 && !containsInt64(IDs, row.Dyn) {
				IDs = append(IDs, row.Dyn)
			}
			for _, col := range row.Columns {
				for _, mod := range col.Modules {
					if mod.Dyn != 0 && !containsInt64(IDs, mod.Dyn) {
						IDs = append(IDs, mod.Dyn)
					}
				}
			}
		}
	}
	return IDs
}

// cacheDynamicData retrieves all the dynamic data for the Tree to be cached by the DataStore.
// This helps us get as much data as we can with a single query.
func (t *Tree) cacheDynamicData(ds DataStore) error {
	if dynIDs := t.DynamicDataIDs(); len(dynIDs) > 0 {
		_, err := ds.GetMulti(dynIDs)
		return err
	}
	return nil
}

// containsInt64 says if slice a has n as an element.
func containsInt64(a []int64, n int64) bool {
	for i := range a {
		if a[i] == n {
			return true
		}
	}
	return false
}
