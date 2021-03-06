// Code generated by protoc-gen-go. DO NOT EDIT.
// source: room/types.proto

package room

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// A Module represents the smallest unit of composition that appears as an element in the HTML of a page.
// We know how to deal with the Data field when we know the Type and have a matching ModuleCompiler to compile
// it or a matching ModuleBuilder to build its final HTML.
//
// A compiled module may be fully static and ready to be written out to the page without building, in which
// case the data contains the final HTML and the type is "_static".
type Module struct {
	// Type indicates broadly the kind of module that this is. It is used to identify what ModuleCompiler or
	// ModuleBuilder to use.
	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	// Dyn is the ID of the record holding the module's dynamic settings. Using the Dyn field is optional for
	// each module type, but if it is non-zero then it must be the ID of the Datum entity which holds the
	// dynamic settings. (This requirement is intended to make it possible to get all of the data for the
	// entire page tree with a single query before the page is built.)
	Dyn int64 `protobuf:"varint,2,opt,name=dyn,proto3" json:"dyn,omitempty"`
	// Data contains all of the configuration data needed for the module. The data can be encoded in any
	// format, and it's up to the corresponding ModuleCompiler and ModuleBuilder to be able to decode it.
	//
	// When a Module is compiled, if Compile returns true then Data must already contain all the pre-built
	// HTML.
	Data                 []byte   `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Module) Reset()         { *m = Module{} }
func (m *Module) String() string { return proto.CompactTextString(m) }
func (*Module) ProtoMessage()    {}
func (*Module) Descriptor() ([]byte, []int) {
	return fileDescriptor_a57991ff5698184a, []int{0}
}

func (m *Module) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Module.Unmarshal(m, b)
}
func (m *Module) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Module.Marshal(b, m, deterministic)
}
func (m *Module) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Module.Merge(m, src)
}
func (m *Module) XXX_Size() int {
	return xxx_messageInfo_Module.Size(m)
}
func (m *Module) XXX_DiscardUnknown() {
	xxx_messageInfo_Module.DiscardUnknown(m)
}

var xxx_messageInfo_Module proto.InternalMessageInfo

func (m *Module) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Module) GetDyn() int64 {
	if m != nil {
		return m.Dyn
	}
	return 0
}

func (m *Module) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type Column struct {
	Modules              []*Module `protobuf:"bytes,1,rep,name=modules,proto3" json:"modules,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Column) Reset()         { *m = Column{} }
func (m *Column) String() string { return proto.CompactTextString(m) }
func (*Column) ProtoMessage()    {}
func (*Column) Descriptor() ([]byte, []int) {
	return fileDescriptor_a57991ff5698184a, []int{1}
}

func (m *Column) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Column.Unmarshal(m, b)
}
func (m *Column) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Column.Marshal(b, m, deterministic)
}
func (m *Column) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Column.Merge(m, src)
}
func (m *Column) XXX_Size() int {
	return xxx_messageInfo_Column.Size(m)
}
func (m *Column) XXX_DiscardUnknown() {
	xxx_messageInfo_Column.DiscardUnknown(m)
}

var xxx_messageInfo_Column proto.InternalMessageInfo

func (m *Column) GetModules() []*Module {
	if m != nil {
		return m.Modules
	}
	return nil
}

// Common is a data structure commonly used in this package in various other types.
type Common struct {
	// Type indicates the broad type of the structure.
	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	// IdAttr is the user's custom ID attribute for the element.
	IdAttr string `protobuf:"bytes,2,opt,name=id_attr,json=idAttr,proto3" json:"id_attr,omitempty"`
	// Classes is the user's custom classes.
	// For compiled sections and rows, this includes classes added by module compilers for styling.
	Classes []string `protobuf:"bytes,3,rep,name=classes,proto3" json:"classes,omitempty"`
	// Styles is the user's styles for the row.
	Styles map[string]string `protobuf:"bytes,4,rep,name=styles,proto3" json:"styles,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Options contains additional options.
	Options              map[string]string `protobuf:"bytes,5,rep,name=options,proto3" json:"options,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Common) Reset()         { *m = Common{} }
func (m *Common) String() string { return proto.CompactTextString(m) }
func (*Common) ProtoMessage()    {}
func (*Common) Descriptor() ([]byte, []int) {
	return fileDescriptor_a57991ff5698184a, []int{2}
}

func (m *Common) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Common.Unmarshal(m, b)
}
func (m *Common) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Common.Marshal(b, m, deterministic)
}
func (m *Common) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Common.Merge(m, src)
}
func (m *Common) XXX_Size() int {
	return xxx_messageInfo_Common.Size(m)
}
func (m *Common) XXX_DiscardUnknown() {
	xxx_messageInfo_Common.DiscardUnknown(m)
}

