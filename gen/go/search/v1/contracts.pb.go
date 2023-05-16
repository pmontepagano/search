// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
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

type GlobalContractFormat int32

const (
	GlobalContractFormat_GLOBAL_CONTRACT_FORMAT_UNSPECIFIED GlobalContractFormat = 0
	GlobalContractFormat_GLOBAL_CONTRACT_FORMAT_FSA         GlobalContractFormat = 1 // System of CFSMs in FSA file format.
	GlobalContractFormat_GLOBAL_CONTRACT_FORMAT_GC          GlobalContractFormat = 2 // Global Choreography.
)

// Enum value maps for GlobalContractFormat.
var (
	GlobalContractFormat_name = map[int32]string{
		0: "GLOBAL_CONTRACT_FORMAT_UNSPECIFIED",
		1: "GLOBAL_CONTRACT_FORMAT_FSA",
		2: "GLOBAL_CONTRACT_FORMAT_GC",
	}
	GlobalContractFormat_value = map[string]int32{
		"GLOBAL_CONTRACT_FORMAT_UNSPECIFIED": 0,
		"GLOBAL_CONTRACT_FORMAT_FSA":         1,
		"GLOBAL_CONTRACT_FORMAT_GC":          2,
	}
)

func (x GlobalContractFormat) Enum() *GlobalContractFormat {
	p := new(GlobalContractFormat)
	*p = x
	return p
}

func (x GlobalContractFormat) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (GlobalContractFormat) Descriptor() protoreflect.EnumDescriptor {
	return file_search_v1_contracts_proto_enumTypes[0].Descriptor()
}

func (GlobalContractFormat) Type() protoreflect.EnumType {
	return &file_search_v1_contracts_proto_enumTypes[0]
}

func (x GlobalContractFormat) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use GlobalContractFormat.Descriptor instead.
func (GlobalContractFormat) EnumDescriptor() ([]byte, []int) {
	return file_search_v1_contracts_proto_rawDescGZIP(), []int{0}
}

type LocalContractFormat int32

const (
	LocalContractFormat_LOCAL_CONTRACT_FORMAT_UNSPECIFIED LocalContractFormat = 0
	LocalContractFormat_LOCAL_CONTRACT_FORMAT_FSA         LocalContractFormat = 1 // Single CFSM in FSA file format (for Service Providers).
)

// Enum value maps for LocalContractFormat.
var (
	LocalContractFormat_name = map[int32]string{
		0: "LOCAL_CONTRACT_FORMAT_UNSPECIFIED",
		1: "LOCAL_CONTRACT_FORMAT_FSA",
	}
	LocalContractFormat_value = map[string]int32{
		"LOCAL_CONTRACT_FORMAT_UNSPECIFIED": 0,
		"LOCAL_CONTRACT_FORMAT_FSA":         1,
	}
)

func (x LocalContractFormat) Enum() *LocalContractFormat {
	p := new(LocalContractFormat)
	*p = x
	return p
}

func (x LocalContractFormat) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (LocalContractFormat) Descriptor() protoreflect.EnumDescriptor {
	return file_search_v1_contracts_proto_enumTypes[1].Descriptor()
}

func (LocalContractFormat) Type() protoreflect.EnumType {
	return &file_search_v1_contracts_proto_enumTypes[1]
}

func (x LocalContractFormat) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use LocalContractFormat.Descriptor instead.
func (LocalContractFormat) EnumDescriptor() ([]byte, []int) {
	return file_search_v1_contracts_proto_rawDescGZIP(), []int{1}
}

type GlobalContract struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Contract []byte               `protobuf:"bytes,1,opt,name=contract,proto3" json:"contract,omitempty"`
	Format   GlobalContractFormat `protobuf:"varint,2,opt,name=format,proto3,enum=search.v1.GlobalContractFormat" json:"format,omitempty"` // string initiator_name = 3;
}

func (x *GlobalContract) Reset() {
	*x = GlobalContract{}
	if protoimpl.UnsafeEnabled {
		mi := &file_search_v1_contracts_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GlobalContract) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GlobalContract) ProtoMessage() {}

func (x *GlobalContract) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use GlobalContract.ProtoReflect.Descriptor instead.
func (*GlobalContract) Descriptor() ([]byte, []int) {
	return file_search_v1_contracts_proto_rawDescGZIP(), []int{0}
}

func (x *GlobalContract) GetContract() []byte {
	if x != nil {
		return x.Contract
	}
	return nil
}

func (x *GlobalContract) GetFormat() GlobalContractFormat {
	if x != nil {
		return x.Format
	}
	return GlobalContractFormat_GLOBAL_CONTRACT_FORMAT_UNSPECIFIED
}

type LocalContract struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Contract []byte              `protobuf:"bytes,1,opt,name=contract,proto3" json:"contract,omitempty"`
	Format   LocalContractFormat `protobuf:"varint,2,opt,name=format,proto3,enum=search.v1.LocalContractFormat" json:"format,omitempty"`
}

func (x *LocalContract) Reset() {
	*x = LocalContract{}
	if protoimpl.UnsafeEnabled {
		mi := &file_search_v1_contracts_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LocalContract) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocalContract) ProtoMessage() {}

