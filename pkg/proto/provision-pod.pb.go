// SPDX-License-Identifier: GPL-3.0

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v5.28.2
// source: provision-pod.proto

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

type ProvisionPodRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pod     *Pod            `protobuf:"bytes,1,opt,name=pod,proto3" json:"pod,omitempty"`
	Payment *PaymentChannel `protobuf:"bytes,3,opt,name=payment,proto3" json:"payment,omitempty"`
}

func (x *ProvisionPodRequest) Reset() {
	*x = ProvisionPodRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_provision_pod_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProvisionPodRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProvisionPodRequest) ProtoMessage() {}

func (x *ProvisionPodRequest) ProtoReflect() protoreflect.Message {
	mi := &file_provision_pod_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProvisionPodRequest.ProtoReflect.Descriptor instead.
func (*ProvisionPodRequest) Descriptor() ([]byte, []int) {
	return file_provision_pod_proto_rawDescGZIP(), []int{0}
}

func (x *ProvisionPodRequest) GetPod() *Pod {
	if x != nil {
		return x.Pod
	}
	return nil
}

func (x *ProvisionPodRequest) GetPayment() *PaymentChannel {
	if x != nil {
		return x.Payment
	}
	return nil
}

type DeletePodRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeletePodRequest) Reset() {
	*x = DeletePodRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_provision_pod_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeletePodRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeletePodRequest) ProtoMessage() {}

func (x *DeletePodRequest) ProtoReflect() protoreflect.Message {
	mi := &file_provision_pod_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeletePodRequest.ProtoReflect.Descriptor instead.
func (*DeletePodRequest) Descriptor() ([]byte, []int) {
	return file_provision_pod_proto_rawDescGZIP(), []int{1}
}

type DeletePodResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Error   string `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *DeletePodResponse) Reset() {
	*x = DeletePodResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_provision_pod_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeletePodResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeletePodResponse) ProtoMessage() {}

func (x *DeletePodResponse) ProtoReflect() protoreflect.Message {
	mi := &file_provision_pod_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeletePodResponse.ProtoReflect.Descriptor instead.
func (*DeletePodResponse) Descriptor() ([]byte, []int) {
	return file_provision_pod_proto_rawDescGZIP(), []int{2}
}

func (x *DeletePodResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *DeletePodResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type UpdatePodRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pod     *Pod            `protobuf:"bytes,1,opt,name=pod,proto3" json:"pod,omitempty"`
	Payment *PaymentChannel `protobuf:"bytes,2,opt,name=payment,proto3" json:"payment,omitempty"`
}

func (x *UpdatePodRequest) Reset() {
	*x = UpdatePodRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_provision_pod_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdatePodRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdatePodRequest) ProtoMessage() {}

func (x *UpdatePodRequest) ProtoReflect() protoreflect.Message {
	mi := &file_provision_pod_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdatePodRequest.ProtoReflect.Descriptor instead.
func (*UpdatePodRequest) Descriptor() ([]byte, []int) {
	return file_provision_pod_proto_rawDescGZIP(), []int{3}
}

func (x *UpdatePodRequest) GetPod() *Pod {
	if x != nil {
		return x.Pod
	}
	return nil
}

func (x *UpdatePodRequest) GetPayment() *PaymentChannel {
	if x != nil {
		return x.Payment
	}
	return nil
}

type UpdatePodResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Error   string `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *UpdatePodResponse) Reset() {
	*x = UpdatePodResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_provision_pod_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdatePodResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdatePodResponse) ProtoMessage() {}