var xxx_messageInfo_Common proto.InternalMessageInfo

func (m *Common) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Common) GetIdAttr() string {
	if m != nil {
		return m.IdAttr
	}
	return ""
}

func (m *Common) GetClasses() []string {
	if m != nil {
		return m.Classes
	}
	return nil
}

func (m *Common) GetStyles() map[string]string {
	if m != nil {
		return m.Styles
	}
	return nil
}

func (m *Common) GetOptions() map[string]string {
	if m != nil {
		return m.Options
	}
	return nil
}

// A Row holds the settings of both static and dynamic page rows. Directly within a Row are columns.
// A compiled Row only needs to be built, along with everything it contains, to be displayed.
type Row struct {
	// Common includes common fields. The Type is the column layout type of the row.
	Common *Common `protobuf:"bytes,1,opt,name=common,proto3" json:"common,omitempty"`
	// Dyn is the ID of the row's dynamic settings, if not zero.
	Dyn int64 `protobuf:"varint,2,opt,name=dyn,proto3" json:"dyn,omitempty"`
	// Columns contains the row's columns. For compiled Rows, this field is not empty only if there is content
	// within the row that is dynamic.
	Columns []*Column `protobuf:"bytes,3,rep,name=columns,proto3" json:"columns,omitempty"`
	// Html is used only for compiled rows when the entire row is static HTML.
	Html                 []byte   `protobuf:"bytes,4,opt,name=html,proto3" json:"html,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Row) Reset()         { *m = Row{} }
func (m *Row) String() string { return proto.CompactTextString(m) }
func (*Row) ProtoMessage()    {}
func (*Row) Descriptor() ([]byte, []int) {
	return fileDescriptor_a57991ff5698184a, []int{3}
}

func (m *Row) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Row.Unmarshal(m, b)
}
func (m *Row) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Row.Marshal(b, m, deterministic)
}
func (m *Row) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Row.Merge(m, src)
}
func (m *Row) XXX_Size() int {
	return xxx_messageInfo_Row.Size(m)
}
func (m *Row) XXX_DiscardUnknown() {
	xxx_messageInfo_Row.DiscardUnknown(m)
}

var xxx_messageInfo_Row proto.InternalMessageInfo

func (m *Row) GetCommon() *Common {
	if m != nil {
		return m.Common
	}
	return nil
}

func (m *Row) GetDyn() int64 {
	if m != nil {
		return m.Dyn
	}
	return 0
}

func (m *Row) GetColumns() []*Column {
	if m != nil {
		return m.Columns
	}
	return nil
}

func (m *Row) GetHtml() []byte {
	if m != nil {
		return m.Html
	}
	return nil
}

// A Section holds the settings of a section in a Tree. Within a Section are Row elements.
// A compiled Section only needs to be built, along with everything it contains, to be displayed.
type Section struct {
	// Common includes common fields. The Type must be either "standard" or "full".
	Common *Common `protobuf:"bytes,1,opt,name=common,proto3" json:"common,omitempty"`
	// Dyn is the ID of the section's dynamic settings, if not zero.
	Dyn int64 `protobuf:"varint,2,opt,name=dyn,proto3" json:"dyn,omitempty"`
	// Rows contains the section's rows.
	Rows []*Row `protobuf:"bytes,3,rep,name=rows,proto3" json:"rows,omitempty"`
	// Html is used only for compiled rows when the entire row is static HTML.
	Html                 []byte   `protobuf:"bytes,4,opt,name=html,proto3" json:"html,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Section) Reset()         { *m = Section{} }
func (m *Section) String() string { return proto.CompactTextString(m) }
func (*Section) ProtoMessage()    {}
func (*Section) Descriptor() ([]byte, []int) {
	return fileDescriptor_a57991ff5698184a, []int{4}
}

func (m *Section) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Section.Unmarshal(m, b)
}
func (m *Section) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Section.Marshal(b, m, deterministic)
}
func (m *Section) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Section.Merge(m, src)
}
func (m *Section) XXX_Size() int {
	return xxx_messageInfo_Section.Size(m)
}
func (m *Section) XXX_DiscardUnknown() {
	xxx_messageInfo_Section.DiscardUnknown(m)
}

