// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: ignite/services/plugin/grpc/v1/service.proto

package v1

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ManifestRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ManifestRequest) Reset() {
	*x = ManifestRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ManifestRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ManifestRequest) ProtoMessage() {}

func (x *ManifestRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ManifestRequest.ProtoReflect.Descriptor instead.
func (*ManifestRequest) Descriptor() ([]byte, []int) {
	return file_ignite_services_plugin_grpc_v1_service_proto_rawDescGZIP(), []int{0}
}

type ManifestResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Manifest *Manifest `protobuf:"bytes,1,opt,name=manifest,proto3" json:"manifest,omitempty"`
}

func (x *ManifestResponse) Reset() {
	*x = ManifestResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ManifestResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ManifestResponse) ProtoMessage() {}

func (x *ManifestResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ManifestResponse.ProtoReflect.Descriptor instead.
func (*ManifestResponse) Descriptor() ([]byte, []int) {
	return file_ignite_services_plugin_grpc_v1_service_proto_rawDescGZIP(), []int{1}
}

func (x *ManifestResponse) GetManifest() *Manifest {
	if x != nil {
		return x.Manifest
	}
	return nil
}

type ExecuteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd       *ExecutedCommand `protobuf:"bytes,1,opt,name=cmd,proto3" json:"cmd,omitempty"`
	ClientApi uint32           `protobuf:"varint,2,opt,name=client_api,json=clientApi,proto3" json:"client_api,omitempty"`
}

func (x *ExecuteRequest) Reset() {
	*x = ExecuteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExecuteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecuteRequest) ProtoMessage() {}

func (x *ExecuteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecuteRequest.ProtoReflect.Descriptor instead.
func (*ExecuteRequest) Descriptor() ([]byte, []int) {
	return file_ignite_services_plugin_grpc_v1_service_proto_rawDescGZIP(), []int{2}
}

func (x *ExecuteRequest) GetCmd() *ExecutedCommand {
	if x != nil {
		return x.Cmd
	}
	return nil
}

func (x *ExecuteRequest) GetClientApi() uint32 {
	if x != nil {
		return x.ClientApi
	}
	return 0
}

type ExecuteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ExecuteResponse) Reset() {
	*x = ExecuteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExecuteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecuteResponse) ProtoMessage() {}

func (x *ExecuteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecuteResponse.ProtoReflect.Descriptor instead.
func (*ExecuteResponse) Descriptor() ([]byte, []int) {
	return file_ignite_services_plugin_grpc_v1_service_proto_rawDescGZIP(), []int{3}
}

type ExecuteHookPreRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hook      *ExecutedHook `protobuf:"bytes,1,opt,name=hook,proto3" json:"hook,omitempty"`
	ClientApi uint32        `protobuf:"varint,2,opt,name=client_api,json=clientApi,proto3" json:"client_api,omitempty"`
}

func (x *ExecuteHookPreRequest) Reset() {
	*x = ExecuteHookPreRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExecuteHookPreRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecuteHookPreRequest) ProtoMessage() {}

func (x *ExecuteHookPreRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecuteHookPreRequest.ProtoReflect.Descriptor instead.
func (*ExecuteHookPreRequest) Descriptor() ([]byte, []int) {
	return file_ignite_services_plugin_grpc_v1_service_proto_rawDescGZIP(), []int{4}
}

func (x *ExecuteHookPreRequest) GetHook() *ExecutedHook {
	if x != nil {
		return x.Hook
	}
	return nil
}

func (x *ExecuteHookPreRequest) GetClientApi() uint32 {
	if x != nil {
		return x.ClientApi
	}
	return 0
}

type ExecuteHookPreResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ExecuteHookPreResponse) Reset() {
	*x = ExecuteHookPreResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExecuteHookPreResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecuteHookPreResponse) ProtoMessage() {}

func (x *ExecuteHookPreResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecuteHookPreResponse.ProtoReflect.Descriptor instead.
func (*ExecuteHookPreResponse) Descriptor() ([]byte, []int) {
	return file_ignite_services_plugin_grpc_v1_service_proto_rawDescGZIP(), []int{5}
}

type ExecuteHookPostRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hook      *ExecutedHook `protobuf:"bytes,1,opt,name=hook,proto3" json:"hook,omitempty"`
	ClientApi uint32        `protobuf:"varint,2,opt,name=client_api,json=clientApi,proto3" json:"client_api,omitempty"`
}

