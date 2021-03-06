// Code generated by protoc-gen-go. DO NOT EDIT.
// source: service.proto

/*
Package pricing is a generated protocol buffer package.

It is generated from these files:
	service.proto

It has these top-level messages:
	ItemPriceRequest
	ItemPriceResponse
	Item
	ItemPrice
*/
package pricing

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ItemPriceRequest struct {
	RegionId int32 `protobuf:"varint,1,opt,name=RegionId" json:"RegionId,omitempty"`
	ItemId   int32 `protobuf:"varint,2,opt,name=ItemId" json:"ItemId,omitempty"`
}

func (m *ItemPriceRequest) Reset()                    { *m = ItemPriceRequest{} }
func (m *ItemPriceRequest) String() string            { return proto.CompactTextString(m) }
func (*ItemPriceRequest) ProtoMessage()               {}
func (*ItemPriceRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ItemPriceRequest) GetRegionId() int32 {
	if m != nil {
		return m.RegionId
	}
	return 0
}

func (m *ItemPriceRequest) GetItemId() int32 {
	if m != nil {
		return m.ItemId
	}
	return 0
}

type ItemPriceResponse struct {
	Item *Item `protobuf:"bytes,1,opt,name=Item" json:"Item,omitempty"`
}

func (m *ItemPriceResponse) Reset()                    { *m = ItemPriceResponse{} }
func (m *ItemPriceResponse) String() string            { return proto.CompactTextString(m) }
func (*ItemPriceResponse) ProtoMessage()               {}
func (*ItemPriceResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ItemPriceResponse) GetItem() *Item {
	if m != nil {
		return m.Item
	}
	return nil
}

type Item struct {
	ItemId int32      `protobuf:"varint,1,opt,name=ItemId" json:"ItemId,omitempty"`
	Buy    *ItemPrice `protobuf:"bytes,2,opt,name=Buy" json:"Buy,omitempty"`
	Sell   *ItemPrice `protobuf:"bytes,3,opt,name=Sell" json:"Sell,omitempty"`
}

func (m *Item) Reset()                    { *m = Item{} }
func (m *Item) String() string            { return proto.CompactTextString(m) }
func (*Item) ProtoMessage()               {}
func (*Item) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Item) GetItemId() int32 {
	if m != nil {
		return m.ItemId
	}
	return 0
}

func (m *Item) GetBuy() *ItemPrice {
	if m != nil {
		return m.Buy
	}
	return nil
}

func (m *Item) GetSell() *ItemPrice {
	if m != nil {
		return m.Sell
	}
	return nil
}

type ItemPrice struct {
	Min float64 `protobuf:"fixed64,1,opt,name=Min" json:"Min,omitempty"`
	Max float64 `protobuf:"fixed64,2,opt,name=Max" json:"Max,omitempty"`
	Avg float64 `protobuf:"fixed64,3,opt,name=Avg" json:"Avg,omitempty"`
	Vol int64   `protobuf:"varint,4,opt,name=Vol" json:"Vol,omitempty"`
	Ord int64   `protobuf:"varint,5,opt,name=Ord" json:"Ord,omitempty"`
}

func (m *ItemPrice) Reset()                    { *m = ItemPrice{} }
func (m *ItemPrice) String() string            { return proto.CompactTextString(m) }
func (*ItemPrice) ProtoMessage()               {}
func (*ItemPrice) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *ItemPrice) GetMin() float64 {
	if m != nil {
		return m.Min
	}
	return 0
}

func (m *ItemPrice) GetMax() float64 {
	if m != nil {
		return m.Max
	}
	return 0
}

func (m *ItemPrice) GetAvg() float64 {
	if m != nil {
		return m.Avg
	}
	return 0
}

func (m *ItemPrice) GetVol() int64 {
	if m != nil {
		return m.Vol
	}
	return 0
}

func (m *ItemPrice) GetOrd() int64 {
	if m != nil {
		return m.Ord
	}
	return 0
}

