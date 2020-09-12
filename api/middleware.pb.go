// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.3
// source: api/middleware.proto

package api

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Result int32

const (
	Result_OK  Result = 0
	Result_ERR Result = 1
)

// Enum value maps for Result.
var (
	Result_name = map[int32]string{
		0: "OK",
		1: "ERR",
	}
	Result_value = map[string]int32{
		"OK":  0,
		"ERR": 1,
	}
)

func (x Result) Enum() *Result {
	p := new(Result)
	*p = x
	return p
}

func (x Result) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Result) Descriptor() protoreflect.EnumDescriptor {
	return file_api_middleware_proto_enumTypes[0].Descriptor()
}

func (Result) Type() protoreflect.EnumType {
	return &file_api_middleware_proto_enumTypes[0]
}

func (x Result) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Result.Descriptor instead.
func (Result) EnumDescriptor() ([]byte, []int) {
	return file_api_middleware_proto_rawDescGZIP(), []int{0}
}

type InitChannelResponse_Result int32

const (
	InitChannelResponse_ACK InitChannelResponse_Result = 0
	InitChannelResponse_ERR InitChannelResponse_Result = 1
)

// Enum value maps for InitChannelResponse_Result.
var (
	InitChannelResponse_Result_name = map[int32]string{
		0: "ACK",
		1: "ERR",
	}
	InitChannelResponse_Result_value = map[string]int32{
		"ACK": 0,
		"ERR": 1,
	}
)

func (x InitChannelResponse_Result) Enum() *InitChannelResponse_Result {
	p := new(InitChannelResponse_Result)
	*p = x
	return p
}

func (x InitChannelResponse_Result) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (InitChannelResponse_Result) Descriptor() protoreflect.EnumDescriptor {
	return file_api_middleware_proto_enumTypes[1].Descriptor()
}

func (InitChannelResponse_Result) Type() protoreflect.EnumType {
	return &file_api_middleware_proto_enumTypes[1]
}

func (x InitChannelResponse_Result) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use InitChannelResponse_Result.Descriptor instead.
func (InitChannelResponse_Result) EnumDescriptor() ([]byte, []int) {
	return file_api_middleware_proto_rawDescGZIP(), []int{8, 0}
}

type StartChannelResponse_Result int32

const (
	StartChannelResponse_ACK StartChannelResponse_Result = 0
	StartChannelResponse_ERR StartChannelResponse_Result = 1
)

// Enum value maps for StartChannelResponse_Result.
var (
	StartChannelResponse_Result_name = map[int32]string{
		0: "ACK",
		1: "ERR",
	}
	StartChannelResponse_Result_value = map[string]int32{
		"ACK": 0,
		"ERR": 1,
	}
)

func (x StartChannelResponse_Result) Enum() *StartChannelResponse_Result {
	p := new(StartChannelResponse_Result)
	*p = x
	return p
}

func (x StartChannelResponse_Result) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (StartChannelResponse_Result) Descriptor() protoreflect.EnumDescriptor {
	return file_api_middleware_proto_enumTypes[2].Descriptor()
}

func (StartChannelResponse_Result) Type() protoreflect.EnumType {
	return &file_api_middleware_proto_enumTypes[2]
}

func (x StartChannelResponse_Result) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use StartChannelResponse_Result.Descriptor instead.
func (StartChannelResponse_Result) EnumDescriptor() ([]byte, []int) {
	return file_api_middleware_proto_rawDescGZIP(), []int{10, 0}
}

type AppSendResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result Result `protobuf:"varint,1,opt,name=result,proto3,enum=api.Result" json:"result,omitempty"`
}

func (x *AppSendResponse) Reset() {
	*x = AppSendResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_middleware_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppSendResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppSendResponse) ProtoMessage() {}

