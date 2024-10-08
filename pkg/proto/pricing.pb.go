// SPDX-License-Identifier: GPL-3.0

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v5.28.2
// source: pricing.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PricingTables struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tables []*PricingTable `protobuf:"bytes,1,rep,name=tables,proto3" json:"tables,omitempty"`
}

func (x *PricingTables) Reset() {
	*x = PricingTables{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pricing_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PricingTables) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PricingTables) ProtoMessage() {}

func (x *PricingTables) ProtoReflect() protoreflect.Message {
	mi := &file_pricing_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PricingTables.ProtoReflect.Descriptor instead.
func (*PricingTables) Descriptor() ([]byte, []int) {
	return file_pricing_proto_rawDescGZIP(), []int{0}
}

func (x *PricingTables) GetTables() []*PricingTable {
	if x != nil {
		return x.Tables
	}
	return nil
}

type PricingTable struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Resources              []*PricedResource `protobuf:"bytes,1,rep,name=resources,proto3" json:"resources,omitempty"`
	PaymentContractAddress []byte            `protobuf:"bytes,2,opt,name=paymentContractAddress,proto3" json:"paymentContractAddress,omitempty"`
}

func (x *PricingTable) Reset() {
	*x = PricingTable{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pricing_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PricingTable) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PricingTable) ProtoMessage() {}

func (x *PricingTable) ProtoReflect() protoreflect.Message {
	mi := &file_pricing_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PricingTable.ProtoReflect.Descriptor instead.
func (*PricingTable) Descriptor() ([]byte, []int) {
	return file_pricing_proto_rawDescGZIP(), []int{1}
}

func (x *PricingTable) GetResources() []*PricedResource {
	if x != nil {
		return x.Resources
	}
	return nil
}

func (x *PricingTable) GetPaymentContractAddress() []byte {
	if x != nil {
		return x.PaymentContractAddress
	}
	return nil
}

type PricedResource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Resource            string `protobuf:"bytes,1,opt,name=resource,proto3" json:"resource,omitempty"`
	PriceForReservation uint64 `protobuf:"varint,2,opt,name=priceForReservation,proto3" json:"priceForReservation,omitempty"`
	PriceForUsage       uint64 `protobuf:"varint,3,opt,name=priceForUsage,proto3" json:"priceForUsage,omitempty"` // uint64 priceForLimit = 4;
}

func (x *PricedResource) Reset() {
	*x = PricedResource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pricing_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PricedResource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PricedResource) ProtoMessage() {}

func (x *PricedResource) ProtoReflect() protoreflect.Message {
	mi := &file_pricing_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PricedResource.ProtoReflect.Descriptor instead.
func (*PricedResource) Descriptor() ([]byte, []int) {
	return file_pricing_proto_rawDescGZIP(), []int{2}
}

func (x *PricedResource) GetResource() string {
	if x != nil {
		return x.Resource
	}
	return ""
}

func (x *PricedResource) GetPriceForReservation() uint64 {
	if x != nil {
		return x.PriceForReservation
	}
	return 0
}

func (x *PricedResource) GetPriceForUsage() uint64 {
	if x != nil {
		return x.PriceForUsage
	}
	return 0
}

var File_pricing_proto protoreflect.FileDescriptor

var file_pricing_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x70, 0x72, 0x69, 0x63, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x19, 0x61, 0x70, 0x6f, 0x63, 0x72, 0x79, 0x70, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x76, 0x30, 0x2e, 0x70, 0x72, 0x69, 0x63, 0x69, 0x6e, 0x67, 0x22, 0x50, 0x0a, 0x0d, 0x50, 0x72,
	0x69, 0x63, 0x69, 0x6e, 0x67, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x12, 0x3f, 0x0a, 0x06, 0x74,
	0x61, 0x62, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x61, 0x70,
	0x6f, 0x63, 0x72, 0x79, 0x70, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x30, 0x2e,
	0x70, 0x72, 0x69, 0x63, 0x69, 0x6e, 0x67, 0x2e, 0x50, 0x72, 0x69, 0x63, 0x69, 0x6e, 0x67, 0x54,
	0x61, 0x62, 0x6c, 0x65, 0x52, 0x06, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x22, 0x8f, 0x01, 0x0a,
	0x0c, 0x50, 0x72, 0x69, 0x63, 0x69, 0x6e, 0x67, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x47, 0x0a,
	0x09, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x29, 0x2e, 0x61, 0x70, 0x6f, 0x63, 0x72, 0x79, 0x70, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x76, 0x30, 0x2e, 0x70, 0x72, 0x69, 0x63, 0x69, 0x6e, 0x67, 0x2e, 0x50, 0x72, 0x69,
	0x63, 0x65, 0x64, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x09, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x12, 0x36, 0x0a, 0x16, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e,
	0x74, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x16, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x43,
	0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x84,
	0x01, 0x0a, 0x0e, 0x50, 0x72, 0x69, 0x63, 0x65, 0x64, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x30, 0x0a,
	0x13, 0x70, 0x72, 0x69, 0x63, 0x65, 0x46, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x65, 0x72, 0x76, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x13, 0x70, 0x72, 0x69, 0x63,
	0x65, 0x46, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x24, 0x0a, 0x0d, 0x70, 0x72, 0x69, 0x63, 0x65, 0x46, 0x6f, 0x72, 0x55, 0x73, 0x61, 0x67, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0d, 0x70, 0x72, 0x69, 0x63, 0x65, 0x46, 0x6f, 0x72,
	0x55, 0x73, 0x61, 0x67, 0x65, 0x42, 0x2c, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x6d, 0x72, 0x61, 0x64, 0x65, 0x2d, 0x63, 0x6f, 0x6f, 0x70,
	0x2f, 0x61, 0x70, 0x6f, 0x63, 0x72, 0x79, 0x70, 0x68, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pricing_proto_rawDescOnce sync.Once
	file_pricing_proto_rawDescData = file_pricing_proto_rawDesc
)

func file_pricing_proto_rawDescGZIP() []byte {
	file_pricing_proto_rawDescOnce.Do(func() {
		file_pricing_proto_rawDescData = protoimpl.X.CompressGZIP(file_pricing_proto_rawDescData)
	})
	return file_pricing_proto_rawDescData
}

var file_pricing_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_pricing_proto_goTypes = []interface{}{
	(*PricingTables)(nil),  // 0: apocryph.proto.v0.pricing.PricingTables
	(*PricingTable)(nil),   // 1: apocryph.proto.v0.pricing.PricingTable
	(*PricedResource)(nil), // 2: apocryph.proto.v0.pricing.PricedResource
}
var file_pricing_proto_depIdxs = []int32{
	1, // 0: apocryph.proto.v0.pricing.PricingTables.tables:type_name -> apocryph.proto.v0.pricing.PricingTable
	2, // 1: apocryph.proto.v0.pricing.PricingTable.resources:type_name -> apocryph.proto.v0.pricing.PricedResource
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_pricing_proto_init() }
func file_pricing_proto_init() {
	if File_pricing_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pricing_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PricingTables); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pricing_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PricingTable); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pricing_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PricedResource); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pricing_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pricing_proto_goTypes,
		DependencyIndexes: file_pricing_proto_depIdxs,
		MessageInfos:      file_pricing_proto_msgTypes,
	}.Build()
	File_pricing_proto = out.File
	file_pricing_proto_rawDesc = nil
	file_pricing_proto_goTypes = nil
	file_pricing_proto_depIdxs = nil
}
