// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.29.0
// 	protoc        (unknown)
// source: search/v1/contracts.proto

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

type CFSMTransition_TransitionType int32

const (
	CFSMTransition_TRANSITION_TYPE_UNSPECIFIED CFSMTransition_TransitionType = 0
	CFSMTransition_TRANSITION_TYPE_SEND        CFSMTransition_TransitionType = 1
	CFSMTransition_TRANSITION_TYPE_RECV        CFSMTransition_TransitionType = 2
)

// Enum value maps for CFSMTransition_TransitionType.
var (
	CFSMTransition_TransitionType_name = map[int32]string{
		0: "TRANSITION_TYPE_UNSPECIFIED",
		1: "TRANSITION_TYPE_SEND",
		2: "TRANSITION_TYPE_RECV",
	}
	CFSMTransition_TransitionType_value = map[string]int32{
		"TRANSITION_TYPE_UNSPECIFIED": 0,
		"TRANSITION_TYPE_SEND":        1,
		"TRANSITION_TYPE_RECV":        2,
	}
)

func (x CFSMTransition_TransitionType) Enum() *CFSMTransition_TransitionType {
	p := new(CFSMTransition_TransitionType)
	*p = x
	return p
}

func (x CFSMTransition_TransitionType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CFSMTransition_TransitionType) Descriptor() protoreflect.EnumDescriptor {
	return file_search_v1_contracts_proto_enumTypes[0].Descriptor()
}

func (CFSMTransition_TransitionType) Type() protoreflect.EnumType {
	return &file_search_v1_contracts_proto_enumTypes[0]
}

func (x CFSMTransition_TransitionType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CFSMTransition_TransitionType.Descriptor instead.
func (CFSMTransition_TransitionType) EnumDescriptor() ([]byte, []int) {
	return file_search_v1_contracts_proto_rawDescGZIP(), []int{3, 0}
}

type Contract struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Contract           string   `protobuf:"bytes,1,opt,name=contract,proto3" json:"contract,omitempty"`                                               // local participant must be called "self" or some other distinguished name
	RemoteParticipants []string `protobuf:"bytes,2,rep,name=remote_participants,json=remoteParticipants,proto3" json:"remote_participants,omitempty"` // participant identifiers (must be referenced in the contract)
}

func (x *Contract) Reset() {
	*x = Contract{}
	if protoimpl.UnsafeEnabled {
		mi := &file_search_v1_contracts_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Contract) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Contract) ProtoMessage() {}

func (x *Contract) ProtoReflect() protoreflect.Message {
	mi := &file_search_v1_contracts_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Contract.ProtoReflect.Descriptor instead.
func (*Contract) Descriptor() ([]byte, []int) {
	return file_search_v1_contracts_proto_rawDescGZIP(), []int{0}
}

func (x *Contract) GetContract() string {
	if x != nil {
		return x.Contract
	}
	return ""
}

func (x *Contract) GetRemoteParticipants() []string {
	if x != nil {
		return x.RemoteParticipants
	}
	return nil
}

type CFSM struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The first state is the start state of the CFSM
	States []*CFSMState `protobuf:"bytes,1,rep,name=states,proto3" json:"states,omitempty"`
}

func (x *CFSM) Reset() {
	*x = CFSM{}
	if protoimpl.UnsafeEnabled {
		mi := &file_search_v1_contracts_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CFSM) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CFSM) ProtoMessage() {}

