// Code generated by protoc-gen-go. DO NOT EDIT.
// source: types/version/version.proto

package version

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

// Version is a semantic versioning version of a module.
type Version struct {
	Major                uint32   `protobuf:"varint,1,opt,name=major,proto3" json:"major,omitempty"`
	Minor                uint32   `protobuf:"varint,2,opt,name=minor,proto3" json:"minor,omitempty"`
	Patch                uint32   `protobuf:"varint,3,opt,name=patch,proto3" json:"patch,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Version) Reset()         { *m = Version{} }
func (m *Version) String() string { return proto.CompactTextString(m) }
func (*Version) ProtoMessage()    {}
func (*Version) Descriptor() ([]byte, []int) {
	return fileDescriptor_620bb2c2a9d5cb2c, []int{0}
}

func (m *Version) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Version.Unmarshal(m, b)
}
func (m *Version) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Version.Marshal(b, m, deterministic)
}
func (m *Version) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Version.Merge(m, src)
}
func (m *Version) XXX_Size() int {
	return xxx_messageInfo_Version.Size(m)
}
func (m *Version) XXX_DiscardUnknown() {
	xxx_messageInfo_Version.DiscardUnknown(m)
}

var xxx_messageInfo_Version proto.InternalMessageInfo

func (m *Version) GetMajor() uint32 {
	if m != nil {
		return m.Major
	}
	return 0
}

func (m *Version) GetMinor() uint32 {
	if m != nil {
		return m.Minor
	}
	return 0
}

func (m *Version) GetPatch() uint32 {
	if m != nil {
		return m.Patch
	}
	return 0
}

func init() {
	proto.RegisterType((*Version)(nil), "version.Version")
}

func init() { proto.RegisterFile("types/version/version.proto", fileDescriptor_620bb2c2a9d5cb2c) }

var fileDescriptor_620bb2c2a9d5cb2c = []byte{
	// 142 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2e, 0xa9, 0x2c, 0x48,
	0x2d, 0xd6, 0x2f, 0x4b, 0x2d, 0x2a, 0xce, 0xcc, 0xcf, 0x83, 0xd1, 0x7a, 0x05, 0x45, 0xf9, 0x25,
	0xf9, 0x42, 0xec, 0x50, 0xae, 0x92, 0x37, 0x17, 0x7b, 0x18, 0x84, 0x29, 0x24, 0xc2, 0xc5, 0x9a,
	0x9b, 0x98, 0x95, 0x5f, 0x24, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0x1b, 0x04, 0xe1, 0x80, 0x45, 0x33,
	0xf3, 0xf2, 0x8b, 0x24, 0x98, 0xa0, 0xa2, 0x20, 0x0e, 0x48, 0xb4, 0x20, 0xb1, 0x24, 0x39, 0x43,
	0x82, 0x19, 0x22, 0x0a, 0xe6, 0x38, 0xe9, 0x45, 0xe9, 0xa4, 0x67, 0x96, 0x64, 0x94, 0x26, 0xe9,
	0x25, 0xe7, 0xe7, 0xea, 0xa7, 0x24, 0x67, 0xa4, 0xe6, 0x65, 0xeb, 0xe7, 0x26, 0x56, 0xa5, 0x96,
	0x67, 0x16, 0xa5, 0xea, 0x17, 0x64, 0xa7, 0xeb, 0xa3, 0xb8, 0x29, 0x89, 0x0d, 0xec, 0x18, 0x63,
	0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x15, 0x30, 0x12, 0xf7, 0xab, 0x00, 0x00, 0x00,
}