var xxx_messageInfo_Section proto.InternalMessageInfo

func (m *Section) GetCommon() *Common {
	if m != nil {
		return m.Common
	}
	return nil
}

func (m *Section) GetDyn() int64 {
	if m != nil {
		return m.Dyn
	}
	return 0
}

func (m *Section) GetRows() []*Row {
	if m != nil {
		return m.Rows
	}
	return nil
}

func (m *Section) GetHtml() []byte {
	if m != nil {
		return m.Html
	}
	return nil
}

// A Tree represents the structure of a page. Each Section within a Tree may contain Row elements, which
// contain Column elements, which contain Module elements.
//
// Tree implements the ModuleCompiler and ModuleBuilder interfaces, so a Tree can be nested as a Module.
//
// A compiled Tree represents the structure of a compiled page, though some elements may need to be compiled
// dynamically. Each compiled Section within a Tree contains compiled Row elements, each of which may in turn
// contain Column elements with compiled Module elements.
type Tree struct {
	Sections             []*Section `protobuf:"bytes,1,rep,name=sections,proto3" json:"sections,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Tree) Reset()         { *m = Tree{} }
func (m *Tree) String() string { return proto.CompactTextString(m) }
func (*Tree) ProtoMessage()    {}
func (*Tree) Descriptor() ([]byte, []int) {
	return fileDescriptor_a57991ff5698184a, []int{5}
}

func (m *Tree) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Tree.Unmarshal(m, b)
}
func (m *Tree) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Tree.Marshal(b, m, deterministic)
}
func (m *Tree) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Tree.Merge(m, src)
}
func (m *Tree) XXX_Size() int {
	return xxx_messageInfo_Tree.Size(m)
}
func (m *Tree) XXX_DiscardUnknown() {
	xxx_messageInfo_Tree.DiscardUnknown(m)
}

var xxx_messageInfo_Tree proto.InternalMessageInfo

func (m *Tree) GetSections() []*Section {
	if m != nil {
		return m.Sections
	}
	return nil
}

func init() {
	proto.RegisterType((*Module)(nil), "room.Module")
	proto.RegisterType((*Column)(nil), "room.Column")
	proto.RegisterType((*Common)(nil), "room.Common")
	proto.RegisterMapType((map[string]string)(nil), "room.Common.OptionsEntry")
	proto.RegisterMapType((map[string]string)(nil), "room.Common.StylesEntry")
	proto.RegisterType((*Row)(nil), "room.Row")
	proto.RegisterType((*Section)(nil), "room.Section")
	proto.RegisterType((*Tree)(nil), "room.Tree")
}

func init() { proto.RegisterFile("room/types.proto", fileDescriptor_a57991ff5698184a) }

var fileDescriptor_a57991ff5698184a = []byte{
	// 402 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x92, 0xc1, 0xeb, 0xd3, 0x30,
	0x14, 0xc7, 0xe9, 0xda, 0x5f, 0x6b, 0xdf, 0x26, 0x8c, 0x20, 0x18, 0x07, 0x42, 0xa9, 0x3a, 0xea,
	0xa5, 0x9d, 0xee, 0xa2, 0xbb, 0x39, 0xf1, 0x28, 0x42, 0xe6, 0xc9, 0x8b, 0x74, 0x6d, 0xd8, 0xca,
	0x9a, 0xa6, 0x24, 0xa9, 0xa5, 0xfe, 0x4f, 0xfe, 0x8f, 0x92, 0xa4, 0x1d, 0x13, 0x76, 0xf1, 0x77,
	0x7b, 0xc9, 0x37, 0xdf, 0xf7, 0xf9, 0xf2, 0x5e, 0x60, 0x29, 0x38, 0x67, 0x99, 0x1a, 0x5a, 0x2a,
	0xd3, 0x56, 0x70, 0xc5, 0x91, 0xa7, 0x6f, 0xe2, 0x3d, 0xf8, 0x5f, 0x79, 0xd9, 0xd5, 0x14, 0x21,
	0xf0, 0xb4, 0x8c, 0x9d, 0xc8, 0x49, 0x42, 0x62, 0x6a, 0xb4, 0x04, 0xb7, 0x1c, 0x1a, 0x3c, 0x8b,
	0x9c, 0xc4, 0x25, 0xba, 0xd4, 0xaf, 0xca, 0x5c, 0xe5, 0xd8, 0x8d, 0x9c, 0x64, 0x41, 0x4c, 0x1d,
	0x6f, 0xc0, 0xff, 0xcc, 0xeb, 0x8e, 0x35, 0x68, 0x0d, 0x01, 0x33, 0xdd, 0x24, 0x76, 0x22, 0x37,
	0x99, 0xbf, 0x5f, 0xa4, 0x9a, 0x92, 0x5a, 0x04, 0x99, 0xc4, 0xf8, 0xcf, 0x4c, 0x5b, 0x18, 0xe3,
	0xcd, 0x5d, 0xec, 0x73, 0x08, 0xaa, 0xf2, 0x67, 0xae, 0x94, 0x30, 0xe8, 0x90, 0xf8, 0x55, 0xf9,
	0x49, 0x29, 0x81, 0x30, 0x04, 0x45, 0x9d, 0x4b, 0x49, 0x25, 0x76, 0x23, 0x37, 0x09, 0xc9, 0x74,
	0x44, 0x1b, 0xf0, 0xa5, 0x1a, 0x34, 0xd8, 0x33, 0x60, 0x6c, 0xc1, 0x16, 0x92, 0x1e, 0x8c, 0xf4,
	0xa5, 0x51, 0x62, 0x20, 0xe3, 0x3b, 0xb4, 0x85, 0x80, 0xb7, 0xaa, 0xe2, 0x8d, 0xc4, 0x0f, 0xc6,
	0xf2, 0xe2, 0x1f, 0xcb, 0x37, 0xab, 0x59, 0xcf, 0xf4, 0x72, 0xf5, 0x11, 0xe6, 0x37, 0xbd, 0xf4,
	0x7c, 0x2e, 0x74, 0x18, 0xb3, 0xeb, 0x12, 0x3d, 0x83, 0x87, 0x5f, 0x79, 0xdd, 0xd1, 0x31, 0xb8,
	0x3d, 0xec, 0x66, 0x1f, 0x9c, 0xd5, 0x0e, 0x16, 0xb7, 0x3d, 0xff, 0xc7, 0x1b, 0x0f, 0xe0, 0x12,
	0xde, 0xa3, 0xd7, 0xe0, 0x17, 0x26, 0x9d, 0x71, 0x5d, 0xa7, 0x6b, 0x13, 0x93, 0x51, 0xbb, 0xb3,
	0xb4, 0x35, 0x04, 0x85, 0x59, 0x90, 0x1d, 0xdb, 0x8d, 0x51, 0x5f, 0x92, 0x49, 0xd4, 0xbb, 0x38,
	0x2b, 0x56, 0x63, 0xcf, 0x2e, 0x57, 0xd7, 0xb1, 0x82, 0xe0, 0x40, 0x0b, 0x9d, 0xfb, 0xd1, 0xf8,
	0x97, 0xe0, 0x09, 0xde, 0x4f, 0xec, 0xd0, 0xba, 0x08, 0xef, 0x89, 0xb9, 0xbe, 0x4b, 0x7d, 0x07,
	0xde, 0x77, 0x41, 0x29, 0x7a, 0x0b, 0x4f, 0xa4, 0xa5, 0x4f, 0x3f, 0xea, 0xa9, 0xb5, 0x8f, 0x99,
	0xc8, 0x55, 0xde, 0xbf, 0xf9, 0xf1, 0xea, 0x54, 0xa9, 0x73, 0x77, 0x4c, 0x0b, 0xce, 0xb2, 0xb2,
	0x38, 0xd3, 0xe6, 0x92, 0xb1, 0xfc, 0x37, 0xed, 0x2b, 0x41, 0xb3, 0xf6, 0x72, 0xca, 0xb4, 0xf1,
	0xe8, 0x9b, 0xdf, 0xbf, 0xfd, 0x1b, 0x00, 0x00, 0xff, 0xff, 0x49, 0xf6, 0x63, 0x10, 0x11, 0x03,
	0x00, 0x00,
}
