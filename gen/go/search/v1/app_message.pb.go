// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        (unknown)
// source: search/v1/app_message.proto

package v1

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

// This is what will be exchanged between middlewares
type MessageExchangeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChannelId string `protobuf:"bytes,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"` // We'll use UUIDv4. It's a global ID shared by all participants
	// This is necessary because URLs don't univocally determine apps. There can be multiple applications
	// behind the same middleware (there is a 1:1 mapping between URLs and middlewares)
	SenderId    string      `protobuf:"bytes,2,opt,name=sender_id,json=senderId,proto3" json:"sender_id,omitempty"`          // appid de app emisora
	RecipientId string      `protobuf:"bytes,3,opt,name=recipient_id,json=recipientId,proto3" json:"recipient_id,omitempty"` // appid de app receptora
	Content     *AppMessage `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *MessageExchangeRequest) Reset() {
	*x = MessageExchangeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_search_v1_app_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessageExchangeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageExchangeRequest) ProtoMessage() {}

func (x *MessageExchangeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_search_v1_app_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageExchangeRequest.ProtoReflect.Descriptor instead.
func (*MessageExchangeRequest) Descriptor() ([]byte, []int) {
	return file_search_v1_app_message_proto_rawDescGZIP(), []int{0}
}

func (x *MessageExchangeRequest) GetChannelId() string {
	if x != nil {
		return x.ChannelId
	}
	return ""
}

func (x *MessageExchangeRequest) GetSenderId() string {
	if x != nil {
		return x.SenderId
	}
	return ""
}

func (x *MessageExchangeRequest) GetRecipientId() string {
	if x != nil {
		return x.RecipientId
	}
	return ""
}

func (x *MessageExchangeRequest) GetContent() *AppMessage {
	if x != nil {
		return x.Content
	}
	return nil
}

// This is what will be sent from an app to the middleware
type AppSendRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChannelId string      `protobuf:"bytes,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	Recipient string      `protobuf:"bytes,2,opt,name=recipient,proto3" json:"recipient,omitempty"` // name of the recipient in the local contract
	Message   *AppMessage `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *AppSendRequest) Reset() {
	*x = AppSendRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_search_v1_app_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppSendRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppSendRequest) ProtoMessage() {}

func (x *AppSendRequest) ProtoReflect() protoreflect.Message {
	mi := &file_search_v1_app_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppSendRequest.ProtoReflect.Descriptor instead.
func (*AppSendRequest) Descriptor() ([]byte, []int) {
	return file_search_v1_app_message_proto_rawDescGZIP(), []int{1}
}

func (x *AppSendRequest) GetChannelId() string {
	if x != nil {
		return x.ChannelId
	}
	return ""
}

func (x *AppSendRequest) GetRecipient() string {
	if x != nil {
		return x.Recipient
	}
	return ""
}

func (x *AppSendRequest) GetMessage() *AppMessage {
	if x != nil {
		return x.Message
	}
	return nil
}

// This is what will be sent from the middleware to a local app
type AppRecvResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChannelId string      `protobuf:"bytes,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	Sender    string      `protobuf:"bytes,2,opt,name=sender,proto3" json:"sender,omitempty"` // name of the sender in the local contract
	Message   *AppMessage `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *AppRecvResponse) Reset() {
	*x = AppRecvResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_search_v1_app_message_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppRecvResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppRecvResponse) ProtoMessage() {}

func (x *AppRecvResponse) ProtoReflect() protoreflect.Message {
	mi := &file_search_v1_app_message_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppRecvResponse.ProtoReflect.Descriptor instead.
func (*AppRecvResponse) Descriptor() ([]byte, []int) {
	return file_search_v1_app_message_proto_rawDescGZIP(), []int{2}
}

func (x *AppRecvResponse) GetChannelId() string {
	if x != nil {
		return x.ChannelId
	}
	return ""
}

func (x *AppRecvResponse) GetSender() string {
	if x != nil {
		return x.Sender
	}
	return ""
}

func (x *AppRecvResponse) GetMessage() *AppMessage {
	if x != nil {
		return x.Message
	}
	return nil
}

// This is the message content that is sent by the app (this is copied as-is by the middlewares)
// TODO: we may want to use self-describing messages to have a rich type system for messages!
// https://protobuf.dev/programming-guides/techniques/#self-description
type AppMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Body []byte `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *AppMessage) Reset() {
	*x = AppMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_search_v1_app_message_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppMessage) ProtoMessage() {}

