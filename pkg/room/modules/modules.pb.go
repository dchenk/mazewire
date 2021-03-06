// Code generated by protoc-gen-go. DO NOT EDIT.
// source: room/modules/modules.proto

package modules

import (
	fmt "fmt"
	room "github.com/dchenk/mazewire/pkg/room"
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

// A NavLink contains the information needed to render a navigation link.
type NavLink struct {
	Text string `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	Href string `protobuf:"bytes,2,opt,name=href,proto3" json:"href,omitempty"`
	// Target is either blank or something like "_blank".
	Target               string   `protobuf:"bytes,3,opt,name=target,proto3" json:"target,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NavLink) Reset()         { *m = NavLink{} }
func (m *NavLink) String() string { return proto.CompactTextString(m) }
func (*NavLink) ProtoMessage()    {}
func (*NavLink) Descriptor() ([]byte, []int) {
	return fileDescriptor_7c9af2b85f98e5df, []int{0}
}

func (m *NavLink) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NavLink.Unmarshal(m, b)
}
func (m *NavLink) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NavLink.Marshal(b, m, deterministic)
}
func (m *NavLink) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NavLink.Merge(m, src)
}
func (m *NavLink) XXX_Size() int {
	return xxx_messageInfo_NavLink.Size(m)
}
func (m *NavLink) XXX_DiscardUnknown() {
	xxx_messageInfo_NavLink.DiscardUnknown(m)
}

var xxx_messageInfo_NavLink proto.InternalMessageInfo

func (m *NavLink) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

func (m *NavLink) GetHref() string {
	if m != nil {
		return m.Href
	}
	return ""
}

func (m *NavLink) GetTarget() string {
	if m != nil {
		return m.Target
	}
	return ""
}

// A Nav is a collection of NavLinks.
type Nav struct {
	Common               *room.Common `protobuf:"bytes,1,opt,name=common,proto3" json:"common,omitempty"`
	Links                []*NavLink   `protobuf:"bytes,2,rep,name=links,proto3" json:"links,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Nav) Reset()         { *m = Nav{} }
func (m *Nav) String() string { return proto.CompactTextString(m) }
func (*Nav) ProtoMessage()    {}
func (*Nav) Descriptor() ([]byte, []int) {
	return fileDescriptor_7c9af2b85f98e5df, []int{1}
}

func (m *Nav) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Nav.Unmarshal(m, b)
}
func (m *Nav) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Nav.Marshal(b, m, deterministic)
}
func (m *Nav) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Nav.Merge(m, src)
}
func (m *Nav) XXX_Size() int {
	return xxx_messageInfo_Nav.Size(m)
}
func (m *Nav) XXX_DiscardUnknown() {
	xxx_messageInfo_Nav.DiscardUnknown(m)
}

var xxx_messageInfo_Nav proto.InternalMessageInfo

func (m *Nav) GetCommon() *room.Common {
	if m != nil {
		return m.Common
	}
	return nil
}

func (m *Nav) GetLinks() []*NavLink {
	if m != nil {
		return m.Links
	}
	return nil
}

// HTMLModule is a module that represents any kind of custom HTML. HTMLModule does not have a ModuleBuilder
// implementation because its data is always static.
type HTML struct {
	Tag                  string   `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
	IdAttr               string   `protobuf:"bytes,2,opt,name=id_attr,json=idAttr,proto3" json:"id_attr,omitempty"`
	Html                 []byte   `protobuf:"bytes,3,opt,name=html,proto3" json:"html,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HTML) Reset()         { *m = HTML{} }
func (m *HTML) String() string { return proto.CompactTextString(m) }
func (*HTML) ProtoMessage()    {}
func (*HTML) Descriptor() ([]byte, []int) {
	return fileDescriptor_7c9af2b85f98e5df, []int{2}
}

func (m *HTML) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HTML.Unmarshal(m, b)
}
func (m *HTML) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HTML.Marshal(b, m, deterministic)
}
func (m *HTML) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HTML.Merge(m, src)
}
func (m *HTML) XXX_Size() int {
	return xxx_messageInfo_HTML.Size(m)
}
func (m *HTML) XXX_DiscardUnknown() {
	xxx_messageInfo_HTML.DiscardUnknown(m)
}

var xxx_messageInfo_HTML proto.InternalMessageInfo

func (m *HTML) GetTag() string {
	if m != nil {
		return m.Tag
	}
	return ""
}

func (m *HTML) GetIdAttr() string {
	if m != nil {
		return m.IdAttr
	}
	return ""
}

func (m *HTML) GetHtml() []byte {
	if m != nil {
		return m.Html
	}
	return nil
}

// Image is an image, optionally wrapped in an anchor (<a>) tag.
type Image struct {
	Common               *room.Common `protobuf:"bytes,1,opt,name=common,proto3" json:"common,omitempty"`
	Src                  string       `protobuf:"bytes,2,opt,name=src,proto3" json:"src,omitempty"`
	Alt                  string       `protobuf:"bytes,3,opt,name=alt,proto3" json:"alt,omitempty"`
	LinkUrl              string       `protobuf:"bytes,4,opt,name=link_url,json=linkUrl,proto3" json:"link_url,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Image) Reset()         { *m = Image{} }
func (m *Image) String() string { return proto.CompactTextString(m) }
func (*Image) ProtoMessage()    {}
func (*Image) Descriptor() ([]byte, []int) {
	return fileDescriptor_7c9af2b85f98e5df, []int{3}
}

func (m *Image) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Image.Unmarshal(m, b)
}
func (m *Image) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Image.Marshal(b, m, deterministic)
}
func (m *Image) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Image.Merge(m, src)
}
func (m *Image) XXX_Size() int {
	return xxx_messageInfo_Image.Size(m)
}
func (m *Image) XXX_DiscardUnknown() {
	xxx_messageInfo_Image.DiscardUnknown(m)
}

