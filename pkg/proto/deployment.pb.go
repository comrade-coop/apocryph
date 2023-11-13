// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.24.4
// source: deployment.proto

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

type Deployment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PodManifestFile string                `protobuf:"bytes,1,opt,name=podManifestFile,proto3" json:"podManifestFile,omitempty"`
	Provider        *ProviderConfig       `protobuf:"bytes,2,opt,name=provider,proto3" json:"provider,omitempty"`
	Payment         *PaymentChannelConfig `protobuf:"bytes,3,opt,name=payment,proto3" json:"payment,omitempty"`
	Images          []*UploadedImage      `protobuf:"bytes,4,rep,name=images,proto3" json:"images,omitempty"`
	Secrets         []*UploadedSecret     `protobuf:"bytes,5,rep,name=secrets,proto3" json:"secrets,omitempty"`
	Deployed        *ProvisionPodResponse `protobuf:"bytes,6,opt,name=deployed,proto3" json:"deployed,omitempty"`
}

func (x *Deployment) Reset() {
	*x = Deployment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_deployment_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Deployment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Deployment) ProtoMessage() {}

func (x *Deployment) ProtoReflect() protoreflect.Message {
	mi := &file_deployment_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Deployment.ProtoReflect.Descriptor instead.
func (*Deployment) Descriptor() ([]byte, []int) {
	return file_deployment_proto_rawDescGZIP(), []int{0}
}

func (x *Deployment) GetPodManifestFile() string {
	if x != nil {
		return x.PodManifestFile
	}
	return ""
}

func (x *Deployment) GetProvider() *ProviderConfig {
	if x != nil {
		return x.Provider
	}
	return nil
}

func (x *Deployment) GetPayment() *PaymentChannelConfig {
	if x != nil {
		return x.Payment
	}
	return nil
}

func (x *Deployment) GetImages() []*UploadedImage {
	if x != nil {
		return x.Images
	}
	return nil
}

func (x *Deployment) GetSecrets() []*UploadedSecret {
	if x != nil {
		return x.Secrets
	}
	return nil
}

func (x *Deployment) GetDeployed() *ProvisionPodResponse {
	if x != nil {
		return x.Deployed
	}
	return nil
}

type ProviderConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EthereumAddress []byte `protobuf:"bytes,1,opt,name=ethereumAddress,proto3" json:"ethereumAddress,omitempty"`
	Libp2PAddress   string `protobuf:"bytes,2,opt,name=libp2pAddress,proto3" json:"libp2pAddress,omitempty"`
}

func (x *ProviderConfig) Reset() {
	*x = ProviderConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_deployment_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProviderConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProviderConfig) ProtoMessage() {}

func (x *ProviderConfig) ProtoReflect() protoreflect.Message {
	mi := &file_deployment_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProviderConfig.ProtoReflect.Descriptor instead.
func (*ProviderConfig) Descriptor() ([]byte, []int) {
	return file_deployment_proto_rawDescGZIP(), []int{1}
}

func (x *ProviderConfig) GetEthereumAddress() []byte {
	if x != nil {
		return x.EthereumAddress
	}
	return nil
}

func (x *ProviderConfig) GetLibp2PAddress() string {
	if x != nil {
		return x.Libp2PAddress
	}
	return ""
}

type PaymentChannelConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChainID                []byte `protobuf:"bytes,1,opt,name=chainID,proto3" json:"chainID,omitempty"`
	PaymentContractAddress []byte `protobuf:"bytes,2,opt,name=paymentContractAddress,proto3" json:"paymentContractAddress,omitempty"`
	PublisherAddress       []byte `protobuf:"bytes,3,opt,name=publisherAddress,proto3" json:"publisherAddress,omitempty"`
	PodID                  []byte `protobuf:"bytes,5,opt,name=podID,proto3" json:"podID,omitempty"`
}

func (x *PaymentChannelConfig) Reset() {
	*x = PaymentChannelConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_deployment_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PaymentChannelConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PaymentChannelConfig) ProtoMessage() {}