func init() {
	proto.RegisterType((*ItemPriceRequest)(nil), "pricing.ItemPriceRequest")
	proto.RegisterType((*ItemPriceResponse)(nil), "pricing.ItemPriceResponse")
	proto.RegisterType((*Item)(nil), "pricing.Item")
	proto.RegisterType((*ItemPrice)(nil), "pricing.ItemPrice")
}

func init() { proto.RegisterFile("service.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 265 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0x41, 0x4b, 0xc3, 0x40,
	0x10, 0x85, 0x59, 0x93, 0x46, 0x9d, 0x5a, 0xa8, 0x7b, 0x90, 0x98, 0x53, 0x0d, 0x22, 0x3d, 0xe5,
	0x10, 0xc1, 0xbb, 0x82, 0x4a, 0x0e, 0x12, 0x59, 0xc1, 0x7b, 0x4d, 0x86, 0xb0, 0x10, 0xb3, 0x71,
	0x37, 0x2d, 0xf5, 0xdf, 0xcb, 0xcc, 0x96, 0x1a, 0xa1, 0xde, 0x66, 0xbf, 0x37, 0xbc, 0x2f, 0x4c,
	0x60, 0xe6, 0xd0, 0x6e, 0x74, 0x85, 0x59, 0x6f, 0xcd, 0x60, 0xe4, 0x71, 0x6f, 0x75, 0xa5, 0xbb,
	0x26, 0x7d, 0x82, 0x79, 0x31, 0xe0, 0xe7, 0xab, 0xd5, 0x15, 0x2a, 0xfc, 0x5a, 0xa3, 0x1b, 0x64,
	0x02, 0x27, 0x0a, 0x1b, 0x6d, 0xba, 0xa2, 0x8e, 0xc5, 0x42, 0x2c, 0x27, 0x6a, 0xff, 0x96, 0x17,
	0x10, 0xd1, 0x7e, 0x51, 0xc7, 0x47, 0x9c, 0xec, 0x5e, 0xe9, 0x1d, 0x9c, 0x8f, 0x7a, 0x5c, 0x6f,
	0x3a, 0x87, 0xf2, 0x0a, 0x42, 0x82, 0x5c, 0x32, 0xcd, 0x67, 0xd9, 0x4e, 0x9a, 0x11, 0x54, 0x1c,
	0xa5, 0xad, 0x5f, 0x19, 0xf5, 0x8a, 0x71, 0xaf, 0xbc, 0x86, 0xe0, 0x61, 0xfd, 0xcd, 0xb2, 0x69,
	0x2e, 0xff, 0x34, 0x78, 0x17, 0xc5, 0xf2, 0x06, 0xc2, 0x37, 0x6c, 0xdb, 0x38, 0xf8, 0x77, 0x8d,
	0xf3, 0x14, 0xe1, 0x74, 0x8f, 0xe4, 0x1c, 0x82, 0x17, 0xdd, 0xb1, 0x4f, 0x28, 0x1a, 0x99, 0xac,
	0xb6, 0x2c, 0x23, 0xb2, 0xda, 0x12, 0xb9, 0xdf, 0x34, 0xdc, 0x2b, 0x14, 0x8d, 0x44, 0xde, 0x4d,
	0x1b, 0x87, 0x0b, 0xb1, 0x0c, 0x14, 0x8d, 0x44, 0x4a, 0x5b, 0xc7, 0x13, 0x4f, 0x4a, 0x5b, 0xe7,
	0x25, 0x44, 0xac, 0x70, 0xf2, 0x11, 0xce, 0x9e, 0x71, 0xf8, 0x75, 0x5e, 0x1e, 0xf8, 0x34, 0x7f,
	0xf5, 0x24, 0x39, 0x14, 0xf9, 0x43, 0x7e, 0x44, 0xfc, 0xd7, 0x6e, 0x7f, 0x02, 0x00, 0x00, 0xff,
	0xff, 0x00, 0x58, 0xb2, 0x2e, 0xc6, 0x01, 0x00, 0x00,
}
