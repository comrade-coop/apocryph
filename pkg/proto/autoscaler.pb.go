// SPDX-License-Identifier: GPL-3.0

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v5.28.2
// source: autoscaler.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ConnectClusterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// could get it from within the cluster maybe?
	NodeGateway string   `protobuf:"bytes,1,opt,name=nodeGateway,proto3" json:"nodeGateway,omitempty"`
	Servers     []string `protobuf:"bytes,2,rep,name=servers,proto3" json:"servers,omitempty"`
	Timeout     uint32   `protobuf:"varint,3,opt,name=timeout,proto3" json:"timeout,omitempty"`
}

func (x *ConnectClusterRequest) Reset() {
	*x = ConnectClusterRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_autoscaler_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnectClusterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectClusterRequest) ProtoMessage() {}

func (x *ConnectClusterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_autoscaler_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectClusterRequest.ProtoReflect.Descriptor instead.
func (*ConnectClusterRequest) Descriptor() ([]byte, []int) {
	return file_autoscaler_proto_rawDescGZIP(), []int{0}
}

func (x *ConnectClusterRequest) GetNodeGateway() string {
	if x != nil {
		return x.NodeGateway
	}
	return ""
}

func (x *ConnectClusterRequest) GetServers() []string {
	if x != nil {
		return x.Servers
	}
	return nil
}

func (x *ConnectClusterRequest) GetTimeout() uint32 {
	if x != nil {
		return x.Timeout
	}
	return 0
}

type ConnectClusterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Error   string `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *ConnectClusterResponse) Reset() {
	*x = ConnectClusterResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_autoscaler_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnectClusterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectClusterResponse) ProtoMessage() {}

func (x *ConnectClusterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_autoscaler_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectClusterResponse.ProtoReflect.Descriptor instead.
func (*ConnectClusterResponse) Descriptor() ([]byte, []int) {
	return file_autoscaler_proto_rawDescGZIP(), []int{1}
}

func (x *ConnectClusterResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *ConnectClusterResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type TriggerNodeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PeerID string `protobuf:"bytes,1,opt,name=peerID,proto3" json:"peerID,omitempty"`
}

func (x *TriggerNodeResponse) Reset() {
	*x = TriggerNodeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_autoscaler_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TriggerNodeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TriggerNodeResponse) ProtoMessage() {}

func (x *TriggerNodeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_autoscaler_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TriggerNodeResponse.ProtoReflect.Descriptor instead.
func (*TriggerNodeResponse) Descriptor() ([]byte, []int) {
	return file_autoscaler_proto_rawDescGZIP(), []int{2}
}

func (x *TriggerNodeResponse) GetPeerID() string {
	if x != nil {
		return x.PeerID
	}
	return ""
}

var File_autoscaler_proto protoreflect.FileDescriptor

var file_autoscaler_proto_rawDesc = []byte{
	0x0a, 0x10, 0x61, 0x75, 0x74, 0x6f, 0x73, 0x63, 0x61, 0x6c, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x1c, 0x61, 0x70, 0x6f, 0x63, 0x72, 0x79, 0x70, 0x68, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x76, 0x30, 0x2e, 0x61, 0x75, 0x74, 0x6f, 0x73, 0x63, 0x61, 0x6c, 0x65, 0x72,
	0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6d, 0x0a,
	0x15, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x6e, 0x6f, 0x64, 0x65, 0x47, 0x61,
	0x74, 0x65, 0x77, 0x61, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6e, 0x6f, 0x64,
	0x65, 0x47, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x22, 0x48, 0x0a, 0x16,
	0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x2d, 0x0a, 0x13, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65,
	0x72, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a,
	0x06, 0x70, 0x65, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70,
	0x65, 0x65, 0x72, 0x49, 0x44, 0x32, 0x87, 0x02, 0x0a, 0x11, 0x41, 0x75, 0x74, 0x6f, 0x73, 0x63,
	0x61, 0x6c, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x7b, 0x0a, 0x0e, 0x43,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x12, 0x33, 0x2e,
	0x61, 0x70, 0x6f, 0x63, 0x72, 0x79, 0x70, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76,
	0x30, 0x2e, 0x61, 0x75, 0x74, 0x6f, 0x73, 0x63, 0x61, 0x6c, 0x65, 0x72, 0x2e, 0x43, 0x6f, 0x6e,
	0x6e, 0x65, 0x63, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x34, 0x2e, 0x61, 0x70, 0x6f, 0x63, 0x72, 0x79, 0x70, 0x68, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x30, 0x2e, 0x61, 0x75, 0x74, 0x6f, 0x73, 0x63, 0x61, 0x6c, 0x65,
	0x72, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x75, 0x0a, 0x0b, 0x54, 0x72, 0x69, 0x67,
	0x67, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x33, 0x2e, 0x61, 0x70, 0x6f, 0x63, 0x72, 0x79,
	0x70, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x30, 0x2e, 0x61, 0x75, 0x74, 0x6f,
	0x73, 0x63, 0x61, 0x6c, 0x65, 0x72, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x43, 0x6c,
	0x75, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x31, 0x2e, 0x61,
	0x70, 0x6f, 0x63, 0x72, 0x79, 0x70, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x30,
	0x2e, 0x61, 0x75, 0x74, 0x6f, 0x73, 0x63, 0x61, 0x6c, 0x65, 0x72, 0x2e, 0x54, 0x72, 0x69, 0x67,
	0x67, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42,
	0x2c, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f,
	0x6d, 0x72, 0x61, 0x64, 0x65, 0x2d, 0x63, 0x6f, 0x6f, 0x70, 0x2f, 0x61, 0x70, 0x6f, 0x63, 0x72,
	0x79, 0x70, 0x68, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_autoscaler_proto_rawDescOnce sync.Once
	file_autoscaler_proto_rawDescData = file_autoscaler_proto_rawDesc
)

func file_autoscaler_proto_rawDescGZIP() []byte {
	file_autoscaler_proto_rawDescOnce.Do(func() {
		file_autoscaler_proto_rawDescData = protoimpl.X.CompressGZIP(file_autoscaler_proto_rawDescData)
	})
	return file_autoscaler_proto_rawDescData
}

var file_autoscaler_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_autoscaler_proto_goTypes = []interface{}{
	(*ConnectClusterRequest)(nil),  // 0: apocryph.proto.v0.autoscaler.ConnectClusterRequest
	(*ConnectClusterResponse)(nil), // 1: apocryph.proto.v0.autoscaler.ConnectClusterResponse
	(*TriggerNodeResponse)(nil),    // 2: apocryph.proto.v0.autoscaler.TriggerNodeResponse
}
var file_autoscaler_proto_depIdxs = []int32{
	0, // 0: apocryph.proto.v0.autoscaler.AutoscalerService.ConnectCluster:input_type -> apocryph.proto.v0.autoscaler.ConnectClusterRequest
	0, // 1: apocryph.proto.v0.autoscaler.AutoscalerService.TriggerNode:input_type -> apocryph.proto.v0.autoscaler.ConnectClusterRequest
	1, // 2: apocryph.proto.v0.autoscaler.AutoscalerService.ConnectCluster:output_type -> apocryph.proto.v0.autoscaler.ConnectClusterResponse
	2, // 3: apocryph.proto.v0.autoscaler.AutoscalerService.TriggerNode:output_type -> apocryph.proto.v0.autoscaler.TriggerNodeResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_autoscaler_proto_init() }
func file_autoscaler_proto_init() {
	if File_autoscaler_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_autoscaler_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnectClusterRequest); i {
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
		file_autoscaler_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnectClusterResponse); i {
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
		file_autoscaler_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TriggerNodeResponse); i {
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
			RawDescriptor: file_autoscaler_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_autoscaler_proto_goTypes,
		DependencyIndexes: file_autoscaler_proto_depIdxs,
		MessageInfos:      file_autoscaler_proto_msgTypes,
	}.Build()
	File_autoscaler_proto = out.File
	file_autoscaler_proto_rawDesc = nil
	file_autoscaler_proto_goTypes = nil
	file_autoscaler_proto_depIdxs = nil
}