func (x *ExecuteHookPostRequest) Reset() {
	*x = ExecuteHookPostRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExecuteHookPostRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecuteHookPostRequest) ProtoMessage() {}

func (x *ExecuteHookPostRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecuteHookPostRequest.ProtoReflect.Descriptor instead.
func (*ExecuteHookPostRequest) Descriptor() ([]byte, []int) {
	return file_ignite_services_plugin_grpc_v1_service_proto_rawDescGZIP(), []int{6}
}

func (x *ExecuteHookPostRequest) GetHook() *ExecutedHook {
	if x != nil {
		return x.Hook
	}
	return nil
}

func (x *ExecuteHookPostRequest) GetClientApi() uint32 {
	if x != nil {
		return x.ClientApi
	}
	return 0
}

type ExecuteHookPostResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ExecuteHookPostResponse) Reset() {
	*x = ExecuteHookPostResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExecuteHookPostResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecuteHookPostResponse) ProtoMessage() {}

func (x *ExecuteHookPostResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecuteHookPostResponse.ProtoReflect.Descriptor instead.
func (*ExecuteHookPostResponse) Descriptor() ([]byte, []int) {
	return file_ignite_services_plugin_grpc_v1_service_proto_rawDescGZIP(), []int{7}
}

type ExecuteHookCleanUpRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hook      *ExecutedHook `protobuf:"bytes,1,opt,name=hook,proto3" json:"hook,omitempty"`
	ClientApi uint32        `protobuf:"varint,2,opt,name=client_api,json=clientApi,proto3" json:"client_api,omitempty"`
}

func (x *ExecuteHookCleanUpRequest) Reset() {
	*x = ExecuteHookCleanUpRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExecuteHookCleanUpRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecuteHookCleanUpRequest) ProtoMessage() {}

func (x *ExecuteHookCleanUpRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecuteHookCleanUpRequest.ProtoReflect.Descriptor instead.
func (*ExecuteHookCleanUpRequest) Descriptor() ([]byte, []int) {
	return file_ignite_services_plugin_grpc_v1_service_proto_rawDescGZIP(), []int{8}
}

func (x *ExecuteHookCleanUpRequest) GetHook() *ExecutedHook {
	if x != nil {
		return x.Hook
	}
	return nil
}

func (x *ExecuteHookCleanUpRequest) GetClientApi() uint32 {
	if x != nil {
		return x.ClientApi
	}
	return 0
}

type ExecuteHookCleanUpResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ExecuteHookCleanUpResponse) Reset() {
	*x = ExecuteHookCleanUpResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExecuteHookCleanUpResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecuteHookCleanUpResponse) ProtoMessage() {}

func (x *ExecuteHookCleanUpResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecuteHookCleanUpResponse.ProtoReflect.Descriptor instead.
func (*ExecuteHookCleanUpResponse) Descriptor() ([]byte, []int) {
	return file_ignite_services_plugin_grpc_v1_service_proto_rawDescGZIP(), []int{9}
}

type GetChainInfoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetChainInfoRequest) Reset() {
	*x = GetChainInfoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetChainInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetChainInfoRequest) ProtoMessage() {}

func (x *GetChainInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetChainInfoRequest.ProtoReflect.Descriptor instead.
func (*GetChainInfoRequest) Descriptor() ([]byte, []int) {
	return file_ignite_services_plugin_grpc_v1_service_proto_rawDescGZIP(), []int{10}
}

type GetChainInfoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChainInfo *ChainInfo `protobuf:"bytes,1,opt,name=chain_info,json=chainInfo,proto3" json:"chain_info,omitempty"`
}

func (x *GetChainInfoResponse) Reset() {
	*x = GetChainInfoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetChainInfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetChainInfoResponse) ProtoMessage() {}

