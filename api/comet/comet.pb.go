// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.19.1
// source: comet/comet.proto

package comet

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

type Type int32

const (
	Type_AREA  Type = 0
	Type_ROOM  Type = 1
	Type_PUSH  Type = 2
	Type_CLOSE Type = 3
)

// Enum value maps for Type.
var (
	Type_name = map[int32]string{
		0: "AREA",
		1: "ROOM",
		2: "PUSH",
		3: "CLOSE",
	}
	Type_value = map[string]int32{
		"AREA":  0,
		"ROOM":  1,
		"PUSH":  2,
		"CLOSE": 3,
	}
)

func (x Type) Enum() *Type {
	p := new(Type)
	*p = x
	return p
}

func (x Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Type) Descriptor() protoreflect.EnumDescriptor {
	return file_comet_comet_proto_enumTypes[0].Descriptor()
}

func (Type) Type() protoreflect.EnumType {
	return &file_comet_comet_proto_enumTypes[0]
}

func (x Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Type.Descriptor instead.
func (Type) EnumDescriptor() ([]byte, []int) {
	return file_comet_comet_proto_rawDescGZIP(), []int{0}
}

// todo 可单独拆分job层 目前为每个comet绑定对应的kafka topic
type Msg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type   Type   `protobuf:"varint,1,opt,name=type,proto3,enum=comet.Type" json:"type,omitempty"`
	ToId   uint64 `protobuf:"varint,2,opt,name=to_id,json=toId,proto3" json:"to_id,omitempty"`
	SendId uint64 `protobuf:"varint,3,opt,name=send_id,json=sendId,proto3" json:"send_id,omitempty"`
	Msg    []byte `protobuf:"bytes,4,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *Msg) Reset() {
	*x = Msg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_comet_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Msg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Msg) ProtoMessage() {}

func (x *Msg) ProtoReflect() protoreflect.Message {
	mi := &file_comet_comet_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Msg.ProtoReflect.Descriptor instead.
func (*Msg) Descriptor() ([]byte, []int) {
	return file_comet_comet_proto_rawDescGZIP(), []int{0}
}

func (x *Msg) GetType() Type {
	if x != nil {
		return x.Type
	}
	return Type_AREA
}

func (x *Msg) GetToId() uint64 {
	if x != nil {
		return x.ToId
	}
	return 0
}

func (x *Msg) GetSendId() uint64 {
	if x != nil {
		return x.SendId
	}
	return 0
}

func (x *Msg) GetMsg() []byte {
	if x != nil {
		return x.Msg
	}
	return nil
}

var File_comet_comet_proto protoreflect.FileDescriptor

var file_comet_comet_proto_rawDesc = []byte{
	0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x2f, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x22, 0x66, 0x0a, 0x03, 0x4d, 0x73,
	0x67, 0x12, 0x1f, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x0b, 0x2e, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x12, 0x13, 0x0a, 0x05, 0x74, 0x6f, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x04, 0x74, 0x6f, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x73, 0x65, 0x6e, 0x64, 0x5f,
	0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x49, 0x64,
	0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x6d,
	0x73, 0x67, 0x2a, 0x2f, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x08, 0x0a, 0x04, 0x41, 0x52,
	0x45, 0x41, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x52, 0x4f, 0x4f, 0x4d, 0x10, 0x01, 0x12, 0x08,
	0x0a, 0x04, 0x50, 0x55, 0x53, 0x48, 0x10, 0x02, 0x12, 0x09, 0x0a, 0x05, 0x43, 0x4c, 0x4f, 0x53,
	0x45, 0x10, 0x03, 0x42, 0x29, 0x5a, 0x27, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x32, 0x70, 0x67, 0x63, 0x6e, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x69, 0x6d, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x3b, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_comet_comet_proto_rawDescOnce sync.Once
	file_comet_comet_proto_rawDescData = file_comet_comet_proto_rawDesc
)

func file_comet_comet_proto_rawDescGZIP() []byte {
	file_comet_comet_proto_rawDescOnce.Do(func() {
		file_comet_comet_proto_rawDescData = protoimpl.X.CompressGZIP(file_comet_comet_proto_rawDescData)
	})
	return file_comet_comet_proto_rawDescData
}

var file_comet_comet_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_comet_comet_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_comet_comet_proto_goTypes = []interface{}{
	(Type)(0),   // 0: comet.Type
	(*Msg)(nil), // 1: comet.Msg
}
var file_comet_comet_proto_depIdxs = []int32{
	0, // 0: comet.Msg.type:type_name -> comet.Type
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_comet_comet_proto_init() }
func file_comet_comet_proto_init() {
	if File_comet_comet_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_comet_comet_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Msg); i {
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
			RawDescriptor: file_comet_comet_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_comet_comet_proto_goTypes,
		DependencyIndexes: file_comet_comet_proto_depIdxs,
		EnumInfos:         file_comet_comet_proto_enumTypes,
		MessageInfos:      file_comet_comet_proto_msgTypes,
	}.Build()
	File_comet_comet_proto = out.File
	file_comet_comet_proto_rawDesc = nil
	file_comet_comet_proto_goTypes = nil
	file_comet_comet_proto_depIdxs = nil
}
