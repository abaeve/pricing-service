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
	Min float32 `protobuf:"fixed32,1,opt,name=Min" json:"Min,omitempty"`
	Max float32 `protobuf:"fixed32,2,opt,name=Max" json:"Max,omitempty"`
	Avg float32 `protobuf:"fixed32,3,opt,name=Avg" json:"Avg,omitempty"`
	Vol float32 `protobuf:"fixed32,4,opt,name=Vol" json:"Vol,omitempty"`
	Ord int64   `protobuf:"varint,5,opt,name=Ord" json:"Ord,omitempty"`
}

func (m *ItemPrice) Reset()                    { *m = ItemPrice{} }
func (m *ItemPrice) String() string            { return proto.CompactTextString(m) }
func (*ItemPrice) ProtoMessage()               {}
func (*ItemPrice) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *ItemPrice) GetMin() float32 {
	if m != nil {
		return m.Min
	}
	return 0
}

func (m *ItemPrice) GetMax() float32 {
	if m != nil {
		return m.Max
	}
	return 0
}

func (m *ItemPrice) GetAvg() float32 {
	if m != nil {
		return m.Avg
	}
	return 0
}

func (m *ItemPrice) GetVol() float32 {
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
	// 266 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0x41, 0x4b, 0xc3, 0x40,
	0x10, 0x85, 0x49, 0x93, 0x46, 0x9d, 0x5a, 0xa8, 0x7b, 0x90, 0x98, 0x53, 0x0d, 0x22, 0x3d, 0xe5,
	0x10, 0xc1, 0xbb, 0x82, 0x4a, 0x0e, 0x52, 0x59, 0xa1, 0xf7, 0x9a, 0x0c, 0x61, 0x61, 0xcd, 0xc6,
	0xdd, 0xb4, 0xd4, 0x7f, 0x2f, 0x33, 0x5b, 0x6a, 0x84, 0x7a, 0x9b, 0xfd, 0xde, 0xf0, 0xbe, 0x30,
	0x81, 0xa9, 0x43, 0xbb, 0x55, 0x15, 0xe6, 0x9d, 0x35, 0xbd, 0x11, 0x27, 0x9d, 0x55, 0x95, 0x6a,
	0x9b, 0xec, 0x19, 0x66, 0x65, 0x8f, 0x9f, 0x6f, 0x56, 0x55, 0x28, 0xf1, 0x6b, 0x83, 0xae, 0x17,
	0x29, 0x9c, 0x4a, 0x6c, 0x94, 0x69, 0xcb, 0x3a, 0x09, 0xe6, 0xc1, 0x62, 0x2c, 0x0f, 0x6f, 0x71,
	0x09, 0x31, 0xed, 0x97, 0x75, 0x32, 0xe2, 0x64, 0xff, 0xca, 0xee, 0xe1, 0x62, 0xd0, 0xe3, 0x3a,
	0xd3, 0x3a, 0x14, 0xd7, 0x10, 0x11, 0xe4, 0x92, 0x49, 0x31, 0xcd, 0xf7, 0xd2, 0x9c, 0xa0, 0xe4,
	0x28, 0xd3, 0x7e, 0x65, 0xd0, 0x1b, 0x0c, 0x7b, 0xc5, 0x0d, 0x84, 0x8f, 0x9b, 0x6f, 0x96, 0x4d,
	0x0a, 0xf1, 0xa7, 0xc1, 0xbb, 0x28, 0x16, 0xb7, 0x10, 0xbd, 0xa3, 0xd6, 0x49, 0xf8, 0xef, 0x1a,
	0xe7, 0x19, 0xc2, 0xd9, 0x01, 0x89, 0x19, 0x84, 0xaf, 0xaa, 0x65, 0xdf, 0x48, 0xd2, 0xc8, 0x64,
	0xbd, 0x63, 0x19, 0x91, 0xf5, 0x8e, 0xc8, 0xc3, 0xb6, 0xe1, 0xde, 0x91, 0xa4, 0x91, 0xc8, 0xca,
	0xe8, 0x24, 0xf2, 0x64, 0x65, 0x34, 0x91, 0xa5, 0xad, 0x93, 0xf1, 0x3c, 0x58, 0x84, 0x92, 0xc6,
	0x62, 0x09, 0x31, 0x2b, 0x9c, 0x78, 0x82, 0xf3, 0x17, 0xec, 0x7f, 0x9d, 0x57, 0x47, 0x3e, 0xcd,
	0x5f, 0x3d, 0x4d, 0x8f, 0x45, 0xfe, 0x90, 0x1f, 0x31, 0xff, 0xb5, 0xbb, 0x9f, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x32, 0xec, 0xfc, 0xa2, 0xc6, 0x01, 0x00, 0x00,
}