func (x *PaymentChannelConfig) ProtoReflect() protoreflect.Message {
	mi := &file_deployment_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PaymentChannelConfig.ProtoReflect.Descriptor instead.
func (*PaymentChannelConfig) Descriptor() ([]byte, []int) {
	return file_deployment_proto_rawDescGZIP(), []int{2}
}

func (x *PaymentChannelConfig) GetChainID() []byte {
	if x != nil {
		return x.ChainID
	}
	return nil
}

func (x *PaymentChannelConfig) GetPaymentContractAddress() []byte {
	if x != nil {
		return x.PaymentContractAddress
	}
	return nil
}

func (x *PaymentChannelConfig) GetPublisherAddress() []byte {
	if x != nil {
		return x.PublisherAddress
	}
	return nil
}

func (x *PaymentChannelConfig) GetPodID() []byte {
	if x != nil {
		return x.PodID
	}
	return nil
}

type UploadedImage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SourceUrl string `protobuf:"bytes,1,opt,name=sourceUrl,proto3" json:"sourceUrl,omitempty"`
	Digest    string `protobuf:"bytes,2,opt,name=digest,proto3" json:"digest,omitempty"`
	Cid       []byte `protobuf:"bytes,3,opt,name=cid,proto3" json:"cid,omitempty"`
	Key       *Key   `protobuf:"bytes,4,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *UploadedImage) Reset() {
	*x = UploadedImage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_deployment_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadedImage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadedImage) ProtoMessage() {}

func (x *UploadedImage) ProtoReflect() protoreflect.Message {
	mi := &file_deployment_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadedImage.ProtoReflect.Descriptor instead.
func (*UploadedImage) Descriptor() ([]byte, []int) {
	return file_deployment_proto_rawDescGZIP(), []int{3}
}

func (x *UploadedImage) GetSourceUrl() string {
	if x != nil {
		return x.SourceUrl
	}
	return ""
}

func (x *UploadedImage) GetDigest() string {
	if x != nil {
		return x.Digest
	}
	return ""
}

func (x *UploadedImage) GetCid() []byte {
	if x != nil {
		return x.Cid
	}
	return nil
}

func (x *UploadedImage) GetKey() *Key {
	if x != nil {
		return x.Key
	}
	return nil
}

type UploadedSecret struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VolumeName string `protobuf:"bytes,1,opt,name=volumeName,proto3" json:"volumeName,omitempty"`
	Sha256Sum  []byte `protobuf:"bytes,2,opt,name=sha256sum,proto3" json:"sha256sum,omitempty"`
	Cid        []byte `protobuf:"bytes,3,opt,name=cid,proto3" json:"cid,omitempty"`
	Key        *Key   `protobuf:"bytes,4,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *UploadedSecret) Reset() {
	*x = UploadedSecret{}
	if protoimpl.UnsafeEnabled {
		mi := &file_deployment_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadedSecret) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadedSecret) ProtoMessage() {}