func (x *AppSendResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_middleware_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppSendResponse.ProtoReflect.Descriptor instead.
func (*AppSendResponse) Descriptor() ([]byte, []int) {
	return file_api_middleware_proto_rawDescGZIP(), []int{0}
}

func (x *AppSendResponse) GetResult() Result {
	if x != nil {
		return x.Result
	}
	return Result_OK
}

type AppRecvRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChannelId   string `protobuf:"bytes,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	Participant string `protobuf:"bytes,2,opt,name=participant,proto3" json:"participant,omitempty"`
}

func (x *AppRecvRequest) Reset() {
	*x = AppRecvRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_middleware_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppRecvRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppRecvRequest) ProtoMessage() {}

func (x *AppRecvRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_middleware_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppRecvRequest.ProtoReflect.Descriptor instead.
func (*AppRecvRequest) Descriptor() ([]byte, []int) {
	return file_api_middleware_proto_rawDescGZIP(), []int{1}
}

func (x *AppRecvRequest) GetChannelId() string {
	if x != nil {
		return x.ChannelId
	}
	return ""
}

func (x *AppRecvRequest) GetParticipant() string {
	if x != nil {
		return x.Participant
	}
	return ""
}

type RegisterChannelRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RequirementsContract *Contract `protobuf:"bytes,1,opt,name=requirements_contract,json=requirementsContract,proto3" json:"requirements_contract,omitempty"`
}

func (x *RegisterChannelRequest) Reset() {
	*x = RegisterChannelRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_middleware_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterChannelRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterChannelRequest) ProtoMessage() {}

func (x *RegisterChannelRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_middleware_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterChannelRequest.ProtoReflect.Descriptor instead.
func (*RegisterChannelRequest) Descriptor() ([]byte, []int) {
	return file_api_middleware_proto_rawDescGZIP(), []int{2}
}

func (x *RegisterChannelRequest) GetRequirementsContract() *Contract {
	if x != nil {
		return x.RequirementsContract
	}
	return nil
}

type RegisterChannelResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChannelId string `protobuf:"bytes,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
}

func (x *RegisterChannelResponse) Reset() {
	*x = RegisterChannelResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_middleware_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterChannelResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterChannelResponse) ProtoMessage() {}

func (x *RegisterChannelResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_middleware_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterChannelResponse.ProtoReflect.Descriptor instead.
func (*RegisterChannelResponse) Descriptor() ([]byte, []int) {
	return file_api_middleware_proto_rawDescGZIP(), []int{3}
}

func (x *RegisterChannelResponse) GetChannelId() string {
	if x != nil {
		return x.ChannelId
	}
	return ""
}

type RegisterAppRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProviderContract *Contract `protobuf:"bytes,1,opt,name=provider_contract,json=providerContract,proto3" json:"provider_contract,omitempty"`
}

func (x *RegisterAppRequest) Reset() {
	*x = RegisterAppRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_middleware_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterAppRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterAppRequest) ProtoMessage() {}

func (x *RegisterAppRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_middleware_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterAppRequest.ProtoReflect.Descriptor instead.
func (*RegisterAppRequest) Descriptor() ([]byte, []int) {
	return file_api_middleware_proto_rawDescGZIP(), []int{4}
}

func (x *RegisterAppRequest) GetProviderContract() *Contract {
	if x != nil {
		return x.ProviderContract
	}
	return nil
}

// whenever a new channel that involves this app is started, the middleware needs to notify the local app
type RegisterAppResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to AckOrNew:
	//	*RegisterAppResponse_AppId
	//	*RegisterAppResponse_Notification
	AckOrNew isRegisterAppResponse_AckOrNew `protobuf_oneof:"ack_or_new"`
}

func (x *RegisterAppResponse) Reset() {
	*x = RegisterAppResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_middleware_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterAppResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterAppResponse) ProtoMessage() {}