func (x *AppMessage) ProtoReflect() protoreflect.Message {
	mi := &file_search_v1_app_message_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppMessage.ProtoReflect.Descriptor instead.
func (*AppMessage) Descriptor() ([]byte, []int) {
	return file_search_v1_app_message_proto_rawDescGZIP(), []int{3}
}

func (x *AppMessage) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *AppMessage) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

var File_search_v1_app_message_proto protoreflect.FileDescriptor

var file_search_v1_app_message_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x70, 0x70, 0x5f,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x73,
	0x65, 0x61, 0x72, 0x63, 0x68, 0x2e, 0x76, 0x31, 0x22, 0xa8, 0x01, 0x0a, 0x16, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c,
	0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x21, 0x0a, 0x0c, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74,
	0x49, 0x64, 0x12, 0x2f, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2e, 0x76, 0x31, 0x2e,
	0x41, 0x70, 0x70, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x22, 0x7e, 0x0a, 0x0e, 0x41, 0x70, 0x70, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x68, 0x61, 0x6e, 0x6e,
	0x65, 0x6c, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65,
	0x6e, 0x74, 0x12, 0x2f, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2e, 0x76, 0x31, 0x2e,
	0x41, 0x70, 0x70, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x22, 0x79, 0x0a, 0x0f, 0x41, 0x70, 0x70, 0x52, 0x65, 0x63, 0x76, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65,
	0x6c, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x68, 0x61, 0x6e,
	0x6e, 0x65, 0x6c, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x2f, 0x0a,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15,
	0x2e, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x70, 0x70, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x34,
	0x0a, 0x0a, 0x41, 0x70, 0x70, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04,
	0x62, 0x6f, 0x64, 0x79, 0x42, 0x48, 0x0a, 0x1c, 0x61, 0x72, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x6d,
	0x6f, 0x6e, 0x74, 0x65, 0x70, 0x61, 0x67, 0x61, 0x6e, 0x6f, 0x2e, 0x73, 0x65, 0x61, 0x72, 0x63,
	0x68, 0x2e, 0x76, 0x31, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x70, 0x6d, 0x6f, 0x6e, 0x74, 0x65, 0x70, 0x61, 0x67, 0x61, 0x6e, 0x6f, 0x2f, 0x73, 0x65,
	0x61, 0x72, 0x63, 0x68, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x76, 0x31, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_search_v1_app_message_proto_rawDescOnce sync.Once
	file_search_v1_app_message_proto_rawDescData = file_search_v1_app_message_proto_rawDesc
)

func file_search_v1_app_message_proto_rawDescGZIP() []byte {
	file_search_v1_app_message_proto_rawDescOnce.Do(func() {
		file_search_v1_app_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_search_v1_app_message_proto_rawDescData)
	})
	return file_search_v1_app_message_proto_rawDescData
}

var file_search_v1_app_message_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_search_v1_app_message_proto_goTypes = []interface{}{
	(*MessageExchangeRequest)(nil), // 0: search.v1.MessageExchangeRequest
	(*AppSendRequest)(nil),         // 1: search.v1.AppSendRequest
	(*AppRecvResponse)(nil),        // 2: search.v1.AppRecvResponse
	(*AppMessage)(nil),             // 3: search.v1.AppMessage
}
var file_search_v1_app_message_proto_depIdxs = []int32{
	3, // 0: search.v1.MessageExchangeRequest.content:type_name -> search.v1.AppMessage
	3, // 1: search.v1.AppSendRequest.message:type_name -> search.v1.AppMessage
	3, // 2: search.v1.AppRecvResponse.message:type_name -> search.v1.AppMessage
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_search_v1_app_message_proto_init() }
func file_search_v1_app_message_proto_init() {
	if File_search_v1_app_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_search_v1_app_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessageExchangeRequest); i {
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
		file_search_v1_app_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppSendRequest); i {
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
		file_search_v1_app_message_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppRecvResponse); i {
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
		file_search_v1_app_message_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppMessage); i {
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
			RawDescriptor: file_search_v1_app_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_search_v1_app_message_proto_goTypes,
		DependencyIndexes: file_search_v1_app_message_proto_depIdxs,
		MessageInfos:      file_search_v1_app_message_proto_msgTypes,
	}.Build()
	File_search_v1_app_message_proto = out.File
	file_search_v1_app_message_proto_rawDesc = nil
	file_search_v1_app_message_proto_goTypes = nil
	file_search_v1_app_message_proto_depIdxs = nil
}