func (x *UploadedSecret) ProtoReflect() protoreflect.Message {
	mi := &file_deployment_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadedSecret.ProtoReflect.Descriptor instead.
func (*UploadedSecret) Descriptor() ([]byte, []int) {
	return file_deployment_proto_rawDescGZIP(), []int{4}
}

func (x *UploadedSecret) GetVolumeName() string {
	if x != nil {
		return x.VolumeName
	}
	return ""
}

func (x *UploadedSecret) GetSha256Sum() []byte {
	if x != nil {
		return x.Sha256Sum
	}
	return nil
}

func (x *UploadedSecret) GetCid() []byte {
	if x != nil {
		return x.Cid
	}
	return nil
}

func (x *UploadedSecret) GetKey() *Key {
	if x != nil {
		return x.Key
	}
	return nil
}

var File_deployment_proto protoreflect.FileDescriptor

var file_deployment_proto_rawDesc = []byte{
	0x0a, 0x10, 0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x1c, 0x61, 0x70, 0x6f, 0x63, 0x72, 0x79, 0x70, 0x68, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x76, 0x30, 0x2e, 0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x1a, 0x13, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x2d, 0x70, 0x6f, 0x64, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x09, 0x70, 0x6f, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xad, 0x03, 0x0a, 0x0a, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x12,
	0x28, 0x0a, 0x0f, 0x70, 0x6f, 0x64, 0x4d, 0x61, 0x6e, 0x69, 0x66, 0x65, 0x73, 0x74, 0x46, 0x69,
	0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x70, 0x6f, 0x64, 0x4d, 0x61, 0x6e,
	0x69, 0x66, 0x65, 0x73, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x48, 0x0a, 0x08, 0x70, 0x72, 0x6f,
	0x76, 0x69, 0x64, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x61, 0x70,
	0x6f, 0x63, 0x72, 0x79, 0x70, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x30, 0x2e,
	0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x50, 0x72, 0x6f, 0x76, 0x69,
	0x64, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69,
	0x64, 0x65, 0x72, 0x12, 0x4c, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x32, 0x2e, 0x61, 0x70, 0x6f, 0x63, 0x72, 0x79, 0x70, 0x68, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x30, 0x2e, 0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d,
	0x65, 0x6e, 0x74, 0x2e, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x6e,
	0x65, 0x6c, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e,
	0x74, 0x12, 0x43, 0x0a, 0x06, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x2b, 0x2e, 0x61, 0x70, 0x6f, 0x63, 0x72, 0x79, 0x70, 0x68, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x76, 0x30, 0x2e, 0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x64, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x06,
	0x69, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x12, 0x46, 0x0a, 0x07, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74,
	0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x61, 0x70, 0x6f, 0x63, 0x72, 0x79,
	0x70, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x30, 0x2e, 0x64, 0x65, 0x70, 0x6c,
	0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x64, 0x53,
	0x65, 0x63, 0x72, 0x65, 0x74, 0x52, 0x07, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x73, 0x12, 0x50,
	0x0a, 0x08, 0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x34, 0x2e, 0x61, 0x70, 0x6f, 0x63, 0x72, 0x79, 0x70, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x76, 0x30, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x50, 0x6f,
	0x64, 0x2e, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x73, 0x69, 0x6f, 0x6e, 0x50, 0x6f, 0x64, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x08, 0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x65, 0x64,
	0x22, 0x60, 0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x12, 0x28, 0x0a, 0x0f, 0x65, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0f, 0x65, 0x74, 0x68,
	0x65, 0x72, 0x65, 0x75, 0x6d, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x24, 0x0a, 0x0d,
	0x6c, 0x69, 0x62, 0x70, 0x32, 0x70, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0d, 0x6c, 0x69, 0x62, 0x70, 0x32, 0x70, 0x41, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x22, 0xaa, 0x01, 0x0a, 0x14, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x68,
	0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x18, 0x0a, 0x07, 0x63,
	0x68, 0x61, 0x69, 0x6e, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x63, 0x68,
	0x61, 0x69, 0x6e, 0x49, 0x44, 0x12, 0x36, 0x0a, 0x16, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x16, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x6f,
	0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x2a, 0x0a,
	0x10, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x10, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68,
	0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x6f, 0x64,
	0x49, 0x44, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x70, 0x6f, 0x64, 0x49, 0x44, 0x22,
	0x85, 0x01, 0x0a, 0x0d, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x64, 0x49, 0x6d, 0x61, 0x67,
	0x65, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x55, 0x72, 0x6c, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x55, 0x72, 0x6c, 0x12,
	0x16, 0x0a, 0x06, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x63, 0x69, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x63, 0x69, 0x64, 0x12, 0x2c, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x61, 0x70, 0x6f, 0x63, 0x72, 0x79, 0x70,
	0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x30, 0x2e, 0x70, 0x6f, 0x64, 0x2e, 0x4b,
	0x65, 0x79, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x22, 0x8e, 0x01, 0x0a, 0x0e, 0x55, 0x70, 0x6c, 0x6f,
	0x61, 0x64, 0x65, 0x64, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x76, 0x6f,
	0x6c, 0x75, 0x6d, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x68,
	0x61, 0x32, 0x35, 0x36, 0x73, 0x75, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73,
	0x68, 0x61, 0x32, 0x35, 0x36, 0x73, 0x75, 0x6d, 0x12, 0x10, 0x0a, 0x03, 0x63, 0x69, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x63, 0x69, 0x64, 0x12, 0x2c, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x61, 0x70, 0x6f, 0x63, 0x72, 0x79,
	0x70, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x30, 0x2e, 0x70, 0x6f, 0x64, 0x2e,
	0x4b, 0x65, 0x79, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x42, 0x30, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x6d, 0x72, 0x61, 0x64, 0x65, 0x2d, 0x63,
	0x6f, 0x6f, 0x70, 0x2f, 0x74, 0x72, 0x75, 0x73, 0x74, 0x65, 0x64, 0x2d, 0x70, 0x6f, 0x64, 0x73,
	0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_deployment_proto_rawDescOnce sync.Once
	file_deployment_proto_rawDescData = file_deployment_proto_rawDesc
)

func file_deployment_proto_rawDescGZIP() []byte {
	file_deployment_proto_rawDescOnce.Do(func() {
		file_deployment_proto_rawDescData = protoimpl.X.CompressGZIP(file_deployment_proto_rawDescData)
	})
	return file_deployment_proto_rawDescData
}

var file_deployment_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_deployment_proto_goTypes = []interface{}{
	(*Deployment)(nil),           // 0: apocryph.proto.v0.deployment.Deployment
	(*ProviderConfig)(nil),       // 1: apocryph.proto.v0.deployment.ProviderConfig
	(*PaymentChannelConfig)(nil), // 2: apocryph.proto.v0.deployment.PaymentChannelConfig
	(*UploadedImage)(nil),        // 3: apocryph.proto.v0.deployment.UploadedImage
	(*UploadedSecret)(nil),       // 4: apocryph.proto.v0.deployment.UploadedSecret
	(*ProvisionPodResponse)(nil), // 5: apocryph.proto.v0.provisionPod.ProvisionPodResponse
	(*Key)(nil),                  // 6: apocryph.proto.v0.pod.Key
}
var file_deployment_proto_depIdxs = []int32{
	1, // 0: apocryph.proto.v0.deployment.Deployment.provider:type_name -> apocryph.proto.v0.deployment.ProviderConfig
	2, // 1: apocryph.proto.v0.deployment.Deployment.payment:type_name -> apocryph.proto.v0.deployment.PaymentChannelConfig
	3, // 2: apocryph.proto.v0.deployment.Deployment.images:type_name -> apocryph.proto.v0.deployment.UploadedImage
	4, // 3: apocryph.proto.v0.deployment.Deployment.secrets:type_name -> apocryph.proto.v0.deployment.UploadedSecret
	5, // 4: apocryph.proto.v0.deployment.Deployment.deployed:type_name -> apocryph.proto.v0.provisionPod.ProvisionPodResponse
	6, // 5: apocryph.proto.v0.deployment.UploadedImage.key:type_name -> apocryph.proto.v0.pod.Key
	6, // 6: apocryph.proto.v0.deployment.UploadedSecret.key:type_name -> apocryph.proto.v0.pod.Key
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_deployment_proto_init() }
func file_deployment_proto_init() {
	if File_deployment_proto != nil {
		return
	}
	file_provision_pod_proto_init()
	file_pod_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_deployment_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Deployment); i {
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
		file_deployment_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProviderConfig); i {
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
		file_deployment_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PaymentChannelConfig); i {
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
		file_deployment_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadedImage); i {
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
		file_deployment_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadedSecret); i {
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
			RawDescriptor: file_deployment_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_deployment_proto_goTypes,
		DependencyIndexes: file_deployment_proto_depIdxs,
		MessageInfos:      file_deployment_proto_msgTypes,
	}.Build()
	File_deployment_proto = out.File
	file_deployment_proto_rawDesc = nil
	file_deployment_proto_goTypes = nil
	file_deployment_proto_depIdxs = nil
}