func (x *CFSM) ProtoReflect() protoreflect.Message {
	mi := &file_search_v1_contracts_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CFSM.ProtoReflect.Descriptor instead.
func (*CFSM) Descriptor() ([]byte, []int) {
	return file_search_v1_contracts_proto_rawDescGZIP(), []int{1}
}

func (x *CFSM) GetStates() []*CFSMState {
	if x != nil {
		return x.States
	}
	return nil
}

type CFSMState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string            `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Transitions []*CFSMTransition `protobuf:"bytes,2,rep,name=transitions,proto3" json:"transitions,omitempty"` // outbound transitions
}

func (x *CFSMState) Reset() {
	*x = CFSMState{}
	if protoimpl.UnsafeEnabled {
		mi := &file_search_v1_contracts_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CFSMState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CFSMState) ProtoMessage() {}

func (x *CFSMState) ProtoReflect() protoreflect.Message {
	mi := &file_search_v1_contracts_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CFSMState.ProtoReflect.Descriptor instead.
func (*CFSMState) Descriptor() ([]byte, []int) {
	return file_search_v1_contracts_proto_rawDescGZIP(), []int{2}
}

func (x *CFSMState) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *CFSMState) GetTransitions() []*CFSMTransition {
	if x != nil {
		return x.Transitions
	}
	return nil
}

type CFSMTransition struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ToState      string                        `protobuf:"bytes,1,opt,name=to_state,json=toState,proto3" json:"to_state,omitempty"` // id of the CFSMState to which we transition
	Type         CFSMTransition_TransitionType `protobuf:"varint,2,opt,name=type,proto3,enum=search.v1.CFSMTransition_TransitionType" json:"type,omitempty"`
	Interlocutor string                        `protobuf:"bytes,3,opt,name=interlocutor,proto3" json:"interlocutor,omitempty"` // The CFSM with we exchange the message
	Label        string                        `protobuf:"bytes,4,opt,name=label,proto3" json:"label,omitempty"`               // The message that is sent
}

func (x *CFSMTransition) Reset() {
	*x = CFSMTransition{}
	if protoimpl.UnsafeEnabled {
		mi := &file_search_v1_contracts_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CFSMTransition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CFSMTransition) ProtoMessage() {}

func (x *CFSMTransition) ProtoReflect() protoreflect.Message {
	mi := &file_search_v1_contracts_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CFSMTransition.ProtoReflect.Descriptor instead.
func (*CFSMTransition) Descriptor() ([]byte, []int) {
	return file_search_v1_contracts_proto_rawDescGZIP(), []int{3}
}

func (x *CFSMTransition) GetToState() string {
	if x != nil {
		return x.ToState
	}
	return ""
}

func (x *CFSMTransition) GetType() CFSMTransition_TransitionType {
	if x != nil {
		return x.Type
	}
	return CFSMTransition_TRANSITION_TYPE_UNSPECIFIED
}

func (x *CFSMTransition) GetInterlocutor() string {
	if x != nil {
		return x.Interlocutor
	}
	return ""
}

func (x *CFSMTransition) GetLabel() string {
	if x != nil {
		return x.Label
	}
	return ""
}

var File_search_v1_contracts_proto protoreflect.FileDescriptor

var file_search_v1_contracts_proto_rawDesc = []byte{
	0x0a, 0x19, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6e, 0x74,
	0x72, 0x61, 0x63, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x73, 0x65, 0x61,
	0x72, 0x63, 0x68, 0x2e, 0x76, 0x31, 0x22, 0x57, 0x0a, 0x08, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61,
	0x63, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x12, 0x2f,
	0x0a, 0x13, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x5f, 0x70, 0x61, 0x72, 0x74, 0x69, 0x63, 0x69,
	0x70, 0x61, 0x6e, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x12, 0x72, 0x65, 0x6d,
	0x6f, 0x74, 0x65, 0x50, 0x61, 0x72, 0x74, 0x69, 0x63, 0x69, 0x70, 0x61, 0x6e, 0x74, 0x73, 0x22,
	0x34, 0x0a, 0x04, 0x43, 0x46, 0x53, 0x4d, 0x12, 0x2c, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x65,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x2e, 0x76, 0x31, 0x2e, 0x43, 0x46, 0x53, 0x4d, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x65, 0x73, 0x22, 0x58, 0x0a, 0x09, 0x43, 0x46, 0x53, 0x4d, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x3b, 0x0a, 0x0b, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x2e, 0x76, 0x31, 0x2e, 0x43, 0x46, 0x53, 0x4d, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x0b, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22,
	0x8a, 0x02, 0x0a, 0x0e, 0x43, 0x46, 0x53, 0x4d, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x19, 0x0a, 0x08, 0x74, 0x6f, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x74, 0x6f, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x3c, 0x0a,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x28, 0x2e, 0x73, 0x65,
	0x61, 0x72, 0x63, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x46, 0x53, 0x4d, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x6c, 0x6f, 0x63, 0x75, 0x74, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6c, 0x6f, 0x63, 0x75, 0x74, 0x6f, 0x72, 0x12,
	0x14, 0x0a, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x6c, 0x61, 0x62, 0x65, 0x6c, 0x22, 0x65, 0x0a, 0x0e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x69, 0x74,
	0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1f, 0x0a, 0x1b, 0x54, 0x52, 0x41, 0x4e, 0x53,
	0x49, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45,
	0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x18, 0x0a, 0x14, 0x54, 0x52, 0x41, 0x4e,
	0x53, 0x49, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x53, 0x45, 0x4e, 0x44,
	0x10, 0x01, 0x12, 0x18, 0x0a, 0x14, 0x54, 0x52, 0x41, 0x4e, 0x53, 0x49, 0x54, 0x49, 0x4f, 0x4e,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x52, 0x45, 0x43, 0x56, 0x10, 0x02, 0x42, 0x25, 0x5a, 0x23,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6c, 0x70, 0x6f, 0x6d,
	0x62, 0x6f, 0x2f, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f,
	0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_search_v1_contracts_proto_rawDescOnce sync.Once
	file_search_v1_contracts_proto_rawDescData = file_search_v1_contracts_proto_rawDesc
)

func file_search_v1_contracts_proto_rawDescGZIP() []byte {
	file_search_v1_contracts_proto_rawDescOnce.Do(func() {
		file_search_v1_contracts_proto_rawDescData = protoimpl.X.CompressGZIP(file_search_v1_contracts_proto_rawDescData)
	})
	return file_search_v1_contracts_proto_rawDescData
}

var file_search_v1_contracts_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_search_v1_contracts_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_search_v1_contracts_proto_goTypes = []interface{}{
	(CFSMTransition_TransitionType)(0), // 0: search.v1.CFSMTransition.TransitionType
	(*Contract)(nil),                   // 1: search.v1.Contract
	(*CFSM)(nil),                       // 2: search.v1.CFSM
	(*CFSMState)(nil),                  // 3: search.v1.CFSMState
	(*CFSMTransition)(nil),             // 4: search.v1.CFSMTransition
}
var file_search_v1_contracts_proto_depIdxs = []int32{
	3, // 0: search.v1.CFSM.states:type_name -> search.v1.CFSMState
	4, // 1: search.v1.CFSMState.transitions:type_name -> search.v1.CFSMTransition
	0, // 2: search.v1.CFSMTransition.type:type_name -> search.v1.CFSMTransition.TransitionType
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_search_v1_contracts_proto_init() }
func file_search_v1_contracts_proto_init() {
	if File_search_v1_contracts_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_search_v1_contracts_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Contract); i {
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
		file_search_v1_contracts_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CFSM); i {
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
		file_search_v1_contracts_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CFSMState); i {
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
		file_search_v1_contracts_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CFSMTransition); i {
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
			RawDescriptor: file_search_v1_contracts_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_search_v1_contracts_proto_goTypes,
		DependencyIndexes: file_search_v1_contracts_proto_depIdxs,
		EnumInfos:         file_search_v1_contracts_proto_enumTypes,
		MessageInfos:      file_search_v1_contracts_proto_msgTypes,
	}.Build()
	File_search_v1_contracts_proto = out.File
	file_search_v1_contracts_proto_rawDesc = nil
	file_search_v1_contracts_proto_goTypes = nil
	file_search_v1_contracts_proto_depIdxs = nil
}
