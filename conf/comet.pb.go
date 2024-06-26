// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.19.1
// source: comet.proto

package conf

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ServerMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name           string    `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Addrs          []string  `protobuf:"bytes,2,rep,name=addrs,proto3" json:"addrs,omitempty"`
	Port           int32     `protobuf:"varint,3,opt,name=port,proto3" json:"port,omitempty"`
	MaxConnections uint32    `protobuf:"varint,4,opt,name=max_connections,json=maxConnections,proto3" json:"max_connections,omitempty"` //max conn nums
	Queue          *QueueMsg `protobuf:"bytes,5,opt,name=queue,proto3" json:"queue,omitempty"`
}

func (x *ServerMsg) Reset() {
	*x = ServerMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerMsg) ProtoMessage() {}

func (x *ServerMsg) ProtoReflect() protoreflect.Message {
	mi := &file_comet_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerMsg.ProtoReflect.Descriptor instead.
func (*ServerMsg) Descriptor() ([]byte, []int) {
	return file_comet_proto_rawDescGZIP(), []int{0}
}

func (x *ServerMsg) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ServerMsg) GetAddrs() []string {
	if x != nil {
		return x.Addrs
	}
	return nil
}

func (x *ServerMsg) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

func (x *ServerMsg) GetMaxConnections() uint32 {
	if x != nil {
		return x.MaxConnections
	}
	return 0
}

func (x *ServerMsg) GetQueue() *QueueMsg {
	if x != nil {
		return x.Queue
	}
	return nil
}

type SslMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Enable bool   `protobuf:"varint,1,opt,name=enable,proto3" json:"enable,omitempty"`
	Key    string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Cert   string `protobuf:"bytes,3,opt,name=cert,proto3" json:"cert,omitempty"`
}

func (x *SslMsg) Reset() {
	*x = SslMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SslMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SslMsg) ProtoMessage() {}

func (x *SslMsg) ProtoReflect() protoreflect.Message {
	mi := &file_comet_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SslMsg.ProtoReflect.Descriptor instead.
func (*SslMsg) Descriptor() ([]byte, []int) {
	return file_comet_proto_rawDescGZIP(), []int{1}
}

func (x *SslMsg) GetEnable() bool {
	if x != nil {
		return x.Enable
	}
	return false
}

func (x *SslMsg) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *SslMsg) GetCert() string {
	if x != nil {
		return x.Cert
	}
	return ""
}

type TcpMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SendBuf uint64 `protobuf:"varint,1,opt,name=send_buf,json=sendBuf,proto3" json:"send_buf,omitempty"`
	RecvBuf uint64 `protobuf:"varint,2,opt,name=recv_buf,json=recvBuf,proto3" json:"recv_buf,omitempty"`
}

func (x *TcpMsg) Reset() {
	*x = TcpMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TcpMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TcpMsg) ProtoMessage() {}

func (x *TcpMsg) ProtoReflect() protoreflect.Message {
	mi := &file_comet_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TcpMsg.ProtoReflect.Descriptor instead.
func (*TcpMsg) Descriptor() ([]byte, []int) {
	return file_comet_proto_rawDescGZIP(), []int{2}
}

func (x *TcpMsg) GetSendBuf() uint64 {
	if x != nil {
		return x.SendBuf
	}
	return 0
}

func (x *TcpMsg) GetRecvBuf() uint64 {
	if x != nil {
		return x.RecvBuf
	}
	return 0
}

type LogicClientGrpc struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Addr      string               `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	Timeout   *durationpb.Duration `protobuf:"bytes,2,opt,name=timeout,proto3" json:"timeout,omitempty"`
	ClientNum int32                `protobuf:"varint,3,opt,name=client_num,json=clientNum,proto3" json:"client_num,omitempty"`
}

func (x *LogicClientGrpc) Reset() {
	*x = LogicClientGrpc{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogicClientGrpc) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogicClientGrpc) ProtoMessage() {}

func (x *LogicClientGrpc) ProtoReflect() protoreflect.Message {
	mi := &file_comet_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogicClientGrpc.ProtoReflect.Descriptor instead.
func (*LogicClientGrpc) Descriptor() ([]byte, []int) {
	return file_comet_proto_rawDescGZIP(), []int{3}
}

func (x *LogicClientGrpc) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *LogicClientGrpc) GetTimeout() *durationpb.Duration {
	if x != nil {
		return x.Timeout
	}
	return nil
}

func (x *LogicClientGrpc) GetClientNum() int32 {
	if x != nil {
		return x.ClientNum
	}
	return 0
}

type AppConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Appid string `protobuf:"bytes,1,opt,name=appid,proto3" json:"appid,omitempty"`
	// todo 暂未启用认证
	Pwd string `protobuf:"bytes,2,opt,name=pwd,proto3" json:"pwd,omitempty"`
	// bucket 根据bucketNum哈希 减少app里存储的全区的锁力度
	BucketNum        int32                `protobuf:"varint,3,opt,name=bucketNum,proto3" json:"bucketNum,omitempty"`
	KeepaliveTimeout *durationpb.Duration `protobuf:"bytes,4,opt,name=keepalive_timeout,json=keepaliveTimeout,proto3" json:"keepalive_timeout,omitempty"`
	LogicClientGrpc  *LogicClientGrpc     `protobuf:"bytes,5,opt,name=logic_client_grpc,json=logicClientGrpc,proto3" json:"logic_client_grpc,omitempty"`
}

func (x *AppConfig) Reset() {
	*x = AppConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppConfig) ProtoMessage() {}

func (x *AppConfig) ProtoReflect() protoreflect.Message {
	mi := &file_comet_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppConfig.ProtoReflect.Descriptor instead.
func (*AppConfig) Descriptor() ([]byte, []int) {
	return file_comet_proto_rawDescGZIP(), []int{4}
}

func (x *AppConfig) GetAppid() string {
	if x != nil {
		return x.Appid
	}
	return ""
}

func (x *AppConfig) GetPwd() string {
	if x != nil {
		return x.Pwd
	}
	return ""
}

func (x *AppConfig) GetBucketNum() int32 {
	if x != nil {
		return x.BucketNum
	}
	return 0
}

func (x *AppConfig) GetKeepaliveTimeout() *durationpb.Duration {
	if x != nil {
		return x.KeepaliveTimeout
	}
	return nil
}

func (x *AppConfig) GetLogicClientGrpc() *LogicClientGrpc {
	if x != nil {
		return x.LogicClientGrpc
	}
	return nil
}

type CometConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Server    *ServerMsg   `protobuf:"bytes,1,opt,name=server,proto3" json:"server,omitempty"`
	Queue     *QueueMsg    `protobuf:"bytes,2,opt,name=queue,proto3" json:"queue,omitempty"`
	Ssl       *SslMsg      `protobuf:"bytes,3,opt,name=ssl,proto3" json:"ssl,omitempty"`
	Tcp       *TcpMsg      `protobuf:"bytes,4,opt,name=tcp,proto3" json:"tcp,omitempty"`
	AppConfig []*AppConfig `protobuf:"bytes,6,rep,name=app_config,json=appConfig,proto3" json:"app_config,omitempty"`
	UpData    *UpData      `protobuf:"bytes,7,opt,name=up_data,json=upData,proto3" json:"up_data,omitempty"`
}

func (x *CometConfig) Reset() {
	*x = CometConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CometConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CometConfig) ProtoMessage() {}

func (x *CometConfig) ProtoReflect() protoreflect.Message {
	mi := &file_comet_proto_msgTypes[5]
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
	return file_comet_proto_rawDescGZIP(), []int{5}
}

func (x *CometConfig) GetServer() *ServerMsg {
	if x != nil {
		return x.Server
	}
	return nil
}

func (x *CometConfig) GetQueue() *QueueMsg {
	if x != nil {
		return x.Queue
	}
	return nil
}

func (x *CometConfig) GetSsl() *SslMsg {
	if x != nil {
		return x.Ssl
	}
	return nil
}

func (x *CometConfig) GetTcp() *TcpMsg {
	if x != nil {
		return x.Tcp
	}
	return nil
}

func (x *CometConfig) GetAppConfig() []*AppConfig {
	if x != nil {
		return x.AppConfig
	}
	return nil
}

func (x *CometConfig) GetUpData() *UpData {
	if x != nil {
		return x.UpData
	}
	return nil
}

type UpData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pyroscope *Pyroscope `protobuf:"bytes,1,opt,name=pyroscope,proto3" json:"pyroscope,omitempty"`
	TraceConf *TraceConf `protobuf:"bytes,2,opt,name=trace_conf,json=traceConf,proto3" json:"trace_conf,omitempty"`
}

func (x *UpData) Reset() {
	*x = UpData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpData) ProtoMessage() {}