var xxx_messageInfo_Image proto.InternalMessageInfo

func (m *Image) GetCommon() *room.Common {
	if m != nil {
		return m.Common
	}
	return nil
}

func (m *Image) GetSrc() string {
	if m != nil {
		return m.Src
	}
	return ""
}

func (m *Image) GetAlt() string {
	if m != nil {
		return m.Alt
	}
	return ""
}

func (m *Image) GetLinkUrl() string {
	if m != nil {
		return m.LinkUrl
	}
	return ""
}

// Text is a rich text module created using Quill's Delta format (https://quilljs.com/docs/delta).
type Text struct {
	Common               *room.Common `protobuf:"bytes,1,opt,name=common,proto3" json:"common,omitempty"`
	Ops                  []byte       `protobuf:"bytes,2,opt,name=ops,proto3" json:"ops,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Text) Reset()         { *m = Text{} }
func (m *Text) String() string { return proto.CompactTextString(m) }
func (*Text) ProtoMessage()    {}
func (*Text) Descriptor() ([]byte, []int) {
	return fileDescriptor_7c9af2b85f98e5df, []int{4}
}

func (m *Text) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Text.Unmarshal(m, b)
}
func (m *Text) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Text.Marshal(b, m, deterministic)
}
func (m *Text) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Text.Merge(m, src)
}
func (m *Text) XXX_Size() int {
	return xxx_messageInfo_Text.Size(m)
}
func (m *Text) XXX_DiscardUnknown() {
	xxx_messageInfo_Text.DiscardUnknown(m)
}

var xxx_messageInfo_Text proto.InternalMessageInfo

func (m *Text) GetCommon() *room.Common {
	if m != nil {
		return m.Common
	}
	return nil
}

func (m *Text) GetOps() []byte {
	if m != nil {
		return m.Ops
	}
	return nil
}

func init() {
	proto.RegisterType((*NavLink)(nil), "modules.NavLink")
	proto.RegisterType((*Nav)(nil), "modules.Nav")
	proto.RegisterType((*HTML)(nil), "modules.HTML")
	proto.RegisterType((*Image)(nil), "modules.Image")
	proto.RegisterType((*Text)(nil), "modules.Text")
}

func init() { proto.RegisterFile("room/modules/modules.proto", fileDescriptor_7c9af2b85f98e5df) }

var fileDescriptor_7c9af2b85f98e5df = []byte{
	// 314 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x91, 0x4f, 0x6b, 0xb3, 0x40,
	0x10, 0x87, 0x49, 0x34, 0xfa, 0xbe, 0x93, 0x1c, 0x64, 0x0f, 0xad, 0xcd, 0x29, 0x48, 0x29, 0x81,
	0x52, 0x85, 0xf4, 0x5e, 0x68, 0x4b, 0xa1, 0x81, 0x34, 0x07, 0x9b, 0x5e, 0x7a, 0x09, 0x1b, 0xdd,
	0x1a, 0x71, 0xd7, 0x95, 0x75, 0x4c, 0xff, 0x7c, 0xfa, 0xb2, 0x9b, 0x15, 0x7a, 0xcc, 0xc9, 0x67,
	0x7e, 0x23, 0x33, 0x8f, 0x0e, 0x4c, 0x95, 0x94, 0x22, 0x11, 0x32, 0xef, 0x38, 0x6b, 0xfb, 0x67,
	0xdc, 0x28, 0x89, 0x92, 0xf8, 0xb6, 0x9c, 0x06, 0xe6, 0x25, 0xfc, 0x6e, 0xfa, 0x56, 0xb4, 0x04,
	0x7f, 0x4d, 0x0f, 0xab, 0xb2, 0xae, 0x08, 0x01, 0x17, 0xd9, 0x17, 0x86, 0x83, 0xd9, 0x60, 0xfe,
	0x3f, 0x35, 0xac, 0xb3, 0xbd, 0x62, 0x1f, 0xe1, 0xf0, 0x98, 0x69, 0x26, 0x67, 0xe0, 0x21, 0x55,
	0x05, 0xc3, 0xd0, 0x31, 0xa9, 0xad, 0xa2, 0x57, 0x70, 0xd6, 0xf4, 0x40, 0x2e, 0xc1, 0xcb, 0xa4,
	0x10, 0xb2, 0x36, 0x83, 0xc6, 0x8b, 0x49, 0xac, 0x97, 0xc6, 0x8f, 0x26, 0x4b, 0x6d, 0x8f, 0x5c,
	0xc1, 0x88, 0x97, 0x75, 0xd5, 0x86, 0xc3, 0x99, 0x33, 0x1f, 0x2f, 0x82, 0xb8, 0x37, 0xb6, 0x36,
	0xe9, 0xb1, 0x1d, 0x3d, 0x81, 0xfb, 0xbc, 0x79, 0x59, 0x91, 0x00, 0x1c, 0xa4, 0x85, 0x75, 0xd3,
	0x48, 0xce, 0xc1, 0x2f, 0xf3, 0x2d, 0x45, 0x54, 0xd6, 0xce, 0x2b, 0xf3, 0x7b, 0x44, 0x65, 0x9c,
	0x51, 0x70, 0x63, 0x37, 0x49, 0x0d, 0x47, 0x1c, 0x46, 0x4b, 0x41, 0x0b, 0x76, 0xa2, 0x5d, 0x00,
	0x4e, 0xab, 0x32, 0x3b, 0x57, 0xa3, 0x4e, 0x28, 0xef, 0xbf, 0x58, 0x23, 0xb9, 0x80, 0x7f, 0x5a,
	0x71, 0xdb, 0x29, 0x1e, 0xba, 0x26, 0xf6, 0x75, 0xfd, 0xa6, 0x78, 0x74, 0x07, 0xee, 0x46, 0xff,
	0xbd, 0x93, 0x97, 0xc9, 0xa6, 0x35, 0xcb, 0x26, 0xa9, 0xc6, 0x87, 0x9b, 0xf7, 0xeb, 0xa2, 0xc4,
	0x7d, 0xb7, 0x8b, 0x33, 0x29, 0x92, 0x3c, 0xdb, 0xb3, 0xba, 0x4a, 0x04, 0xfd, 0x61, 0x9f, 0xa5,
	0x62, 0x49, 0x53, 0x15, 0xc9, 0xdf, 0x63, 0xef, 0x3c, 0x73, 0xca, 0xdb, 0xdf, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x8a, 0x15, 0x0c, 0x30, 0x03, 0x02, 0x00, 0x00,
}
