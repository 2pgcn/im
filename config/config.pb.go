// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.19.1
// source: config/config.proto

package config

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

type QueueType int32

const (
	QueueType_REDIS QueueType = 0
	QueueType_KAFKA QueueType = 1
)

// Enum value maps for QueueType.
var (
	QueueType_name = map[int32]string{
		0: "REDIS",
		1: "KAFKA",
	}
	QueueType_value = map[string]int32{
		"REDIS": 0,
		"KAFKA": 1,
	}
)

func (x QueueType) Enum() *QueueType {
	p := new(QueueType)
	*p = x
	return p
}

func (x QueueType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (QueueType) Descriptor() protoreflect.EnumDescriptor {
	return file_config_config_proto_enumTypes[0].Descriptor()
}

func (QueueType) Type() protoreflect.EnumType {
	return &file_config_config_proto_enumTypes[0]
}

func (x QueueType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use QueueType.Descriptor instead.
func (QueueType) EnumDescriptor() ([]byte, []int) {
	return file_config_config_proto_rawDescGZIP(), []int{0}
}

// repeated string topic = 1;
//
//	string group = 2;
//	repeated string brokers = 3;
//	int32 partition = 4;
type CometConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Server          *CometConfigServerMsg       `protobuf:"bytes,1,opt,name=server,proto3" json:"server,omitempty"`
	Queue           *CometConfigQueueMsg        `protobuf:"bytes,2,opt,name=queue,proto3" json:"queue,omitempty"`
	Ssl             *CometConfigSslMsg          `protobuf:"bytes,3,opt,name=ssl,proto3" json:"ssl,omitempty"`
	Tcp             *CometConfigTcpMsg          `protobuf:"bytes,4,opt,name=tcp,proto3" json:"tcp,omitempty"`
	LogicClientGrpc *CometConfigLogicClientGrpc `protobuf:"bytes,5,opt,name=logic_client_grpc,json=logicClientGrpc,proto3" json:"logic_client_grpc,omitempty"`
	IsDev           bool                        `protobuf:"varint,6,opt,name=is_dev,json=isDev,proto3" json:"is_dev,omitempty"`
}

func (x *CometConfig) Reset() {
	*x = CometConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CometConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CometConfig) ProtoMessage() {}

func (x *CometConfig) ProtoReflect() protoreflect.Message {
	mi := &file_config_config_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CometConfig.ProtoReflect.Descriptor instead.
func (*CometConfig) Descriptor() ([]byte, []int) {
	return file_config_config_proto_rawDescGZIP(), []int{0}
}

func (x *CometConfig) GetServer() *CometConfigServerMsg {
	if x != nil {
		return x.Server
	}
	return nil
}

func (x *CometConfig) GetQueue() *CometConfigQueueMsg {
	if x != nil {
		return x.Queue
	}
	return nil
}

func (x *CometConfig) GetSsl() *CometConfigSslMsg {
	if x != nil {
		return x.Ssl
	}
	return nil
}

func (x *CometConfig) GetTcp() *CometConfigTcpMsg {
	if x != nil {
		return x.Tcp
	}
	return nil
}

func (x *CometConfig) GetLogicClientGrpc() *CometConfigLogicClientGrpc {
	if x != nil {
		return x.LogicClientGrpc
	}
	return nil
}

func (x *CometConfig) GetIsDev() bool {
	if x != nil {
		return x.IsDev
	}
	return false
}

type LogicConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *LogicConfig) Reset() {
	*x = LogicConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_config_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogicConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogicConfig) ProtoMessage() {}

func (x *LogicConfig) ProtoReflect() protoreflect.Message {
	mi := &file_config_config_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogicConfig.ProtoReflect.Descriptor instead.
func (*LogicConfig) Descriptor() ([]byte, []int) {
	return file_config_config_proto_rawDescGZIP(), []int{1}
}

type CometConfigServerMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name           string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Addrs          []string `protobuf:"bytes,2,rep,name=addrs,proto3" json:"addrs,omitempty"`
	Port           int32    `protobuf:"varint,3,opt,name=port,proto3" json:"port,omitempty"`
	MaxConnections uint32   `protobuf:"varint,4,opt,name=max_connections,json=maxConnections,proto3" json:"max_connections,omitempty"` //max conn nums
	Goroutines     uint32   `protobuf:"varint,5,opt,name=goroutines,proto3" json:"goroutines,omitempty"`                               //go routines nums default is 0 means use cpu nums
}

func (x *CometConfigServerMsg) Reset() {
	*x = CometConfigServerMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_config_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CometConfigServerMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CometConfigServerMsg) ProtoMessage() {}

func (x *CometConfigServerMsg) ProtoReflect() protoreflect.Message {
	mi := &file_config_config_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CometConfigServerMsg.ProtoReflect.Descriptor instead.
func (*CometConfigServerMsg) Descriptor() ([]byte, []int) {
	return file_config_config_proto_rawDescGZIP(), []int{0, 0}
}

func (x *CometConfigServerMsg) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CometConfigServerMsg) GetAddrs() []string {
	if x != nil {
		return x.Addrs
	}
	return nil
}

func (x *CometConfigServerMsg) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

func (x *CometConfigServerMsg) GetMaxConnections() uint32 {
	if x != nil {
		return x.MaxConnections
	}
	return 0
}

func (x *CometConfigServerMsg) GetGoroutines() uint32 {
	if x != nil {
		return x.Goroutines
	}
	return 0
}

// todo 目前仅为kafka
type CometConfigQueueMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type  string                     `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Redis *CometConfigQueueMsg_Redis `protobuf:"bytes,2,opt,name=redis,proto3" json:"redis,omitempty"`
}

func (x *CometConfigQueueMsg) Reset() {
	*x = CometConfigQueueMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_config_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CometConfigQueueMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CometConfigQueueMsg) ProtoMessage() {}

func (x *CometConfigQueueMsg) ProtoReflect() protoreflect.Message {
	mi := &file_config_config_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CometConfigQueueMsg.ProtoReflect.Descriptor instead.
func (*CometConfigQueueMsg) Descriptor() ([]byte, []int) {
	return file_config_config_proto_rawDescGZIP(), []int{0, 1}
}

func (x *CometConfigQueueMsg) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *CometConfigQueueMsg) GetRedis() *CometConfigQueueMsg_Redis {
	if x != nil {
		return x.Redis
	}
	return nil
}

type CometConfigSslMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Enable bool   `protobuf:"varint,1,opt,name=enable,proto3" json:"enable,omitempty"`
	Key    string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Cert   string `protobuf:"bytes,3,opt,name=cert,proto3" json:"cert,omitempty"`
}

func (x *CometConfigSslMsg) Reset() {
	*x = CometConfigSslMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_config_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CometConfigSslMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CometConfigSslMsg) ProtoMessage() {}

func (x *CometConfigSslMsg) ProtoReflect() protoreflect.Message {
	mi := &file_config_config_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CometConfigSslMsg.ProtoReflect.Descriptor instead.
func (*CometConfigSslMsg) Descriptor() ([]byte, []int) {
	return file_config_config_proto_rawDescGZIP(), []int{0, 2}
}

func (x *CometConfigSslMsg) GetEnable() bool {
	if x != nil {
		return x.Enable
	}
	return false
}

func (x *CometConfigSslMsg) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *CometConfigSslMsg) GetCert() string {
	if x != nil {
		return x.Cert
	}
	return ""
}

type CometConfigTcpMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SendBuf   uint64 `protobuf:"varint,1,opt,name=send_buf,json=sendBuf,proto3" json:"send_buf,omitempty"`
	RecvBuf   uint64 `protobuf:"varint,2,opt,name=recv_buf,json=recvBuf,proto3" json:"recv_buf,omitempty"`
	Keepalive bool   `protobuf:"varint,3,opt,name=keepalive,proto3" json:"keepalive,omitempty"`
}

func (x *CometConfigTcpMsg) Reset() {
	*x = CometConfigTcpMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_config_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CometConfigTcpMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CometConfigTcpMsg) ProtoMessage() {}

func (x *CometConfigTcpMsg) ProtoReflect() protoreflect.Message {
	mi := &file_config_config_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CometConfigTcpMsg.ProtoReflect.Descriptor instead.
func (*CometConfigTcpMsg) Descriptor() ([]byte, []int) {
	return file_config_config_proto_rawDescGZIP(), []int{0, 3}
}

func (x *CometConfigTcpMsg) GetSendBuf() uint64 {
	if x != nil {
		return x.SendBuf
	}
	return 0
}

func (x *CometConfigTcpMsg) GetRecvBuf() uint64 {
	if x != nil {
		return x.RecvBuf
	}
	return 0
}

func (x *CometConfigTcpMsg) GetKeepalive() bool {
	if x != nil {
		return x.Keepalive
	}
	return false
}

type CometConfigLogicClientGrpc struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Addr    string `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	Timeout uint32 `protobuf:"varint,2,opt,name=timeout,proto3" json:"timeout,omitempty"`
}