func (x *UpData) ProtoReflect() protoreflect.Message {
	mi := &file_comet_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpData.ProtoReflect.Descriptor instead.
func (*UpData) Descriptor() ([]byte, []int) {
	return file_comet_proto_rawDescGZIP(), []int{6}
}

func (x *UpData) GetPyroscope() *Pyroscope {
	if x != nil {
		return x.Pyroscope
	}
	return nil
}

func (x *UpData) GetTraceConf() *TraceConf {
	if x != nil {
		return x.TraceConf
	}
	return nil
}

var File_comet_proto protoreflect.FileDescriptor

var file_comet_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x63,
	0x6f, 0x6e, 0x66, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x98, 0x01, 0x0a, 0x09, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4d, 0x73, 0x67, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x64, 0x64, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x05, 0x61, 0x64, 0x64, 0x72, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x27, 0x0a,
	0x0f, 0x6d, 0x61, 0x78, 0x5f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0e, 0x6d, 0x61, 0x78, 0x43, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x24, 0x0a, 0x05, 0x71, 0x75, 0x65, 0x75, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x2e, 0x71, 0x75, 0x65,
	0x75, 0x65, 0x4d, 0x73, 0x67, 0x52, 0x05, 0x71, 0x75, 0x65, 0x75, 0x65, 0x22, 0x46, 0x0a, 0x06,
	0x73, 0x73, 0x6c, 0x4d, 0x73, 0x67, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x12, 0x0a, 0x04, 0x63, 0x65, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x63, 0x65, 0x72, 0x74, 0x22, 0x3e, 0x0a, 0x06, 0x74, 0x63, 0x70, 0x4d, 0x73, 0x67, 0x12, 0x19,
	0x0a, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x5f, 0x62, 0x75, 0x66, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x07, 0x73, 0x65, 0x6e, 0x64, 0x42, 0x75, 0x66, 0x12, 0x19, 0x0a, 0x08, 0x72, 0x65, 0x63,
	0x76, 0x5f, 0x62, 0x75, 0x66, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x72, 0x65, 0x63,
	0x76, 0x42, 0x75, 0x66, 0x22, 0x79, 0x0a, 0x0f, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x43, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x47, 0x72, 0x70, 0x63, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x64, 0x64, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72, 0x12, 0x33, 0x0a, 0x07, 0x74,
	0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74,
	0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x6e, 0x75, 0x6d, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4e, 0x75, 0x6d, 0x22,
	0xdd, 0x01, 0x0a, 0x0a, 0x61, 0x70, 0x70, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x14,
	0x0a, 0x05, 0x61, 0x70, 0x70, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x61,
	0x70, 0x70, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x77, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x70, 0x77, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74,
	0x4e, 0x75, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x62, 0x75, 0x63, 0x6b, 0x65,
	0x74, 0x4e, 0x75, 0x6d, 0x12, 0x46, 0x0a, 0x11, 0x6b, 0x65, 0x65, 0x70, 0x61, 0x6c, 0x69, 0x76,
	0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x10, 0x6b, 0x65, 0x65, 0x70,
	0x61, 0x6c, 0x69, 0x76, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12, 0x41, 0x0a, 0x11,
	0x6c, 0x6f, 0x67, 0x69, 0x63, 0x5f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x67, 0x72, 0x70,
	0x63, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x2e, 0x6c,
	0x6f, 0x67, 0x69, 0x63, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x47, 0x72, 0x70, 0x63, 0x52, 0x0f,
	0x6c, 0x6f, 0x67, 0x69, 0x63, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x47, 0x72, 0x70, 0x63, 0x22,
	0xf4, 0x01, 0x0a, 0x0b, 0x43, 0x6f, 0x6d, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12,
	0x27, 0x0a, 0x06, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0f, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4d, 0x73, 0x67,
	0x52, 0x06, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x24, 0x0a, 0x05, 0x71, 0x75, 0x65, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x2e, 0x71,
	0x75, 0x65, 0x75, 0x65, 0x4d, 0x73, 0x67, 0x52, 0x05, 0x71, 0x75, 0x65, 0x75, 0x65, 0x12, 0x1e,
	0x0a, 0x03, 0x73, 0x73, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x63, 0x6f,
	0x6e, 0x66, 0x2e, 0x73, 0x73, 0x6c, 0x4d, 0x73, 0x67, 0x52, 0x03, 0x73, 0x73, 0x6c, 0x12, 0x1e,
	0x0a, 0x03, 0x74, 0x63, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x63, 0x6f,
	0x6e, 0x66, 0x2e, 0x74, 0x63, 0x70, 0x4d, 0x73, 0x67, 0x52, 0x03, 0x74, 0x63, 0x70, 0x12, 0x2f,
	0x0a, 0x0a, 0x61, 0x70, 0x70, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x06, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x10, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x2e, 0x61, 0x70, 0x70, 0x5f, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x52, 0x09, 0x61, 0x70, 0x70, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12,
	0x25, 0x0a, 0x07, 0x75, 0x70, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0c, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x2e, 0x55, 0x70, 0x44, 0x61, 0x74, 0x61, 0x52, 0x06,
	0x75, 0x70, 0x44, 0x61, 0x74, 0x61, 0x22, 0x67, 0x0a, 0x06, 0x55, 0x70, 0x44, 0x61, 0x74, 0x61,
	0x12, 0x2d, 0x0a, 0x09, 0x70, 0x79, 0x72, 0x6f, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x2e, 0x50, 0x79, 0x72, 0x6f, 0x73,
	0x63, 0x6f, 0x70, 0x65, 0x52, 0x09, 0x70, 0x79, 0x72, 0x6f, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x12,
	0x2e, 0x0a, 0x0a, 0x74, 0x72, 0x61, 0x63, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x2e, 0x54, 0x72, 0x61, 0x63, 0x65,
	0x43, 0x6f, 0x6e, 0x66, 0x52, 0x09, 0x74, 0x72, 0x61, 0x63, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x42,
	0x23, 0x5a, 0x21, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x32, 0x70,
	0x67, 0x63, 0x6e, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x69, 0x6d, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x3b,
	0x63, 0x6f, 0x6e, 0x66, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_comet_proto_rawDescOnce sync.Once
	file_comet_proto_rawDescData = file_comet_proto_rawDesc
)