func (x *UpdatePodResponse) ProtoReflect() protoreflect.Message {
	mi := &file_provision_pod_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdatePodResponse.ProtoReflect.Descriptor instead.
func (*UpdatePodResponse) Descriptor() ([]byte, []int) {
	return file_provision_pod_proto_rawDescGZIP(), []int{4}
}

func (x *UpdatePodResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *UpdatePodResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type PaymentChannel struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChainID          []byte `protobuf:"bytes,1,opt,name=chainID,proto3" json:"chainID,omitempty"`
	ContractAddress  []byte `protobuf:"bytes,2,opt,name=contractAddress,proto3" json:"contractAddress,omitempty"`
	PublisherAddress []byte `protobuf:"bytes,3,opt,name=publisherAddress,proto3" json:"publisherAddress,omitempty"`
	ProviderAddress  []byte `protobuf:"bytes,4,opt,name=providerAddress,proto3" json:"providerAddress,omitempty"`
	PodID            []byte `protobuf:"bytes,5,opt,name=podID,proto3" json:"podID,omitempty"`
}

func (x *PaymentChannel) Reset() {
	*x = PaymentChannel{}
	if protoimpl.UnsafeEnabled {
		mi := &file_provision_pod_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PaymentChannel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PaymentChannel) ProtoMessage() {}

func (x *PaymentChannel) ProtoReflect() protoreflect.Message {
	mi := &file_provision_pod_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PaymentChannel.ProtoReflect.Descriptor instead.
func (*PaymentChannel) Descriptor() ([]byte, []int) {
	return file_provision_pod_proto_rawDescGZIP(), []int{5}
}

func (x *PaymentChannel) GetChainID() []byte {
	if x != nil {
		return x.ChainID
	}
	return nil
}

func (x *PaymentChannel) GetContractAddress() []byte {
	if x != nil {
		return x.ContractAddress
	}
	return nil
}

func (x *PaymentChannel) GetPublisherAddress() []byte {
	if x != nil {
		return x.PublisherAddress
	}
	return nil
}

func (x *PaymentChannel) GetProviderAddress() []byte {
	if x != nil {
		return x.ProviderAddress
	}
	return nil
}

func (x *PaymentChannel) GetPodID() []byte {
	if x != nil {
		return x.PodID
	}
	return nil
}

type ProvisionPodResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error            string                                  `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	Addresses        []*ProvisionPodResponse_ExposedHostPort `protobuf:"bytes,2,rep,name=addresses,proto3" json:"addresses,omitempty"`
	Namespace        string                                  `protobuf:"bytes,3,opt,name=namespace,proto3" json:"namespace,omitempty"`
	VerificationHost string                                  `protobuf:"bytes,4,opt,name=verificationHost,proto3" json:"verificationHost,omitempty"`
}

func (x *ProvisionPodResponse) Reset() {
	*x = ProvisionPodResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_provision_pod_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProvisionPodResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProvisionPodResponse) ProtoMessage() {}

func (x *ProvisionPodResponse) ProtoReflect() protoreflect.Message {
	mi := &file_provision_pod_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProvisionPodResponse.ProtoReflect.Descriptor instead.
func (*ProvisionPodResponse) Descriptor() ([]byte, []int) {
	return file_provision_pod_proto_rawDescGZIP(), []int{6}
}

func (x *ProvisionPodResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

func (x *ProvisionPodResponse) GetAddresses() []*ProvisionPodResponse_ExposedHostPort {
	if x != nil {
		return x.Addresses
	}
	return nil
}

func (x *ProvisionPodResponse) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *ProvisionPodResponse) GetVerificationHost() string {
	if x != nil {
		return x.VerificationHost
	}
	return ""
}

type PodLogRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ContainerName string `protobuf:"bytes,1,opt,name=ContainerName,proto3" json:"ContainerName,omitempty"`
}

func (x *PodLogRequest) Reset() {
	*x = PodLogRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_provision_pod_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PodLogRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PodLogRequest) ProtoMessage() {}

func (x *PodLogRequest) ProtoReflect() protoreflect.Message {
	mi := &file_provision_pod_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PodLogRequest.ProtoReflect.Descriptor instead.
func (*PodLogRequest) Descriptor() ([]byte, []int) {
	return file_provision_pod_proto_rawDescGZIP(), []int{7}
}

func (x *PodLogRequest) GetContainerName() string {
	if x != nil {
		return x.ContainerName
	}
	return ""
}

type PodLogResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LogEntry *LogEntry `protobuf:"bytes,1,opt,name=logEntry,proto3" json:"logEntry,omitempty"`
}

func (x *PodLogResponse) Reset() {
	*x = PodLogResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_provision_pod_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PodLogResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PodLogResponse) ProtoMessage() {}

func (x *PodLogResponse) ProtoReflect() protoreflect.Message {
	mi := &file_provision_pod_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PodLogResponse.ProtoReflect.Descriptor instead.
func (*PodLogResponse) Descriptor() ([]byte, []int) {
	return file_provision_pod_proto_rawDescGZIP(), []int{8}
}

func (x *PodLogResponse) GetLogEntry() *LogEntry {
	if x != nil {
		return x.LogEntry
	}
	return nil
}

type LogEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NanosecondsUnixEpoch uint64 `protobuf:"varint,1,opt,name=NanosecondsUnixEpoch,proto3" json:"NanosecondsUnixEpoch,omitempty"`
	Line                 string `protobuf:"bytes,2,opt,name=line,proto3" json:"line,omitempty"`
}

func (x *LogEntry) Reset() {
	*x = LogEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_provision_pod_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogEntry) ProtoMessage() {}

func (x *LogEntry) ProtoReflect() protoreflect.Message {
	mi := &file_provision_pod_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogEntry.ProtoReflect.Descriptor instead.
func (*LogEntry) Descriptor() ([]byte, []int) {
	return file_provision_pod_proto_rawDescGZIP(), []int{9}
}

func (x *LogEntry) GetNanosecondsUnixEpoch() uint64 {
	if x != nil {
		return x.NanosecondsUnixEpoch
	}
	return 0
}

func (x *LogEntry) GetLine() string {
	if x != nil {
		return x.Line
	}
	return ""
}

type ProvisionPodResponse_ExposedHostPort struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Multiaddr     string `protobuf:"bytes,1,opt,name=multiaddr,proto3" json:"multiaddr,omitempty"`
	ContainerName string `protobuf:"bytes,2,opt,name=containerName,proto3" json:"containerName,omitempty"`
	ContainerPort uint64 `protobuf:"varint,3,opt,name=containerPort,proto3" json:"containerPort,omitempty"`
}

func (x *ProvisionPodResponse_ExposedHostPort) Reset() {
	*x = ProvisionPodResponse_ExposedHostPort{}
	if protoimpl.UnsafeEnabled {
		mi := &file_provision_pod_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProvisionPodResponse_ExposedHostPort) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProvisionPodResponse_ExposedHostPort) ProtoMessage() {}

func (x *ProvisionPodResponse_ExposedHostPort) ProtoReflect() protoreflect.Message {
	mi := &file_provision_pod_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProvisionPodResponse_ExposedHostPort.ProtoReflect.Descriptor instead.
func (*ProvisionPodResponse_ExposedHostPort) Descriptor() ([]byte, []int) {
	return file_provision_pod_proto_rawDescGZIP(), []int{6, 0}
}

func (x *ProvisionPodResponse_ExposedHostPort) GetMultiaddr() string {
	if x != nil {
		return x.Multiaddr
	}
	return ""
}

func (x *ProvisionPodResponse_ExposedHostPort) GetContainerName() string {
	if x != nil {
		return x.ContainerName
	}
	return ""
}

func (x *ProvisionPodResponse_ExposedHostPort) GetContainerPort() uint64 {
	if x != nil {
		return x.ContainerPort
	}
	return 0
}

var File_provision_pod_proto protoreflect.FileDescriptor

var file_provision_pod_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x2d, 0x70, 0x6f, 0x64, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1e, 0x61, 0x70, 0x6f, 0x63, 0x72, 0x79, 0x70, 0x68, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x30, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x73, 0x69,
	0x6f, 0x6e, 0x50, 0x6f, 0x64, 0x1a, 0x09, 0x70, 0x6f, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x8d, 0x01, 0x0a, 0x13, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x50, 0x6f,
	0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2c, 0x0a, 0x03, 0x70, 0x6f, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x61, 0x70, 0x6f, 0x63, 0x72, 0x79, 0x70, 0x68,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x30, 0x2e, 0x70, 0x6f, 0x64, 0x2e, 0x50, 0x6f,
	0x64, 0x52, 0x03, 0x70, 0x6f, 0x64, 0x12, 0x48, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x61, 0x70, 0x6f, 0x63, 0x72, 0x79,
	0x70, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x30, 0x2e, 0x70, 0x72, 0x6f, 0x76,
	0x69, 0x73, 0x69, 0x6f, 0x6e, 0x50, 0x6f, 0x64, 0x2e, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x22, 0x12, 0x0a, 0x10, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x50, 0x6f, 0x64, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x22, 0x43, 0x0a, 0x11, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x50, 0x6f,
	0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x8a, 0x01, 0x0a, 0x10, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x50, 0x6f, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2c,
	0x0a, 0x03, 0x70, 0x6f, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x61, 0x70,
	0x6f, 0x63, 0x72, 0x79, 0x70, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x30, 0x2e,
	0x70, 0x6f, 0x64, 0x2e, 0x50, 0x6f, 0x64, 0x52, 0x03, 0x70, 0x6f, 0x64, 0x12, 0x48, 0x0a, 0x07,
	0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2e, 0x2e,
	0x61, 0x70, 0x6f, 0x63, 0x72, 0x79, 0x70, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76,
	0x30, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x50, 0x6f, 0x64, 0x2e, 0x50,
	0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x07, 0x70,
	0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x22, 0x43, 0x0a, 0x11, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x50, 0x6f, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73,
	0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0xc0, 0x01, 0x0a, 0x0e,
	0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x18,
	0x0a, 0x07, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x07, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x44, 0x12, 0x28, 0x0a, 0x0f, 0x63, 0x6f, 0x6e, 0x74,
	0x72, 0x61, 0x63, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x0f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x12, 0x2a, 0x0a, 0x10, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x72, 0x41,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x10, 0x70, 0x75,
	0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x28,
	0x0a, 0x0f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65,
	0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x6f, 0x64, 0x49,
	0x44, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x70, 0x6f, 0x64, 0x49, 0x44, 0x22, 0xd7,
	0x02, 0x0a, 0x14, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x50, 0x6f, 0x64, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x62, 0x0a,
	0x09, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x44, 0x2e, 0x61, 0x70, 0x6f, 0x63, 0x72, 0x79, 0x70, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x76, 0x30, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x50, 0x6f,
	0x64, 0x2e, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x50, 0x6f, 0x64, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x45, 0x78, 0x70, 0x6f, 0x73, 0x65, 0x64, 0x48, 0x6f,
	0x73, 0x74, 0x50, 0x6f, 0x72, 0x74, 0x52, 0x09, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65,
	0x73, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12,
	0x2a, 0x0a, 0x10, 0x76, 0x65, 0x72, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x48,
	0x6f, 0x73, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x76, 0x65, 0x72, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x6f, 0x73, 0x74, 0x1a, 0x7b, 0x0a, 0x0f, 0x45,
	0x78, 0x70, 0x6f, 0x73, 0x65, 0x64, 0x48, 0x6f, 0x73, 0x74, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x1c,
	0x0a, 0x09, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x61, 0x64, 0x64, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x61, 0x64, 0x64, 0x72, 0x12, 0x24, 0x0a, 0x0d,
	0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x24, 0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x50,
	0x6f, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x61,
	0x69, 0x6e, 0x65, 0x72, 0x50, 0x6f, 0x72, 0x74, 0x22, 0x35, 0x0a, 0x0d, 0x50, 0x6f, 0x64, 0x4c,
	0x6f, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x24, 0x0a, 0x0d, 0x43, 0x6f, 0x6e,
	0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0d, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x22,
	0x56, 0x0a, 0x0e, 0x50, 0x6f, 0x64, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x44, 0x0a, 0x08, 0x6c, 0x6f, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x61, 0x70, 0x6f, 0x63, 0x72, 0x79, 0x70, 0x68, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x30, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x73, 0x69, 0x6f,
	0x6e, 0x50, 0x6f, 0x64, 0x2e, 0x4c, 0x6f, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x6c,
	0x6f, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x22, 0x52, 0x0a, 0x08, 0x4c, 0x6f, 0x67, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x32, 0x0a, 0x14, 0x4e, 0x61, 0x6e, 0x6f, 0x73, 0x65, 0x63, 0x6f, 0x6e,
	0x64, 0x73, 0x55, 0x6e, 0x69, 0x78, 0x45, 0x70, 0x6f, 0x63, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x14, 0x4e, 0x61, 0x6e, 0x6f, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x55, 0x6e,
	0x69, 0x78, 0x45, 0x70, 0x6f, 0x63, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x6c, 0x69, 0x6e, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6c, 0x69, 0x6e, 0x65, 0x32, 0xe6, 0x03, 0x0a, 0x13,
	0x50, 0x72, 0x6f, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x50, 0x6f, 0x64, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x79, 0x0a, 0x0c, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e,
	0x50, 0x6f, 0x64, 0x12, 0x33, 0x2e, 0x61, 0x70, 0x6f, 0x63, 0x72, 0x79, 0x70, 0x68, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x30, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x73, 0x69, 0x6f,
	0x6e, 0x50, 0x6f, 0x64, 0x2e, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x50, 0x6f,
	0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x34, 0x2e, 0x61, 0x70, 0x6f, 0x63, 0x72,
	0x79, 0x70, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x30, 0x2e, 0x70, 0x72, 0x6f,
	0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x50, 0x6f, 0x64, 0x2e, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x73,
	0x69, 0x6f, 0x6e, 0x50, 0x6f, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x73,
	0x0a, 0x09, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x6f, 0x64, 0x12, 0x30, 0x2e, 0x61, 0x70,
	0x6f, 0x63, 0x72, 0x79, 0x70, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x30, 0x2e,
	0x70, 0x72, 0x6f, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x50, 0x6f, 0x64, 0x2e, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x50, 0x6f, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x34, 0x2e,
	0x61, 0x70, 0x6f, 0x63, 0x72, 0x79, 0x70, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76,
	0x30, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x50, 0x6f, 0x64, 0x2e, 0x50,
	0x72, 0x6f, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x50, 0x6f, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x70, 0x0a, 0x09, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x50, 0x6f, 0x64,
	0x12, 0x30, 0x2e, 0x61, 0x70, 0x6f, 0x63, 0x72, 0x79, 0x70, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x76, 0x30, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x50, 0x6f,
	0x64, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x50, 0x6f, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x31, 0x2e, 0x61, 0x70, 0x6f, 0x63, 0x72, 0x79, 0x70, 0x68, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x30, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e,
	0x50, 0x6f, 0x64, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x50, 0x6f, 0x64, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x6d, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x50, 0x6f, 0x64, 0x4c,
	0x6f, 0x67, 0x73, 0x12, 0x2d, 0x2e, 0x61, 0x70, 0x6f, 0x63, 0x72, 0x79, 0x70, 0x68, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x30, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x73, 0x69, 0x6f,
	0x6e, 0x50, 0x6f, 0x64, 0x2e, 0x50, 0x6f, 0x64, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x2e, 0x2e, 0x61, 0x70, 0x6f, 0x63, 0x72, 0x79, 0x70, 0x68, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x30, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e,
	0x50, 0x6f, 0x64, 0x2e, 0x50, 0x6f, 0x64, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x30, 0x01, 0x42, 0x2c, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x6d, 0x72, 0x61, 0x64, 0x65, 0x2d, 0x63, 0x6f, 0x6f, 0x70, 0x2f,
	0x61, 0x70, 0x6f, 0x63, 0x72, 0x79, 0x70, 0x68, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_provision_pod_proto_rawDescOnce sync.Once
	file_provision_pod_proto_rawDescData = file_provision_pod_proto_rawDesc
)

func file_provision_pod_proto_rawDescGZIP() []byte {
	file_provision_pod_proto_rawDescOnce.Do(func() {
		file_provision_pod_proto_rawDescData = protoimpl.X.CompressGZIP(file_provision_pod_proto_rawDescData)
	})
	return file_provision_pod_proto_rawDescData
}

var file_provision_pod_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_provision_pod_proto_goTypes = []interface{}{
	(*ProvisionPodRequest)(nil),                  // 0: apocryph.proto.v0.provisionPod.ProvisionPodRequest
	(*DeletePodRequest)(nil),                     // 1: apocryph.proto.v0.provisionPod.DeletePodRequest
	(*DeletePodResponse)(nil),                    // 2: apocryph.proto.v0.provisionPod.DeletePodResponse
	(*UpdatePodRequest)(nil),                     // 3: apocryph.proto.v0.provisionPod.UpdatePodRequest
	(*UpdatePodResponse)(nil),                    // 4: apocryph.proto.v0.provisionPod.UpdatePodResponse
	(*PaymentChannel)(nil),                       // 5: apocryph.proto.v0.provisionPod.PaymentChannel
	(*ProvisionPodResponse)(nil),                 // 6: apocryph.proto.v0.provisionPod.ProvisionPodResponse
	(*PodLogRequest)(nil),                        // 7: apocryph.proto.v0.provisionPod.PodLogRequest
	(*PodLogResponse)(nil),                       // 8: apocryph.proto.v0.provisionPod.PodLogResponse
	(*LogEntry)(nil),                             // 9: apocryph.proto.v0.provisionPod.LogEntry
	(*ProvisionPodResponse_ExposedHostPort)(nil), // 10: apocryph.proto.v0.provisionPod.ProvisionPodResponse.ExposedHostPort
	(*Pod)(nil),                                  // 11: apocryph.proto.v0.pod.Pod
}
var file_provision_pod_proto_depIdxs = []int32{
	11, // 0: apocryph.proto.v0.provisionPod.ProvisionPodRequest.pod:type_name -> apocryph.proto.v0.pod.Pod
	5,  // 1: apocryph.proto.v0.provisionPod.ProvisionPodRequest.payment:type_name -> apocryph.proto.v0.provisionPod.PaymentChannel
	11, // 2: apocryph.proto.v0.provisionPod.UpdatePodRequest.pod:type_name -> apocryph.proto.v0.pod.Pod
	5,  // 3: apocryph.proto.v0.provisionPod.UpdatePodRequest.payment:type_name -> apocryph.proto.v0.provisionPod.PaymentChannel
	10, // 4: apocryph.proto.v0.provisionPod.ProvisionPodResponse.addresses:type_name -> apocryph.proto.v0.provisionPod.ProvisionPodResponse.ExposedHostPort
	9,  // 5: apocryph.proto.v0.provisionPod.PodLogResponse.logEntry:type_name -> apocryph.proto.v0.provisionPod.LogEntry
	0,  // 6: apocryph.proto.v0.provisionPod.ProvisionPodService.ProvisionPod:input_type -> apocryph.proto.v0.provisionPod.ProvisionPodRequest
	3,  // 7: apocryph.proto.v0.provisionPod.ProvisionPodService.UpdatePod:input_type -> apocryph.proto.v0.provisionPod.UpdatePodRequest
	1,  // 8: apocryph.proto.v0.provisionPod.ProvisionPodService.DeletePod:input_type -> apocryph.proto.v0.provisionPod.DeletePodRequest
	7,  // 9: apocryph.proto.v0.provisionPod.ProvisionPodService.GetPodLogs:input_type -> apocryph.proto.v0.provisionPod.PodLogRequest
	6,  // 10: apocryph.proto.v0.provisionPod.ProvisionPodService.ProvisionPod:output_type -> apocryph.proto.v0.provisionPod.ProvisionPodResponse
	6,  // 11: apocryph.proto.v0.provisionPod.ProvisionPodService.UpdatePod:output_type -> apocryph.proto.v0.provisionPod.ProvisionPodResponse
	2,  // 12: apocryph.proto.v0.provisionPod.ProvisionPodService.DeletePod:output_type -> apocryph.proto.v0.provisionPod.DeletePodResponse
	8,  // 13: apocryph.proto.v0.provisionPod.ProvisionPodService.GetPodLogs:output_type -> apocryph.proto.v0.provisionPod.PodLogResponse
	10, // [10:14] is the sub-list for method output_type
	6,  // [6:10] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_provision_pod_proto_init() }
func file_provision_pod_proto_init() {
	if File_provision_pod_proto != nil {
		return
	}
	file_pod_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_provision_pod_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProvisionPodRequest); i {
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
		file_provision_pod_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeletePodRequest); i {
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
		file_provision_pod_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeletePodResponse); i {
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
		file_provision_pod_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdatePodRequest); i {
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
		file_provision_pod_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdatePodResponse); i {
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
		file_provision_pod_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PaymentChannel); i {
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
		file_provision_pod_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProvisionPodResponse); i {
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
		file_provision_pod_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PodLogRequest); i {
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
		file_provision_pod_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PodLogResponse); i {
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
		file_provision_pod_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogEntry); i {
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
		file_provision_pod_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProvisionPodResponse_ExposedHostPort); i {
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
			RawDescriptor: file_provision_pod_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_provision_pod_proto_goTypes,
		DependencyIndexes: file_provision_pod_proto_depIdxs,
		MessageInfos:      file_provision_pod_proto_msgTypes,
	}.Build()
	File_provision_pod_proto = out.File
	file_provision_pod_proto_rawDesc = nil
	file_provision_pod_proto_goTypes = nil
	file_provision_pod_proto_depIdxs = nil
}