func (x *GetChainInfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetChainInfoResponse.ProtoReflect.Descriptor instead.
func (*GetChainInfoResponse) Descriptor() ([]byte, []int) {
	return file_ignite_services_plugin_grpc_v1_service_proto_rawDescGZIP(), []int{11}
}

func (x *GetChainInfoResponse) GetChainInfo() *ChainInfo {
	if x != nil {
		return x.ChainInfo
	}
	return nil
}

var File_ignite_services_plugin_grpc_v1_service_proto protoreflect.FileDescriptor

var file_ignite_services_plugin_grpc_v1_service_proto_rawDesc = []byte{
	0x0a, 0x2c, 0x69, 0x67, 0x6e, 0x69, 0x74, 0x65, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x73, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31,
	0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1e,
	0x69, 0x67, 0x6e, 0x69, 0x74, 0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e,
	0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x1a, 0x2f,
	0x69, 0x67, 0x6e, 0x69, 0x74, 0x65, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2f,
	0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x63,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x2e, 0x69, 0x67, 0x6e, 0x69, 0x74, 0x65, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73,
	0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x76, 0x31, 0x2f,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x11, 0x0a, 0x0f, 0x4d, 0x61, 0x6e, 0x69, 0x66, 0x65, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x22, 0x58, 0x0a, 0x10, 0x4d, 0x61, 0x6e, 0x69, 0x66, 0x65, 0x73, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x44, 0x0a, 0x08, 0x6d, 0x61, 0x6e, 0x69, 0x66, 0x65,
	0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x69, 0x67, 0x6e, 0x69, 0x74,
	0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x61, 0x6e, 0x69, 0x66, 0x65,
	0x73, 0x74, 0x52, 0x08, 0x6d, 0x61, 0x6e, 0x69, 0x66, 0x65, 0x73, 0x74, 0x22, 0x72, 0x0a, 0x0e,
	0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x41,
	0x0a, 0x03, 0x63, 0x6d, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2f, 0x2e, 0x69, 0x67,
	0x6e, 0x69, 0x74, 0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x70, 0x6c,
	0x75, 0x67, 0x69, 0x6e, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x78, 0x65,
	0x63, 0x75, 0x74, 0x65, 0x64, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x03, 0x63, 0x6d,
	0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x61, 0x70, 0x69, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x41, 0x70, 0x69,
	0x22, 0x11, 0x0a, 0x0f, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x78, 0x0a, 0x15, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x48, 0x6f,
	0x6f, 0x6b, 0x50, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x40, 0x0a, 0x04,
	0x68, 0x6f, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x69, 0x67, 0x6e,
	0x69, 0x74, 0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x70, 0x6c, 0x75,
	0x67, 0x69, 0x6e, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x78, 0x65, 0x63,
	0x75, 0x74, 0x65, 0x64, 0x48, 0x6f, 0x6f, 0x6b, 0x52, 0x04, 0x68, 0x6f, 0x6f, 0x6b, 0x12, 0x1d,
	0x0a, 0x0a, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x61, 0x70, 0x69, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x41, 0x70, 0x69, 0x22, 0x18, 0x0a,
	0x16, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x48, 0x6f, 0x6f, 0x6b, 0x50, 0x72, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x79, 0x0a, 0x16, 0x45, 0x78, 0x65, 0x63, 0x75,
	0x74, 0x65, 0x48, 0x6f, 0x6f, 0x6b, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x40, 0x0a, 0x04, 0x68, 0x6f, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x2c, 0x2e, 0x69, 0x67, 0x6e, 0x69, 0x74, 0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x73, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31,
	0x2e, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x64, 0x48, 0x6f, 0x6f, 0x6b, 0x52, 0x04, 0x68,
	0x6f, 0x6f, 0x6b, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x61, 0x70,
	0x69, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x41,
	0x70, 0x69, 0x22, 0x19, 0x0a, 0x17, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x48, 0x6f, 0x6f,
	0x6b, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x7c, 0x0a,
	0x19, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x48, 0x6f, 0x6f, 0x6b, 0x43, 0x6c, 0x65, 0x61,
	0x6e, 0x55, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x40, 0x0a, 0x04, 0x68, 0x6f,
	0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x69, 0x67, 0x6e, 0x69, 0x74,
	0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74,
	0x65, 0x64, 0x48, 0x6f, 0x6f, 0x6b, 0x52, 0x04, 0x68, 0x6f, 0x6f, 0x6b, 0x12, 0x1d, 0x0a, 0x0a,
	0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x61, 0x70, 0x69, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x41, 0x70, 0x69, 0x22, 0x1c, 0x0a, 0x1a, 0x45,
	0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x48, 0x6f, 0x6f, 0x6b, 0x43, 0x6c, 0x65, 0x61, 0x6e, 0x55,
	0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x15, 0x0a, 0x13, 0x47, 0x65, 0x74,
	0x43, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x22, 0x60, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x48, 0x0a, 0x0a, 0x63, 0x68, 0x61, 0x69,
	0x6e, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x69,
	0x67, 0x6e, 0x69, 0x74, 0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x70,
	0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x68,
	0x61, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x09, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x6e,
	0x66, 0x6f, 0x32, 0x81, 0x05, 0x0a, 0x10, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x6d, 0x0a, 0x08, 0x4d, 0x61, 0x6e, 0x69, 0x66,
	0x65, 0x73, 0x74, 0x12, 0x2f, 0x2e, 0x69, 0x67, 0x6e, 0x69, 0x74, 0x65, 0x2e, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x67, 0x72, 0x70,
	0x63, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x61, 0x6e, 0x69, 0x66, 0x65, 0x73, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x30, 0x2e, 0x69, 0x67, 0x6e, 0x69, 0x74, 0x65, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x61, 0x6e, 0x69, 0x66, 0x65, 0x73, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x6a, 0x0a, 0x07, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74,
	0x65, 0x12, 0x2e, 0x2e, 0x69, 0x67, 0x6e, 0x69, 0x74, 0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x73, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e,
	0x76, 0x31, 0x2e, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x2f, 0x2e, 0x69, 0x67, 0x6e, 0x69, 0x74, 0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x73, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e,
	0x76, 0x31, 0x2e, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x7f, 0x0a, 0x0e, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x48, 0x6f, 0x6f,
	0x6b, 0x50, 0x72, 0x65, 0x12, 0x35, 0x2e, 0x69, 0x67, 0x6e, 0x69, 0x74, 0x65, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x48, 0x6f, 0x6f,
	0x6b, 0x50, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x36, 0x2e, 0x69, 0x67,
	0x6e, 0x69, 0x74, 0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x70, 0x6c,
	0x75, 0x67, 0x69, 0x6e, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x78, 0x65,
	0x63, 0x75, 0x74, 0x65, 0x48, 0x6f, 0x6f, 0x6b, 0x50, 0x72, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x82, 0x01, 0x0a, 0x0f, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x48,
	0x6f, 0x6f, 0x6b, 0x50, 0x6f, 0x73, 0x74, 0x12, 0x36, 0x2e, 0x69, 0x67, 0x6e, 0x69, 0x74, 0x65,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e,
	0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65,
	0x48, 0x6f, 0x6f, 0x6b, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x37, 0x2e, 0x69, 0x67, 0x6e, 0x69, 0x74, 0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x73, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31,
	0x2e, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x48, 0x6f, 0x6f, 0x6b, 0x50, 0x6f, 0x73, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x8b, 0x01, 0x0a, 0x12, 0x45, 0x78, 0x65,
	0x63, 0x75, 0x74, 0x65, 0x48, 0x6f, 0x6f, 0x6b, 0x43, 0x6c, 0x65, 0x61, 0x6e, 0x55, 0x70, 0x12,
	0x39, 0x2e, 0x69, 0x67, 0x6e, 0x69, 0x74, 0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x73, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31,
	0x2e, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x48, 0x6f, 0x6f, 0x6b, 0x43, 0x6c, 0x65, 0x61,
	0x6e, 0x55, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x3a, 0x2e, 0x69, 0x67, 0x6e,
	0x69, 0x74, 0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x70, 0x6c, 0x75,
	0x67, 0x69, 0x6e, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x78, 0x65, 0x63,
	0x75, 0x74, 0x65, 0x48, 0x6f, 0x6f, 0x6b, 0x43, 0x6c, 0x65, 0x61, 0x6e, 0x55, 0x70, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x8d, 0x01, 0x0a, 0x10, 0x43, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x41, 0x50, 0x49, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x79, 0x0a, 0x0c, 0x47,
	0x65, 0x74, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x33, 0x2e, 0x69, 0x67,
	0x6e, 0x69, 0x74, 0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x70, 0x6c,
	0x75, 0x67, 0x69, 0x6e, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74,
	0x43, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x34, 0x2e, 0x69, 0x67, 0x6e, 0x69, 0x74, 0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x2e, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x76,
	0x31, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x3a, 0x5a, 0x38, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x67, 0x6e, 0x69, 0x74, 0x65, 0x2f, 0x63, 0x6c, 0x69, 0x2f,
	0x76, 0x32, 0x39, 0x2f, 0x69, 0x67, 0x6e, 0x69, 0x74, 0x65, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x73, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f,
	0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ignite_services_plugin_grpc_v1_service_proto_rawDescOnce sync.Once
	file_ignite_services_plugin_grpc_v1_service_proto_rawDescData = file_ignite_services_plugin_grpc_v1_service_proto_rawDesc
)

func file_ignite_services_plugin_grpc_v1_service_proto_rawDescGZIP() []byte {
	file_ignite_services_plugin_grpc_v1_service_proto_rawDescOnce.Do(func() {
		file_ignite_services_plugin_grpc_v1_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_ignite_services_plugin_grpc_v1_service_proto_rawDescData)
	})
	return file_ignite_services_plugin_grpc_v1_service_proto_rawDescData
}

