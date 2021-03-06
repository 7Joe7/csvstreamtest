// Code generated by protoc-gen-go.
// source: common/message/client.proto
// DO NOT EDIT!

/*
Package model is a generated protocol buffer package.

It is generated from these files:
	common/message/client.proto
	common/message/importreport.proto

It has these top-level messages:
	Client
	ImportReport
*/
package model

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

type Client struct {
	// @inject_tag: csv:"id"
	Id int32 `protobuf:"varint,1,opt,name=id" json:"id,omitempty" csv:"id"`
	// @inject_tag: csv:"name"
	Name string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty" csv:"name"`
	// @inject_tag: csv:"email"
	Email string `protobuf:"bytes,3,opt,name=email" json:"email,omitempty" csv:"email"`
	// @inject_tag: csv:"mobile_number"
	MobileNumber string `protobuf:"bytes,4,opt,name=mobileNumber" json:"mobileNumber,omitempty" csv:"mobile_number"`
}

func (m *Client) Reset()                    { *m = Client{} }
func (m *Client) String() string            { return proto.CompactTextString(m) }
func (*Client) ProtoMessage()               {}
func (*Client) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Client) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Client) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Client) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *Client) GetMobileNumber() string {
	if m != nil {
		return m.MobileNumber
	}
	return ""
}

func init() {
	proto.RegisterType((*Client)(nil), "model.Client")
}

func init() { proto.RegisterFile("common/message/client.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 177 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0xce, 0xb1, 0xab, 0xc2, 0x30,
	0x10, 0xc7, 0x71, 0xda, 0xd7, 0x16, 0x5e, 0x10, 0x87, 0xe0, 0x10, 0x70, 0x29, 0x9d, 0x0a, 0x62,
	0x33, 0x38, 0x74, 0xd7, 0xdd, 0xa1, 0xa3, 0x5b, 0x92, 0x9e, 0x35, 0x92, 0xeb, 0x49, 0x92, 0xfa,
	0xf7, 0x8b, 0x11, 0x07, 0xb7, 0xbb, 0xef, 0x6f, 0xf9, 0xb0, 0xad, 0x21, 0x44, 0x9a, 0x25, 0x42,
	0x08, 0x6a, 0x02, 0x69, 0x9c, 0x85, 0x39, 0x76, 0x0f, 0x4f, 0x91, 0x78, 0x89, 0x34, 0x82, 0x6b,
	0xae, 0xac, 0x3a, 0xa5, 0xcc, 0xd7, 0x2c, 0xb7, 0xa3, 0xc8, 0xea, 0xac, 0x2d, 0x87, 0xdc, 0x8e,
	0x9c, 0xb3, 0x62, 0x56, 0x08, 0x22, 0xaf, 0xb3, 0xf6, 0x7f, 0x48, 0x37, 0xdf, 0xb0, 0x12, 0x50,
	0x59, 0x27, 0xfe, 0x52, 0xfc, 0x3c, 0xbc, 0x61, 0x2b, 0x24, 0x6d, 0x1d, 0x9c, 0x17, 0xd4, 0xe0,
	0x45, 0x91, 0xc6, 0x9f, 0x76, 0xdc, 0x5f, 0x76, 0x93, 0x8d, 0xb7, 0x45, 0x77, 0x86, 0x50, 0xf6,
	0x77, 0x82, 0x5e, 0x9a, 0xf0, 0x0c, 0xd1, 0x83, 0xc2, 0x08, 0x21, 0xca, 0x2f, 0xf6, 0xcd, 0xd2,
	0x55, 0x42, 0x1e, 0x5e, 0x01, 0x00, 0x00, 0xff, 0xff, 0x9b, 0x30, 0x5f, 0xf3, 0xc3, 0x00, 0x00,
	0x00,
}