func (x *LocalContract) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use LocalContract.ProtoReflect.Descriptor instead.
func (*LocalContract) Descriptor() ([]byte, []int) {
	return file_search_v1_contracts_proto_rawDescGZIP(), []int{1}
}

func (x *LocalContract) GetContract() []byte {
	if x != nil {
		return x.Contract
	}
	return nil
}

func (x *LocalContract) GetFormat() LocalContractFormat {
	if x != nil {
		return x.Format
	}
	return LocalContractFormat_LOCAL_CONTRACT_FORMAT_UNSPECIFIED
}

var File_search_v1_contracts_proto protoreflect.FileDescriptor

var file_search_v1_contracts_proto_rawDesc = []byte{
	0x0a, 0x19, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6e, 0x74,
	0x72, 0x61, 0x63, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x73, 0x65, 0x61,
	0x72, 0x63, 0x68, 0x2e, 0x76, 0x31, 0x22, 0x65, 0x0a, 0x0e, 0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c,
	0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6f, 0x6e, 0x74,
	0x72, 0x61, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x63, 0x6f, 0x6e, 0x74,
	0x72, 0x61, 0x63, 0x74, 0x12, 0x37, 0x0a, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x1f, 0x2e, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2e, 0x76, 0x31,
	0x2e, 0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x46,
	0x6f, 0x72, 0x6d, 0x61, 0x74, 0x52, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x22, 0x63, 0x0a,
	0x0d, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x12, 0x1a,
	0x0a, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x12, 0x36, 0x0a, 0x06, 0x66, 0x6f,
	0x72, 0x6d, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x73, 0x65, 0x61,
	0x72, 0x63, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x43, 0x6f, 0x6e, 0x74,
	0x72, 0x61, 0x63, 0x74, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x52, 0x06, 0x66, 0x6f, 0x72, 0x6d,
	0x61, 0x74, 0x2a, 0x7d, 0x0a, 0x14, 0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x43, 0x6f, 0x6e, 0x74,
	0x72, 0x61, 0x63, 0x74, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x26, 0x0a, 0x22, 0x47, 0x4c,
	0x4f, 0x42, 0x41, 0x4c, 0x5f, 0x43, 0x4f, 0x4e, 0x54, 0x52, 0x41, 0x43, 0x54, 0x5f, 0x46, 0x4f,
	0x52, 0x4d, 0x41, 0x54, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44,
	0x10, 0x00, 0x12, 0x1e, 0x0a, 0x1a, 0x47, 0x4c, 0x4f, 0x42, 0x41, 0x4c, 0x5f, 0x43, 0x4f, 0x4e,
	0x54, 0x52, 0x41, 0x43, 0x54, 0x5f, 0x46, 0x4f, 0x52, 0x4d, 0x41, 0x54, 0x5f, 0x46, 0x53, 0x41,
	0x10, 0x01, 0x12, 0x1d, 0x0a, 0x19, 0x47, 0x4c, 0x4f, 0x42, 0x41, 0x4c, 0x5f, 0x43, 0x4f, 0x4e,
	0x54, 0x52, 0x41, 0x43, 0x54, 0x5f, 0x46, 0x4f, 0x52, 0x4d, 0x41, 0x54, 0x5f, 0x47, 0x43, 0x10,
	0x02, 0x2a, 0x5b, 0x0a, 0x13, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61,
	0x63, 0x74, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x25, 0x0a, 0x21, 0x4c, 0x4f, 0x43, 0x41,
	0x4c, 0x5f, 0x43, 0x4f, 0x4e, 0x54, 0x52, 0x41, 0x43, 0x54, 0x5f, 0x46, 0x4f, 0x52, 0x4d, 0x41,
	0x54, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12,
	0x1d, 0x0a, 0x19, 0x4c, 0x4f, 0x43, 0x41, 0x4c, 0x5f, 0x43, 0x4f, 0x4e, 0x54, 0x52, 0x41, 0x43,
	0x54, 0x5f, 0x46, 0x4f, 0x52, 0x4d, 0x41, 0x54, 0x5f, 0x46, 0x53, 0x41, 0x10, 0x01, 0x42, 0x25,
	0x5a, 0x23, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6c, 0x70,
	0x6f, 0x6d, 0x62, 0x6f, 0x2f, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2f, 0x67, 0x65, 0x6e, 0x2f,
	0x67, 0x6f, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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

var file_search_v1_contracts_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_search_v1_contracts_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_search_v1_contracts_proto_goTypes = []interface{}{
	(GlobalContractFormat)(0), // 0: search.v1.GlobalContractFormat
	(LocalContractFormat)(0),  // 1: search.v1.LocalContractFormat
	(*GlobalContract)(nil),    // 2: search.v1.GlobalContract
	(*LocalContract)(nil),     // 3: search.v1.LocalContract
}
var file_search_v1_contracts_proto_depIdxs = []int32{
	0, // 0: search.v1.GlobalContract.format:type_name -> search.v1.GlobalContractFormat
	1, // 1: search.v1.LocalContract.format:type_name -> search.v1.LocalContractFormat
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_search_v1_contracts_proto_init() }
func file_search_v1_contracts_proto_init() {
	if File_search_v1_contracts_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_search_v1_contracts_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GlobalContract); i {
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
			switch v := v.(*LocalContract); i {
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
			NumEnums:      2,
			NumMessages:   2,
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
