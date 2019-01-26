package main

// THIS FILE WAS PRODUCED BY THE MSGP CODE GENERATION TOOL (github.com/dchenk/msgp).
// DO NOT EDIT.

import (
	"github.com/dchenk/msgp/msgp"
)

// EncodeMsg implements msgp.Encoder
func (z RespSiteChangeHome) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "ok"
	err = en.Append(0x81, 0xa2, 0x6f, 0x6b)
	if err != nil {
		return
	}
	err = en.WriteBool(z.Ok)
	if err != nil {
		return
	}
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z RespSiteChangeHome) Msgsize() (s int) {
	s = 1 + 3 + msgp.BoolSize
	return
}

// EncodeMsg implements msgp.Encoder
func (z RespSiteCreate) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "new_id"
	err = en.Append(0x81, 0xa6, 0x6e, 0x65, 0x77, 0x5f, 0x69, 0x64)
	if err != nil {
		return
	}
	err = en.WriteInt64(z.NewId)
	if err != nil {
		return
	}
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z RespSiteCreate) Msgsize() (s int) {
	s = 1 + 7 + msgp.Int64Size
	return
}

// DecodeMsg implements msgp.Decoder
func (z *SiteChangeHome) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch string(field) {
		case "site":
			z.Site, err = dc.ReadInt64()
			if err != nil {
				return
			}
		case "new_id":
			z.HomeNewID, err = dc.ReadUint32()
			if err != nil {
				return
			}
		case "old_slug":
			z.OldHomeSlug, err = dc.ReadString()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z SiteChangeHome) Msgsize() (s int) {
	s = 1 + 5 + msgp.Int64Size + 7 + msgp.Uint32Size + 9 + msgp.StringPrefixSize + len(z.OldHomeSlug)
	return
}

// DecodeMsg implements msgp.Decoder
func (z *SiteCreate) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch string(field) {
		case "domain":
			z.Domain, err = dc.ReadString()
			if err != nil {
				return
			}
		case "name":
			z.Name, err = dc.ReadString()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z SiteCreate) Msgsize() (s int) {
	s = 1 + 7 + msgp.StringPrefixSize + len(z.Domain) + 5 + msgp.StringPrefixSize + len(z.Name)
	return
}

// DecodeMsg implements msgp.Decoder
func (z *SiteDelete) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch string(field) {
		case "site_id":
			z.SiteID, err = dc.ReadInt64()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z SiteDelete) Msgsize() (s int) {
	s = 1 + 8 + msgp.Int64Size
	return
}

// DecodeMsg implements msgp.Decoder
func (z *SiteGetTheme) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch string(field) {
		case "site":
			z.Site, err = dc.ReadInt64()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z SiteGetTheme) Msgsize() (s int) {
	s = 1 + 5 + msgp.Int64Size
	return
}