func (x *RegisterAppResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_middleware_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterAppResponse.ProtoReflect.Descriptor instead.
func (*RegisterAppResponse) Descriptor() ([]byte, []int) {
	return file_api_middleware_proto_rawDescGZIP(), []int{5}
}

func (m *RegisterAppResponse) GetAckOrNew() isRegisterAppResponse_AckOrNew {
	if m != nil {
		return m.AckOrNew
	}
	return nil
}

func (x *RegisterAppResponse) GetAppId() string {
	if x, ok := x.GetAckOrNew().(*RegisterAppResponse_AppId); ok {
		return x.AppId
	}
	return ""
}

func (x *RegisterAppResponse) GetNotification() *InitChannelNotification {
	if x, ok := x.GetAckOrNew().(*RegisterAppResponse_Notification); ok {
		return x.Notification
	}
	return nil
}

type isRegisterAppResponse_AckOrNew interface {
	isRegisterAppResponse_AckOrNew()
}

type RegisterAppResponse_AppId struct {
	AppId string `protobuf:"bytes,1,opt,name=app_id,json=appId,proto3,oneof"`
}

type RegisterAppResponse_Notification struct {
	Notification *InitChannelNotification `protobuf:"bytes,2,opt,name=notification,proto3,oneof"`
}

func (*RegisterAppResponse_AppId) isRegisterAppResponse_AckOrNew() {}

func (*RegisterAppResponse_Notification) isRegisterAppResponse_AckOrNew() {}

// this is what a registered app receives whenever a new channel is initiated for that app
// the app has to communicate with the middleware using UseChannel with this new channel_id
type InitChannelNotification struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChannelId string `protobuf:"bytes,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
}

func (x *InitChannelNotification) Reset() {
	*x = InitChannelNotification{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_middleware_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InitChannelNotification) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InitChannelNotification) ProtoMessage() {}

func (x *InitChannelNotification) ProtoReflect() protoreflect.Message {
	mi := &file_api_middleware_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InitChannelNotification.ProtoReflect.Descriptor instead.
func (*InitChannelNotification) Descriptor() ([]byte, []int) {
	return file_api_middleware_proto_rawDescGZIP(), []int{6}
}

func (x *InitChannelNotification) GetChannelId() string {
	if x != nil {
		return x.ChannelId
	}
	return ""
}

// This is something that is sent by the Broker to providers to notify that a new channel is starting
type InitChannelRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChannelId    string                        `protobuf:"bytes,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	AppId        string                        `protobuf:"bytes,2,opt,name=app_id,json=appId,proto3" json:"app_id,omitempty"` // which app behind the middleware is being notified
	Participants map[string]*RemoteParticipant `protobuf:"bytes,3,rep,name=participants,proto3" json:"participants,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *InitChannelRequest) Reset() {
	*x = InitChannelRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_middleware_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InitChannelRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InitChannelRequest) ProtoMessage() {}

func (x *InitChannelRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_middleware_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InitChannelRequest.ProtoReflect.Descriptor instead.
func (*InitChannelRequest) Descriptor() ([]byte, []int) {
	return file_api_middleware_proto_rawDescGZIP(), []int{7}
}

func (x *InitChannelRequest) GetChannelId() string {
	if x != nil {
		return x.ChannelId
	}
	return ""
}

func (x *InitChannelRequest) GetAppId() string {
	if x != nil {
		return x.AppId
	}
	return ""
}

func (x *InitChannelRequest) GetParticipants() map[string]*RemoteParticipant {
	if x != nil {
		return x.Participants
	}
	return nil
}

type InitChannelResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result InitChannelResponse_Result `protobuf:"varint,1,opt,name=result,proto3,enum=api.InitChannelResponse_Result" json:"result,omitempty"`
}

func (x *InitChannelResponse) Reset() {
	*x = InitChannelResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_middleware_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InitChannelResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InitChannelResponse) ProtoMessage() {}

func (x *InitChannelResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_middleware_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InitChannelResponse.ProtoReflect.Descriptor instead.
func (*InitChannelResponse) Descriptor() ([]byte, []int) {
	return file_api_middleware_proto_rawDescGZIP(), []int{8}
}

func (x *InitChannelResponse) GetResult() InitChannelResponse_Result {
	if x != nil {
		return x.Result
	}
	return InitChannelResponse_ACK
}

type StartChannelRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChannelId string `protobuf:"bytes,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	AppId     string `protobuf:"bytes,2,opt,name=app_id,json=appId,proto3" json:"app_id,omitempty"`
}

func (x *StartChannelRequest) Reset() {
	*x = StartChannelRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_middleware_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartChannelRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartChannelRequest) ProtoMessage() {}

func (x *StartChannelRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_middleware_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartChannelRequest.ProtoReflect.Descriptor instead.
func (*StartChannelRequest) Descriptor() ([]byte, []int) {
	return file_api_middleware_proto_rawDescGZIP(), []int{9}
}

func (x *StartChannelRequest) GetChannelId() string {
	if x != nil {
		return x.ChannelId
	}
	return ""
}

func (x *StartChannelRequest) GetAppId() string {
	if x != nil {
		return x.AppId
	}
	return ""
}

type StartChannelResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result StartChannelResponse_Result `protobuf:"varint,1,opt,name=result,proto3,enum=api.StartChannelResponse_Result" json:"result,omitempty"`
}

func (x *StartChannelResponse) Reset() {
	*x = StartChannelResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_middleware_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartChannelResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartChannelResponse) ProtoMessage() {}