var file_ignite_services_plugin_grpc_v1_service_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_ignite_services_plugin_grpc_v1_service_proto_goTypes = []any{
	(*ManifestRequest)(nil),            // 0: ignite.services.plugin.grpc.v1.ManifestRequest
	(*ManifestResponse)(nil),           // 1: ignite.services.plugin.grpc.v1.ManifestResponse
	(*ExecuteRequest)(nil),             // 2: ignite.services.plugin.grpc.v1.ExecuteRequest
	(*ExecuteResponse)(nil),            // 3: ignite.services.plugin.grpc.v1.ExecuteResponse
	(*ExecuteHookPreRequest)(nil),      // 4: ignite.services.plugin.grpc.v1.ExecuteHookPreRequest
	(*ExecuteHookPreResponse)(nil),     // 5: ignite.services.plugin.grpc.v1.ExecuteHookPreResponse
	(*ExecuteHookPostRequest)(nil),     // 6: ignite.services.plugin.grpc.v1.ExecuteHookPostRequest
	(*ExecuteHookPostResponse)(nil),    // 7: ignite.services.plugin.grpc.v1.ExecuteHookPostResponse
	(*ExecuteHookCleanUpRequest)(nil),  // 8: ignite.services.plugin.grpc.v1.ExecuteHookCleanUpRequest
	(*ExecuteHookCleanUpResponse)(nil), // 9: ignite.services.plugin.grpc.v1.ExecuteHookCleanUpResponse
	(*GetChainInfoRequest)(nil),        // 10: ignite.services.plugin.grpc.v1.GetChainInfoRequest
	(*GetChainInfoResponse)(nil),       // 11: ignite.services.plugin.grpc.v1.GetChainInfoResponse
	(*Manifest)(nil),                   // 12: ignite.services.plugin.grpc.v1.Manifest
	(*ExecutedCommand)(nil),            // 13: ignite.services.plugin.grpc.v1.ExecutedCommand
	(*ExecutedHook)(nil),               // 14: ignite.services.plugin.grpc.v1.ExecutedHook
	(*ChainInfo)(nil),                  // 15: ignite.services.plugin.grpc.v1.ChainInfo
}
var file_ignite_services_plugin_grpc_v1_service_proto_depIdxs = []int32{
	12, // 0: ignite.services.plugin.grpc.v1.ManifestResponse.manifest:type_name -> ignite.services.plugin.grpc.v1.Manifest
	13, // 1: ignite.services.plugin.grpc.v1.ExecuteRequest.cmd:type_name -> ignite.services.plugin.grpc.v1.ExecutedCommand
	14, // 2: ignite.services.plugin.grpc.v1.ExecuteHookPreRequest.hook:type_name -> ignite.services.plugin.grpc.v1.ExecutedHook
	14, // 3: ignite.services.plugin.grpc.v1.ExecuteHookPostRequest.hook:type_name -> ignite.services.plugin.grpc.v1.ExecutedHook
	14, // 4: ignite.services.plugin.grpc.v1.ExecuteHookCleanUpRequest.hook:type_name -> ignite.services.plugin.grpc.v1.ExecutedHook
	15, // 5: ignite.services.plugin.grpc.v1.GetChainInfoResponse.chain_info:type_name -> ignite.services.plugin.grpc.v1.ChainInfo
	0,  // 6: ignite.services.plugin.grpc.v1.InterfaceService.Manifest:input_type -> ignite.services.plugin.grpc.v1.ManifestRequest
	2,  // 7: ignite.services.plugin.grpc.v1.InterfaceService.Execute:input_type -> ignite.services.plugin.grpc.v1.ExecuteRequest
	4,  // 8: ignite.services.plugin.grpc.v1.InterfaceService.ExecuteHookPre:input_type -> ignite.services.plugin.grpc.v1.ExecuteHookPreRequest
	6,  // 9: ignite.services.plugin.grpc.v1.InterfaceService.ExecuteHookPost:input_type -> ignite.services.plugin.grpc.v1.ExecuteHookPostRequest
	8,  // 10: ignite.services.plugin.grpc.v1.InterfaceService.ExecuteHookCleanUp:input_type -> ignite.services.plugin.grpc.v1.ExecuteHookCleanUpRequest
	10, // 11: ignite.services.plugin.grpc.v1.ClientAPIService.GetChainInfo:input_type -> ignite.services.plugin.grpc.v1.GetChainInfoRequest
	1,  // 12: ignite.services.plugin.grpc.v1.InterfaceService.Manifest:output_type -> ignite.services.plugin.grpc.v1.ManifestResponse
	3,  // 13: ignite.services.plugin.grpc.v1.InterfaceService.Execute:output_type -> ignite.services.plugin.grpc.v1.ExecuteResponse
	5,  // 14: ignite.services.plugin.grpc.v1.InterfaceService.ExecuteHookPre:output_type -> ignite.services.plugin.grpc.v1.ExecuteHookPreResponse
	7,  // 15: ignite.services.plugin.grpc.v1.InterfaceService.ExecuteHookPost:output_type -> ignite.services.plugin.grpc.v1.ExecuteHookPostResponse
	9,  // 16: ignite.services.plugin.grpc.v1.InterfaceService.ExecuteHookCleanUp:output_type -> ignite.services.plugin.grpc.v1.ExecuteHookCleanUpResponse
	11, // 17: ignite.services.plugin.grpc.v1.ClientAPIService.GetChainInfo:output_type -> ignite.services.plugin.grpc.v1.GetChainInfoResponse
	12, // [12:18] is the sub-list for method output_type
	6,  // [6:12] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_ignite_services_plugin_grpc_v1_service_proto_init() }
func file_ignite_services_plugin_grpc_v1_service_proto_init() {
	if File_ignite_services_plugin_grpc_v1_service_proto != nil {
		return
	}
	file_ignite_services_plugin_grpc_v1_client_api_proto_init()
	file_ignite_services_plugin_grpc_v1_interface_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*ManifestRequest); i {
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
		file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*ManifestResponse); i {
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
		file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*ExecuteRequest); i {
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
		file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*ExecuteResponse); i {
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
		file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*ExecuteHookPreRequest); i {
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
		file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*ExecuteHookPreResponse); i {
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
		file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*ExecuteHookPostRequest); i {
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
		file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*ExecuteHookPostResponse); i {
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
		file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*ExecuteHookCleanUpRequest); i {
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
		file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[9].Exporter = func(v any, i int) any {
			switch v := v.(*ExecuteHookCleanUpResponse); i {
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
		file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[10].Exporter = func(v any, i int) any {
			switch v := v.(*GetChainInfoRequest); i {
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
		file_ignite_services_plugin_grpc_v1_service_proto_msgTypes[11].Exporter = func(v any, i int) any {
			switch v := v.(*GetChainInfoResponse); i {
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
			RawDescriptor: file_ignite_services_plugin_grpc_v1_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_ignite_services_plugin_grpc_v1_service_proto_goTypes,
		DependencyIndexes: file_ignite_services_plugin_grpc_v1_service_proto_depIdxs,
		MessageInfos:      file_ignite_services_plugin_grpc_v1_service_proto_msgTypes,
	}.Build()
	File_ignite_services_plugin_grpc_v1_service_proto = out.File
	file_ignite_services_plugin_grpc_v1_service_proto_rawDesc = nil
	file_ignite_services_plugin_grpc_v1_service_proto_goTypes = nil
	file_ignite_services_plugin_grpc_v1_service_proto_depIdxs = nil
}