func (x *CometConfigLogicClientGrpc) Reset() {
	*x = CometConfigLogicClientGrpc{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_config_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CometConfigLogicClientGrpc) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CometConfigLogicClientGrpc) ProtoMessage() {}

func (x *CometConfigLogicClientGrpc) ProtoReflect() protoreflect.Message {
	mi := &file_config_config_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CometConfigLogicClientGrpc.ProtoReflect.Descriptor instead.
func (*CometConfigLogicClientGrpc) Descriptor() ([]byte, []int) {
	return file_config_config_proto_rawDescGZIP(), []int{0, 4}
}

func (x *CometConfigLogicClientGrpc) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *CometConfigLogicClientGrpc) GetTimeout() uint32 {
	if x != nil {
		return x.Timeout
	}
	return 0
}

type CometConfigQueueMsg_Redis struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Addr     string `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	Db       int32  `protobuf:"varint,3,opt,name=db,proto3" json:"db,omitempty"`
}

func (x *CometConfigQueueMsg_Redis) Reset() {
	*x = CometConfigQueueMsg_Redis{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_config_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CometConfigQueueMsg_Redis) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CometConfigQueueMsg_Redis) ProtoMessage() {}

func (x *CometConfigQueueMsg_Redis) ProtoReflect() protoreflect.Message {
	mi := &file_config_config_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CometConfigQueueMsg_Redis.ProtoReflect.Descriptor instead.
func (*CometConfigQueueMsg_Redis) Descriptor() ([]byte, []int) {
	return file_config_config_proto_rawDescGZIP(), []int{0, 1, 0}
}

func (x *CometConfigQueueMsg_Redis) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *CometConfigQueueMsg_Redis) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *CometConfigQueueMsg_Redis) GetDb() int32 {
	if x != nil {
		return x.Db
	}
	return 0
}

type LogicConfigGrpcServer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Host string `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
	Port int32  `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *LogicConfigGrpcServer) Reset() {
	*x = LogicConfigGrpcServer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_config_config_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogicConfigGrpcServer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogicConfigGrpcServer) ProtoMessage() {}

func (x *LogicConfigGrpcServer) ProtoReflect() protoreflect.Message {
	mi := &file_config_config_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogicConfigGrpcServer.ProtoReflect.Descriptor instead.
func (*LogicConfigGrpcServer) Descriptor() ([]byte, []int) {
	return file_config_config_proto_rawDescGZIP(), []int{1, 0}
}

func (x *LogicConfigGrpcServer) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *LogicConfigGrpcServer) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

var File_config_config_proto protoreflect.FileDescriptor

var file_config_config_proto_rawDesc = []byte{
	0x0a, 0x13, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0xdc, 0x06,
	0x0a, 0x0b, 0x43, 0x6f, 0x6d, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x35, 0x0a,
	0x06, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x43, 0x6f, 0x6d, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4d, 0x73, 0x67, 0x52, 0x06, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x12, 0x32, 0x0a, 0x05, 0x71, 0x75, 0x65, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x43, 0x6f, 0x6d,
	0x65, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x4d, 0x73,
	0x67, 0x52, 0x05, 0x71, 0x75, 0x65, 0x75, 0x65, 0x12, 0x2c, 0x0a, 0x03, 0x73, 0x73, 0x6c, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x43,
	0x6f, 0x6d, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x73, 0x73, 0x6c, 0x4d, 0x73,
	0x67, 0x52, 0x03, 0x73, 0x73, 0x6c, 0x12, 0x2c, 0x0a, 0x03, 0x74, 0x63, 0x70, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x43, 0x6f, 0x6d,
	0x65, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x74, 0x63, 0x70, 0x4d, 0x73, 0x67, 0x52,
	0x03, 0x74, 0x63, 0x70, 0x12, 0x4f, 0x0a, 0x11, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x5f, 0x63, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x5f, 0x67, 0x72, 0x70, 0x63, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x23, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x43, 0x6f, 0x6d, 0x65, 0x74, 0x43, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x2e, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x47, 0x72, 0x70, 0x63, 0x52, 0x0f, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x43, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x47, 0x72, 0x70, 0x63, 0x12, 0x15, 0x0a, 0x06, 0x69, 0x73, 0x5f, 0x64, 0x65, 0x76, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x69, 0x73, 0x44, 0x65, 0x76, 0x1a, 0x92, 0x01, 0x0a,
	0x09, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4d, 0x73, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x61, 0x64, 0x64, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x61,
	0x64, 0x64, 0x72, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x27, 0x0a, 0x0f, 0x6d, 0x61, 0x78, 0x5f,
	0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x0e, 0x6d, 0x61, 0x78, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x67, 0x6f, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x65, 0x73, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x67, 0x6f, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x65,
	0x73, 0x1a, 0xa1, 0x01, 0x0a, 0x08, 0x71, 0x75, 0x65, 0x75, 0x65, 0x4d, 0x73, 0x67, 0x12, 0x12,
	0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x12, 0x38, 0x0a, 0x05, 0x72, 0x65, 0x64, 0x69, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x22, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x43, 0x6f, 0x6d, 0x65, 0x74,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x4d, 0x73, 0x67, 0x2e,
	0x52, 0x65, 0x64, 0x69, 0x73, 0x52, 0x05, 0x72, 0x65, 0x64, 0x69, 0x73, 0x1a, 0x47, 0x0a, 0x05,
	0x52, 0x65, 0x64, 0x69, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x64, 0x64, 0x72, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73,
	0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73,
	0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x64, 0x62, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x02, 0x64, 0x62, 0x1a, 0x46, 0x0a, 0x06, 0x73, 0x73, 0x6c, 0x4d, 0x73, 0x67, 0x12,
	0x16, 0x0a, 0x06, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x06, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x65, 0x72,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x65, 0x72, 0x74, 0x1a, 0x5c, 0x0a,
	0x06, 0x74, 0x63, 0x70, 0x4d, 0x73, 0x67, 0x12, 0x19, 0x0a, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x5f,
	0x62, 0x75, 0x66, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x73, 0x65, 0x6e, 0x64, 0x42,
	0x75, 0x66, 0x12, 0x19, 0x0a, 0x08, 0x72, 0x65, 0x63, 0x76, 0x5f, 0x62, 0x75, 0x66, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x72, 0x65, 0x63, 0x76, 0x42, 0x75, 0x66, 0x12, 0x1c, 0x0a,
	0x09, 0x6b, 0x65, 0x65, 0x70, 0x61, 0x6c, 0x69, 0x76, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x09, 0x6b, 0x65, 0x65, 0x70, 0x61, 0x6c, 0x69, 0x76, 0x65, 0x1a, 0x3f, 0x0a, 0x0f, 0x6c,
	0x6f, 0x67, 0x69, 0x63, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x47, 0x72, 0x70, 0x63, 0x12, 0x12,
	0x0a, 0x04, 0x61, 0x64, 0x64, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x64,
	0x64, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x22, 0x43, 0x0a, 0x0b,
	0x4c, 0x6f, 0x67, 0x69, 0x63, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x1a, 0x34, 0x0a, 0x0a, 0x67,
	0x72, 0x70, 0x63, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x6f, 0x73,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x6f, 0x72,
	0x74, 0x2a, 0x21, 0x0a, 0x09, 0x51, 0x75, 0x65, 0x75, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x09,
	0x0a, 0x05, 0x52, 0x45, 0x44, 0x49, 0x53, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x4b, 0x41, 0x46,
	0x4b, 0x41, 0x10, 0x01, 0x42, 0x27, 0x5a, 0x25, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x32, 0x70, 0x67, 0x63, 0x6e, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x69, 0x6d, 0x2f,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x3b, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_config_config_proto_rawDescOnce sync.Once
	file_config_config_proto_rawDescData = file_config_config_proto_rawDesc
)

func file_config_config_proto_rawDescGZIP() []byte {
	file_config_config_proto_rawDescOnce.Do(func() {
		file_config_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_config_config_proto_rawDescData)
	})
	return file_config_config_proto_rawDescData
}

var file_config_config_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_config_config_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_config_config_proto_goTypes = []interface{}{
	(QueueType)(0),                     // 0: config.QueueType
	(*CometConfig)(nil),                // 1: config.CometConfig
	(*LogicConfig)(nil),                // 2: config.LogicConfig
	(*CometConfigServerMsg)(nil),       // 3: config.CometConfig.serverMsg
	(*CometConfigQueueMsg)(nil),        // 4: config.CometConfig.queueMsg
	(*CometConfigSslMsg)(nil),          // 5: config.CometConfig.sslMsg
	(*CometConfigTcpMsg)(nil),          // 6: config.CometConfig.tcpMsg
	(*CometConfigLogicClientGrpc)(nil), // 7: config.CometConfig.logicClientGrpc
	(*CometConfigQueueMsg_Redis)(nil),  // 8: config.CometConfig.queueMsg.Redis
	(*LogicConfigGrpcServer)(nil),      // 9: config.LogicConfig.grpcServer
}
var file_config_config_proto_depIdxs = []int32{
	3, // 0: config.CometConfig.server:type_name -> config.CometConfig.serverMsg
	4, // 1: config.CometConfig.queue:type_name -> config.CometConfig.queueMsg
	5, // 2: config.CometConfig.ssl:type_name -> config.CometConfig.sslMsg
	6, // 3: config.CometConfig.tcp:type_name -> config.CometConfig.tcpMsg
	7, // 4: config.CometConfig.logic_client_grpc:type_name -> config.CometConfig.logicClientGrpc
	8, // 5: config.CometConfig.queueMsg.redis:type_name -> config.CometConfig.queueMsg.Redis
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_config_config_proto_init() }
func file_config_config_proto_init() {
	if File_config_config_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_config_config_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CometConfig); i {
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
		file_config_config_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogicConfig); i {
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
		file_config_config_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CometConfigServerMsg); i {
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
		file_config_config_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CometConfigQueueMsg); i {
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
		file_config_config_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CometConfigSslMsg); i {
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
		file_config_config_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CometConfigTcpMsg); i {
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
		file_config_config_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CometConfigLogicClientGrpc); i {
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
		file_config_config_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CometConfigQueueMsg_Redis); i {
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
		file_config_config_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogicConfigGrpcServer); i {
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
			RawDescriptor: file_config_config_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_config_config_proto_goTypes,
		DependencyIndexes: file_config_config_proto_depIdxs,
		EnumInfos:         file_config_config_proto_enumTypes,
		MessageInfos:      file_config_config_proto_msgTypes,
	}.Build()
	File_config_config_proto = out.File
	file_config_config_proto_rawDesc = nil
	file_config_config_proto_goTypes = nil
	file_config_config_proto_depIdxs = nil
}