func (x *StartChannelResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_middleware_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartChannelResponse.ProtoReflect.Descriptor instead.
func (*StartChannelResponse) Descriptor() ([]byte, []int) {
	return file_api_middleware_proto_rawDescGZIP(), []int{10}
}

func (x *StartChannelResponse) GetResult() StartChannelResponse_Result {
	if x != nil {
		return x.Result
	}
	return StartChannelResponse_ACK
}

var File_api_middleware_proto protoreflect.FileDescriptor

var file_api_middleware_proto_rawDesc = []byte{
	0x0a, 0x14, 0x61, 0x70, 0x69, 0x2f, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x61, 0x70, 0x69, 0x1a, 0x15, 0x61, 0x70, 0x69,
	0x2f, 0x61, 0x70, 0x70, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x13, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x10, 0x61, 0x70, 0x69, 0x2f, 0x62, 0x72, 0x6f,
	0x6b, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x36, 0x0a, 0x0f, 0x41, 0x70, 0x70,
	0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x23, 0x0a, 0x06,
	0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0b, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x22, 0x51, 0x0a, 0x0e, 0x41, 0x70, 0x70, 0x52, 0x65, 0x63, 0x76, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c,
	0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x61, 0x72, 0x74, 0x69, 0x63, 0x69, 0x70, 0x61, 0x6e,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x61, 0x72, 0x74, 0x69, 0x63, 0x69,
	0x70, 0x61, 0x6e, 0x74, 0x22, 0x5c, 0x0a, 0x16, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72,
	0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x42,
	0x0a, 0x15, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x5f, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x52, 0x14, 0x72, 0x65,
	0x71, 0x75, 0x69, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61,
	0x63, 0x74, 0x22, 0x38, 0x0a, 0x17, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x43, 0x68,
	0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1d, 0x0a,
	0x0a, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x49, 0x64, 0x22, 0x50, 0x0a, 0x12,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x41, 0x70, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x3a, 0x0a, 0x11, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x5f, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x52, 0x10, 0x70, 0x72,
	0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x22, 0x80,
	0x01, 0x0a, 0x13, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x41, 0x70, 0x70, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x17, 0x0a, 0x06, 0x61, 0x70, 0x70, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x05, 0x61, 0x70, 0x70, 0x49, 0x64, 0x12,
	0x42, 0x0a, 0x0c, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x49, 0x6e, 0x69, 0x74,
	0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52, 0x0c, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x42, 0x0c, 0x0a, 0x0a, 0x61, 0x63, 0x6b, 0x5f, 0x6f, 0x72, 0x5f, 0x6e, 0x65,
	0x77, 0x22, 0x38, 0x0a, 0x17, 0x49, 0x6e, 0x69, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c,
	0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x0a,
	0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x49, 0x64, 0x22, 0xf2, 0x01, 0x0a, 0x12,
	0x49, 0x6e, 0x69, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x49,
	0x64, 0x12, 0x15, 0x0a, 0x06, 0x61, 0x70, 0x70, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x61, 0x70, 0x70, 0x49, 0x64, 0x12, 0x4d, 0x0a, 0x0c, 0x70, 0x61, 0x72, 0x74,
	0x69, 0x63, 0x69, 0x70, 0x61, 0x6e, 0x74, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x29,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x49, 0x6e, 0x69, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x50, 0x61, 0x72, 0x74, 0x69, 0x63, 0x69, 0x70,
	0x61, 0x6e, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x0c, 0x70, 0x61, 0x72, 0x74, 0x69,
	0x63, 0x69, 0x70, 0x61, 0x6e, 0x74, 0x73, 0x1a, 0x57, 0x0a, 0x11, 0x50, 0x61, 0x72, 0x74, 0x69,
	0x63, 0x69, 0x70, 0x61, 0x6e, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2c,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x50, 0x61, 0x72, 0x74, 0x69, 0x63,
	0x69, 0x70, 0x61, 0x6e, 0x74, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x22, 0x6a, 0x0a, 0x13, 0x49, 0x6e, 0x69, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x37, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x49, 0x6e,
	0x69, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x22, 0x1a, 0x0a, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x07, 0x0a, 0x03, 0x41, 0x43,
	0x4b, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x03, 0x45, 0x52, 0x52, 0x10, 0x01, 0x22, 0x4b, 0x0a, 0x13,
	0x53, 0x74, 0x61, 0x72, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c,
	0x49, 0x64, 0x12, 0x15, 0x0a, 0x06, 0x61, 0x70, 0x70, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x61, 0x70, 0x70, 0x49, 0x64, 0x22, 0x6c, 0x0a, 0x14, 0x53, 0x74, 0x61,
	0x72, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x38, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x20, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x43, 0x68, 0x61,
	0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x52, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x1a, 0x0a, 0x06, 0x52,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x07, 0x0a, 0x03, 0x41, 0x43, 0x4b, 0x10, 0x00, 0x12, 0x07,
	0x0a, 0x03, 0x45, 0x52, 0x52, 0x10, 0x01, 0x2a, 0x19, 0x0a, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x12, 0x06, 0x0a, 0x02, 0x4f, 0x4b, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x03, 0x45, 0x52, 0x52,
	0x10, 0x01, 0x32, 0xa5, 0x02, 0x0a, 0x11, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4d, 0x69,
	0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x12, 0x4e, 0x0a, 0x0f, 0x52, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x1b, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65,
	0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x44, 0x0a, 0x0b, 0x52, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x41, 0x70, 0x70, 0x12, 0x17, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x41, 0x70, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x18, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x41,
	0x70, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x30, 0x01, 0x12, 0x3d,
	0x0a, 0x07, 0x41, 0x70, 0x70, 0x53, 0x65, 0x6e, 0x64, 0x12, 0x1a, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x4f, 0x75, 0x74, 0x1a, 0x14, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x41, 0x70, 0x70, 0x53,
	0x65, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x3b, 0x0a,
	0x07, 0x41, 0x70, 0x70, 0x52, 0x65, 0x63, 0x76, 0x12, 0x13, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x41,
	0x70, 0x70, 0x52, 0x65, 0x63, 0x76, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x6e, 0x22, 0x00, 0x32, 0xfe, 0x01, 0x0a, 0x10, 0x50,
	0x75, 0x62, 0x6c, 0x69, 0x63, 0x4d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x77, 0x61, 0x72, 0x65, 0x12,
	0x42, 0x0a, 0x0b, 0x49, 0x6e, 0x69, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x17,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x49, 0x6e, 0x69, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x49, 0x6e,
	0x69, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x45, 0x0a, 0x0c, 0x53, 0x74, 0x61, 0x72, 0x74, 0x43, 0x68, 0x61, 0x6e,
	0x6e, 0x65, 0x6c, 0x12, 0x18, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x43,
	0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x5f, 0x0a, 0x0f, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x22, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x57, 0x69, 0x74, 0x68, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x73, 0x1a, 0x22, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x41, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x57, 0x69, 0x74, 0x68, 0x48, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x73, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x42, 0x1f, 0x5a, 0x1d, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6c, 0x70, 0x6f, 0x6d, 0x62,
	0x6f, 0x2f, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_middleware_proto_rawDescOnce sync.Once
	file_api_middleware_proto_rawDescData = file_api_middleware_proto_rawDesc
)

func file_api_middleware_proto_rawDescGZIP() []byte {
	file_api_middleware_proto_rawDescOnce.Do(func() {
		file_api_middleware_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_middleware_proto_rawDescData)
	})
	return file_api_middleware_proto_rawDescData
}