func file_comet_proto_rawDescGZIP() []byte {
	file_comet_proto_rawDescOnce.Do(func() {
		file_comet_proto_rawDescData = protoimpl.X.CompressGZIP(file_comet_proto_rawDescData)
	})
	return file_comet_proto_rawDescData
}

var file_comet_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_comet_proto_goTypes = []interface{}{
	(*ServerMsg)(nil),           // 0: conf.serverMsg
	(*SslMsg)(nil),              // 1: conf.sslMsg
	(*TcpMsg)(nil),              // 2: conf.tcpMsg
	(*LogicClientGrpc)(nil),     // 3: conf.logicClientGrpc
	(*AppConfig)(nil),           // 4: conf.app_config
	(*CometConfig)(nil),         // 5: conf.CometConfig
	(*UpData)(nil),              // 6: conf.UpData
	(*QueueMsg)(nil),            // 7: conf.queueMsg
	(*durationpb.Duration)(nil), // 8: google.protobuf.Duration
	(*Pyroscope)(nil),           // 9: conf.Pyroscope
	(*TraceConf)(nil),           // 10: conf.TraceConf
}
var file_comet_proto_depIdxs = []int32{
	7,  // 0: conf.serverMsg.queue:type_name -> conf.queueMsg
	8,  // 1: conf.logicClientGrpc.timeout:type_name -> google.protobuf.Duration
	8,  // 2: conf.app_config.keepalive_timeout:type_name -> google.protobuf.Duration
	3,  // 3: conf.app_config.logic_client_grpc:type_name -> conf.logicClientGrpc
	0,  // 4: conf.CometConfig.server:type_name -> conf.serverMsg
	7,  // 5: conf.CometConfig.queue:type_name -> conf.queueMsg
	1,  // 6: conf.CometConfig.ssl:type_name -> conf.sslMsg
	2,  // 7: conf.CometConfig.tcp:type_name -> conf.tcpMsg
	4,  // 8: conf.CometConfig.app_config:type_name -> conf.app_config
	6,  // 9: conf.CometConfig.up_data:type_name -> conf.UpData
	9,  // 10: conf.UpData.pyroscope:type_name -> conf.Pyroscope
	10, // 11: conf.UpData.trace_conf:type_name -> conf.TraceConf
	12, // [12:12] is the sub-list for method output_type
	12, // [12:12] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_comet_proto_init() }
func file_comet_proto_init() {
	if File_comet_proto != nil {
		return
	}
	file_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_comet_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerMsg); i {
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
		file_comet_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SslMsg); i {
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
		file_comet_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TcpMsg); i {
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
		file_comet_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogicClientGrpc); i {
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
		file_comet_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppConfig); i {
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
		file_comet_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
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
		file_comet_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpData); i {
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
			RawDescriptor: file_comet_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_comet_proto_goTypes,
		DependencyIndexes: file_comet_proto_depIdxs,
		MessageInfos:      file_comet_proto_msgTypes,
	}.Build()
	File_comet_proto = out.File
	file_comet_proto_rawDesc = nil
	file_comet_proto_goTypes = nil
	file_comet_proto_depIdxs = nil
}