var file_api_middleware_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_api_middleware_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_api_middleware_proto_goTypes = []interface{}{
	(Result)(0),                           // 0: api.Result
	(InitChannelResponse_Result)(0),       // 1: api.InitChannelResponse.Result
	(StartChannelResponse_Result)(0),      // 2: api.StartChannelResponse.Result
	(*AppSendResponse)(nil),               // 3: api.AppSendResponse
	(*AppRecvRequest)(nil),                // 4: api.AppRecvRequest
	(*RegisterChannelRequest)(nil),        // 5: api.RegisterChannelRequest
	(*RegisterChannelResponse)(nil),       // 6: api.RegisterChannelResponse
	(*RegisterAppRequest)(nil),            // 7: api.RegisterAppRequest
	(*RegisterAppResponse)(nil),           // 8: api.RegisterAppResponse
	(*InitChannelNotification)(nil),       // 9: api.InitChannelNotification
	(*InitChannelRequest)(nil),            // 10: api.InitChannelRequest
	(*InitChannelResponse)(nil),           // 11: api.InitChannelResponse
	(*StartChannelRequest)(nil),           // 12: api.StartChannelRequest
	(*StartChannelResponse)(nil),          // 13: api.StartChannelResponse
	nil,                                   // 14: api.InitChannelRequest.ParticipantsEntry
	(*Contract)(nil),                      // 15: api.Contract
	(*RemoteParticipant)(nil),             // 16: api.RemoteParticipant
	(*ApplicationMessageOut)(nil),         // 17: api.ApplicationMessageOut
	(*ApplicationMessageWithHeaders)(nil), // 18: api.ApplicationMessageWithHeaders
	(*ApplicationMessageIn)(nil),          // 19: api.ApplicationMessageIn
}
var file_api_middleware_proto_depIdxs = []int32{
	0,  // 0: api.AppSendResponse.result:type_name -> api.Result
	15, // 1: api.RegisterChannelRequest.requirements_contract:type_name -> api.Contract
	15, // 2: api.RegisterAppRequest.provider_contract:type_name -> api.Contract
	9,  // 3: api.RegisterAppResponse.notification:type_name -> api.InitChannelNotification
	14, // 4: api.InitChannelRequest.participants:type_name -> api.InitChannelRequest.ParticipantsEntry
	1,  // 5: api.InitChannelResponse.result:type_name -> api.InitChannelResponse.Result
	2,  // 6: api.StartChannelResponse.result:type_name -> api.StartChannelResponse.Result
	16, // 7: api.InitChannelRequest.ParticipantsEntry.value:type_name -> api.RemoteParticipant
	5,  // 8: api.PrivateMiddleware.RegisterChannel:input_type -> api.RegisterChannelRequest
	7,  // 9: api.PrivateMiddleware.RegisterApp:input_type -> api.RegisterAppRequest
	17, // 10: api.PrivateMiddleware.AppSend:input_type -> api.ApplicationMessageOut
	4,  // 11: api.PrivateMiddleware.AppRecv:input_type -> api.AppRecvRequest
	10, // 12: api.PublicMiddleware.InitChannel:input_type -> api.InitChannelRequest
	12, // 13: api.PublicMiddleware.StartChannel:input_type -> api.StartChannelRequest
	18, // 14: api.PublicMiddleware.MessageExchange:input_type -> api.ApplicationMessageWithHeaders
	6,  // 15: api.PrivateMiddleware.RegisterChannel:output_type -> api.RegisterChannelResponse
	8,  // 16: api.PrivateMiddleware.RegisterApp:output_type -> api.RegisterAppResponse
	3,  // 17: api.PrivateMiddleware.AppSend:output_type -> api.AppSendResponse
	19, // 18: api.PrivateMiddleware.AppRecv:output_type -> api.ApplicationMessageIn
	11, // 19: api.PublicMiddleware.InitChannel:output_type -> api.InitChannelResponse
	13, // 20: api.PublicMiddleware.StartChannel:output_type -> api.StartChannelResponse
	18, // 21: api.PublicMiddleware.MessageExchange:output_type -> api.ApplicationMessageWithHeaders
	15, // [15:22] is the sub-list for method output_type
	8,  // [8:15] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_api_middleware_proto_init() }
func file_api_middleware_proto_init() {
	if File_api_middleware_proto != nil {
		return
	}
	file_api_app_message_proto_init()
	file_api_contracts_proto_init()
	file_api_broker_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_api_middleware_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppSendResponse); i {
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
		file_api_middleware_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppRecvRequest); i {
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
		file_api_middleware_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterChannelRequest); i {
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
		file_api_middleware_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterChannelResponse); i {
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
		file_api_middleware_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterAppRequest); i {
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
		file_api_middleware_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterAppResponse); i {
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
		file_api_middleware_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InitChannelNotification); i {
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
		file_api_middleware_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InitChannelRequest); i {
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
		file_api_middleware_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InitChannelResponse); i {
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
		file_api_middleware_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartChannelRequest); i {
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
		file_api_middleware_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartChannelResponse); i {
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
	file_api_middleware_proto_msgTypes[5].OneofWrappers = []interface{}{
		(*RegisterAppResponse_AppId)(nil),
		(*RegisterAppResponse_Notification)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_middleware_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_api_middleware_proto_goTypes,
		DependencyIndexes: file_api_middleware_proto_depIdxs,
		EnumInfos:         file_api_middleware_proto_enumTypes,
		MessageInfos:      file_api_middleware_proto_msgTypes,
	}.Build()
	File_api_middleware_proto = out.File
	file_api_middleware_proto_rawDesc = nil
	file_api_middleware_proto_goTypes = nil
	file_api_middleware_proto_depIdxs = nil
}
