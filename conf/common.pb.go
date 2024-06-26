// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.19.1
// source: common.proto

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

type QueueType int32

const (
	QueueType_KAFKA QueueType = 0
	QueueType_NSQ   QueueType = 1
	QueueType_REDIS QueueType = 2
	// for local test
	QueueType_SOCK QueueType = 3
)

// Enum value maps for QueueType.
var (
	QueueType_name = map[int32]string{
		0: "KAFKA",
		1: "NSQ",
		2: "REDIS",
		3: "SOCK",
	}
	QueueType_value = map[string]int32{
		"KAFKA": 0,
		"NSQ":   1,
		"REDIS": 2,
		"SOCK":  3,
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
	return file_common_proto_enumTypes[0].Descriptor()
}

func (QueueType) Type() protoreflect.EnumType {
	return &file_common_proto_enumTypes[0]
}

func (x QueueType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use QueueType.Descriptor instead.
func (QueueType) EnumDescriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{0}
}

type QueueMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type  QueueType `protobuf:"varint,1,opt,name=type,proto3,enum=conf.QueueType" json:"type,omitempty"`
	Queue *Queue    `protobuf:"bytes,2,opt,name=queue,proto3" json:"queue,omitempty"`
	Redis *Redis    `protobuf:"bytes,3,opt,name=redis,proto3" json:"redis,omitempty"`
	Nsq   *Nsq      `protobuf:"bytes,4,opt,name=nsq,proto3" json:"nsq,omitempty"`
	Kafka *Kafka    `protobuf:"bytes,5,opt,name=kafka,proto3" json:"kafka,omitempty"`
	Sock  *Sock     `protobuf:"bytes,6,opt,name=sock,proto3" json:"sock,omitempty"`
}

func (x *QueueMsg) Reset() {
	*x = QueueMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueueMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueMsg) ProtoMessage() {}

func (x *QueueMsg) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueMsg.ProtoReflect.Descriptor instead.
func (*QueueMsg) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{0}
}

func (x *QueueMsg) GetType() QueueType {
	if x != nil {
		return x.Type
	}
	return QueueType_KAFKA
}

func (x *QueueMsg) GetQueue() *Queue {
	if x != nil {
		return x.Queue
	}
	return nil
}

func (x *QueueMsg) GetRedis() *Redis {
	if x != nil {
		return x.Redis
	}
	return nil
}

func (x *QueueMsg) GetNsq() *Nsq {
	if x != nil {
		return x.Nsq
	}
	return nil
}

func (x *QueueMsg) GetKafka() *Kafka {
	if x != nil {
		return x.Kafka
	}
	return nil
}

func (x *QueueMsg) GetSock() *Sock {
	if x != nil {
		return x.Sock
	}
	return nil
}

type Queue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Topic        string `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
	Channel      string `protobuf:"bytes,2,opt,name=channel,proto3" json:"channel,omitempty"`
	RecvQueueLen int32  `protobuf:"varint,3,opt,name=recv_queue_len,json=recvQueueLen,proto3" json:"recv_queue_len,omitempty"`
	RecvQueueNum int32  `protobuf:"varint,4,opt,name=recv_queue_num,json=recvQueueNum,proto3" json:"recv_queue_num,omitempty"`
	SendQueueLen int32  `protobuf:"varint,5,opt,name=send_queue_len,json=sendQueueLen,proto3" json:"send_queue_len,omitempty"`
	SendQueueNum int32  `protobuf:"varint,6,opt,name=send_queue_num,json=sendQueueNum,proto3" json:"send_queue_num,omitempty"`
}

func (x *Queue) Reset() {
	*x = Queue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Queue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Queue) ProtoMessage() {}

func (x *Queue) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Queue.ProtoReflect.Descriptor instead.
func (*Queue) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{1}
}

func (x *Queue) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *Queue) GetChannel() string {
	if x != nil {
		return x.Channel
	}
	return ""
}

func (x *Queue) GetRecvQueueLen() int32 {
	if x != nil {
		return x.RecvQueueLen
	}
	return 0
}

func (x *Queue) GetRecvQueueNum() int32 {
	if x != nil {
		return x.RecvQueueNum
	}
	return 0
}

func (x *Queue) GetSendQueueLen() int32 {
	if x != nil {
		return x.SendQueueLen
	}
	return 0
}

func (x *Queue) GetSendQueueNum() int32 {
	if x != nil {
		return x.SendQueueNum
	}
	return 0
}

type Redis struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Addr         string               `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	Passwd       string               `protobuf:"bytes,2,opt,name=passwd,proto3" json:"passwd,omitempty"`
	ReadTimeout  *durationpb.Duration `protobuf:"bytes,3,opt,name=read_timeout,json=readTimeout,proto3" json:"read_timeout,omitempty"`
	WriteTimeout *durationpb.Duration `protobuf:"bytes,4,opt,name=write_timeout,json=writeTimeout,proto3" json:"write_timeout,omitempty"`
	KeepLive     *durationpb.Duration `protobuf:"bytes,5,opt,name=keep_live,json=keepLive,proto3" json:"keep_live,omitempty"`
	PoolSize     int32                `protobuf:"varint,6,opt,name=pool_size,json=poolSize,proto3" json:"pool_size,omitempty"`
	UseDb        int32                `protobuf:"varint,7,opt,name=use_db,json=useDb,proto3" json:"use_db,omitempty"`
}

func (x *Redis) Reset() {
	*x = Redis{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Redis) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Redis) ProtoMessage() {}

func (x *Redis) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Redis.ProtoReflect.Descriptor instead.
func (*Redis) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{2}
}

func (x *Redis) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *Redis) GetPasswd() string {
	if x != nil {
		return x.Passwd
	}
	return ""
}

func (x *Redis) GetReadTimeout() *durationpb.Duration {
	if x != nil {
		return x.ReadTimeout
	}
	return nil
}

func (x *Redis) GetWriteTimeout() *durationpb.Duration {
	if x != nil {
		return x.WriteTimeout
	}
	return nil
}

func (x *Redis) GetKeepLive() *durationpb.Duration {
	if x != nil {
		return x.KeepLive
	}
	return nil
}

func (x *Redis) GetPoolSize() int32 {
	if x != nil {
		return x.PoolSize
	}
	return 0
}

func (x *Redis) GetUseDb() int32 {
	if x != nil {
		return x.UseDb
	}
	return 0
}

type RedisQueue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Rdscfg *Redis `protobuf:"bytes,1,opt,name=rdscfg,proto3" json:"rdscfg,omitempty"`
	Queue  *Queue `protobuf:"bytes,2,opt,name=queue,proto3" json:"queue,omitempty"`
}

func (x *RedisQueue) Reset() {
	*x = RedisQueue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RedisQueue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RedisQueue) ProtoMessage() {}

func (x *RedisQueue) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RedisQueue.ProtoReflect.Descriptor instead.
func (*RedisQueue) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{3}
}

func (x *RedisQueue) GetRdscfg() *Redis {
	if x != nil {
		return x.Rdscfg
	}
	return nil
}

func (x *RedisQueue) GetQueue() *Queue {
	if x != nil {
		return x.Queue
	}
	return nil
}

type Sasl struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username   string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password   string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	InstanceId string `protobuf:"bytes,3,opt,name=instance_id,json=instanceId,proto3" json:"instance_id,omitempty"`
}

func (x *Sasl) Reset() {
	*x = Sasl{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Sasl) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Sasl) ProtoMessage() {}

func (x *Sasl) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Sasl.ProtoReflect.Descriptor instead.
func (*Sasl) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{4}
}

func (x *Sasl) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *Sasl) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *Sasl) GetInstanceId() string {
	if x != nil {
		return x.InstanceId
	}
	return ""
}

type Kafka struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Topic            string   `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
	Sasl             *Sasl    `protobuf:"bytes,2,opt,name=sasl,proto3" json:"sasl,omitempty"`
	BootstrapServers []string `protobuf:"bytes,3,rep,name=bootstrap_servers,json=bootstrapServers,proto3" json:"bootstrap_servers,omitempty"`
	ConsumerGroupId  string   `protobuf:"bytes,4,opt,name=consumer_group_id,json=consumerGroupId,proto3" json:"consumer_group_id,omitempty"`
	PartitionNum     int32    `protobuf:"varint,5,opt,name=partition_num,json=partitionNum,proto3" json:"partition_num,omitempty"`
}

func (x *Kafka) Reset() {
	*x = Kafka{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Kafka) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Kafka) ProtoMessage() {}

func (x *Kafka) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Kafka.ProtoReflect.Descriptor instead.
func (*Kafka) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{5}
}

func (x *Kafka) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *Kafka) GetSasl() *Sasl {
	if x != nil {
		return x.Sasl
	}
	return nil
}

func (x *Kafka) GetBootstrapServers() []string {
	if x != nil {
		return x.BootstrapServers
	}
	return nil
}

func (x *Kafka) GetConsumerGroupId() string {
	if x != nil {
		return x.ConsumerGroupId
	}
	return ""
}

func (x *Kafka) GetPartitionNum() int32 {
	if x != nil {
		return x.PartitionNum
	}
	return 0
}

type Nsq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Topic       string   `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
	Channel     string   `protobuf:"bytes,2,opt,name=channel,proto3" json:"channel,omitempty"`
	NsqdAddress []string `protobuf:"bytes,3,rep,name=nsqd_address,json=nsqdAddress,proto3" json:"nsqd_address,omitempty"`
	Lookupd     []string `protobuf:"bytes,4,rep,name=lookupd,proto3" json:"lookupd,omitempty"`
}

func (x *Nsq) Reset() {
	*x = Nsq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Nsq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Nsq) ProtoMessage() {}

func (x *Nsq) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Nsq.ProtoReflect.Descriptor instead.
func (*Nsq) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{6}
}

func (x *Nsq) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *Nsq) GetChannel() string {
	if x != nil {
		return x.Channel
	}
	return ""
}

func (x *Nsq) GetNsqdAddress() []string {
	if x != nil {
		return x.NsqdAddress
	}
	return nil
}

func (x *Nsq) GetLookupd() []string {
	if x != nil {
		return x.Lookupd
	}
	return nil
}

type Sock struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
}

func (x *Sock) Reset() {
	*x = Sock{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Sock) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Sock) ProtoMessage() {}

func (x *Sock) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Sock.ProtoReflect.Descriptor instead.
func (*Sock) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{7}
}

func (x *Sock) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

// 上报数据
type Pyroscope struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
}

func (x *Pyroscope) Reset() {
	*x = Pyroscope{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Pyroscope) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pyroscope) ProtoMessage() {}

func (x *Pyroscope) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pyroscope.ProtoReflect.Descriptor instead.
func (*Pyroscope) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{8}
}

func (x *Pyroscope) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

type TraceConf struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SlsProjectHeader         string `protobuf:"bytes,1,opt,name=slsProjectHeader,proto3" json:"slsProjectHeader,omitempty"`
	SlsInstanceIDHeader      string `protobuf:"bytes,2,opt,name=slsInstanceIDHeader,proto3" json:"slsInstanceIDHeader,omitempty"`
	SlsAccessKeyIDHeader     string `protobuf:"bytes,3,opt,name=slsAccessKeyIDHeader,proto3" json:"slsAccessKeyIDHeader,omitempty"`
	SlsAccessKeySecretHeader string `protobuf:"bytes,4,opt,name=slsAccessKeySecretHeader,proto3" json:"slsAccessKeySecretHeader,omitempty"`
	Endpoint                 string `protobuf:"bytes,5,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
}

func (x *TraceConf) Reset() {
	*x = TraceConf{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TraceConf) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TraceConf) ProtoMessage() {}

func (x *TraceConf) ProtoReflect() protoreflect.Message {
	mi := &file_common_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TraceConf.ProtoReflect.Descriptor instead.
func (*TraceConf) Descriptor() ([]byte, []int) {
	return file_common_proto_rawDescGZIP(), []int{9}
}

func (x *TraceConf) GetSlsProjectHeader() string {
	if x != nil {
		return x.SlsProjectHeader
	}
	return ""
}

func (x *TraceConf) GetSlsInstanceIDHeader() string {
	if x != nil {
		return x.SlsInstanceIDHeader
	}
	return ""
}

func (x *TraceConf) GetSlsAccessKeyIDHeader() string {
	if x != nil {
		return x.SlsAccessKeyIDHeader
	}
	return ""
}

func (x *TraceConf) GetSlsAccessKeySecretHeader() string {
	if x != nil {
		return x.SlsAccessKeySecretHeader
	}
	return ""
}

func (x *TraceConf) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

var File_common_proto protoreflect.FileDescriptor

var file_common_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04,
	0x63, 0x6f, 0x6e, 0x66, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd5, 0x01, 0x0a, 0x08, 0x71, 0x75, 0x65, 0x75, 0x65, 0x4d, 0x73,
	0x67, 0x12, 0x23, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x0f, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x21, 0x0a, 0x05, 0x71, 0x75, 0x65, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x2e, 0x51, 0x75, 0x65,
	0x75, 0x65, 0x52, 0x05, 0x71, 0x75, 0x65, 0x75, 0x65, 0x12, 0x21, 0x0a, 0x05, 0x72, 0x65, 0x64,
	0x69, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x2e,
	0x52, 0x65, 0x64, 0x69, 0x73, 0x52, 0x05, 0x72, 0x65, 0x64, 0x69, 0x73, 0x12, 0x1b, 0x0a, 0x03,
	0x6e, 0x73, 0x71, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x63, 0x6f, 0x6e, 0x66,
	0x2e, 0x4e, 0x73, 0x71, 0x52, 0x03, 0x6e, 0x73, 0x71, 0x12, 0x21, 0x0a, 0x05, 0x6b, 0x61, 0x66,
	0x6b, 0x61, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x2e,
	0x4b, 0x61, 0x66, 0x6b, 0x61, 0x52, 0x05, 0x6b, 0x61, 0x66, 0x6b, 0x61, 0x12, 0x1e, 0x0a, 0x04,
	0x73, 0x6f, 0x63, 0x6b, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x63, 0x6f, 0x6e,
	0x66, 0x2e, 0x53, 0x6f, 0x63, 0x6b, 0x52, 0x04, 0x73, 0x6f, 0x63, 0x6b, 0x22, 0xcf, 0x01, 0x0a,
	0x05, 0x51, 0x75, 0x65, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x18, 0x0a, 0x07,
	0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63,
	0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x24, 0x0a, 0x0e, 0x72, 0x65, 0x63, 0x76, 0x5f, 0x71,
	0x75, 0x65, 0x75, 0x65, 0x5f, 0x6c, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c,
	0x72, 0x65, 0x63, 0x76, 0x51, 0x75, 0x65, 0x75, 0x65, 0x4c, 0x65, 0x6e, 0x12, 0x24, 0x0a, 0x0e,
	0x72, 0x65, 0x63, 0x76, 0x5f, 0x71, 0x75, 0x65, 0x75, 0x65, 0x5f, 0x6e, 0x75, 0x6d, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x72, 0x65, 0x63, 0x76, 0x51, 0x75, 0x65, 0x75, 0x65, 0x4e,
	0x75, 0x6d, 0x12, 0x24, 0x0a, 0x0e, 0x73, 0x65, 0x6e, 0x64, 0x5f, 0x71, 0x75, 0x65, 0x75, 0x65,
	0x5f, 0x6c, 0x65, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x73, 0x65, 0x6e, 0x64,
	0x51, 0x75, 0x65, 0x75, 0x65, 0x4c, 0x65, 0x6e, 0x12, 0x24, 0x0a, 0x0e, 0x73, 0x65, 0x6e, 0x64,
	0x5f, 0x71, 0x75, 0x65, 0x75, 0x65, 0x5f, 0x6e, 0x75, 0x6d, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0c, 0x73, 0x65, 0x6e, 0x64, 0x51, 0x75, 0x65, 0x75, 0x65, 0x4e, 0x75, 0x6d, 0x22, 0x9d,
	0x02, 0x0a, 0x05, 0x52, 0x65, 0x64, 0x69, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x64, 0x64, 0x72,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72, 0x12, 0x16, 0x0a, 0x06,
	0x70, 0x61, 0x73, 0x73, 0x77, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x61,
	0x73, 0x73, 0x77, 0x64, 0x12, 0x3c, 0x0a, 0x0c, 0x72, 0x65, 0x61, 0x64, 0x5f, 0x74, 0x69, 0x6d,
	0x65, 0x6f, 0x75, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0b, 0x72, 0x65, 0x61, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x6f,
	0x75, 0x74, 0x12, 0x3e, 0x0a, 0x0d, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65,
	0x6f, 0x75, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0c, 0x77, 0x72, 0x69, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x6f,
	0x75, 0x74, 0x12, 0x36, 0x0a, 0x09, 0x6b, 0x65, 0x65, 0x70, 0x5f, 0x6c, 0x69, 0x76, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x08, 0x6b, 0x65, 0x65, 0x70, 0x4c, 0x69, 0x76, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x6f,
	0x6f, 0x6c, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70,
	0x6f, 0x6f, 0x6c, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x15, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x5f, 0x64,
	0x62, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x75, 0x73, 0x65, 0x44, 0x62, 0x22, 0x54,
	0x0a, 0x0a, 0x52, 0x65, 0x64, 0x69, 0x73, 0x51, 0x75, 0x65, 0x75, 0x65, 0x12, 0x23, 0x0a, 0x06,
	0x72, 0x64, 0x73, 0x63, 0x66, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x63,
	0x6f, 0x6e, 0x66, 0x2e, 0x52, 0x65, 0x64, 0x69, 0x73, 0x52, 0x06, 0x72, 0x64, 0x73, 0x63, 0x66,
	0x67, 0x12, 0x21, 0x0a, 0x05, 0x71, 0x75, 0x65, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0b, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x52, 0x05, 0x71,
	0x75, 0x65, 0x75, 0x65, 0x22, 0x5f, 0x0a, 0x04, 0x73, 0x61, 0x73, 0x6c, 0x12, 0x1a, 0x0a, 0x08,
	0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73,
	0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73,
	0x77, 0x6f, 0x72, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65,
	0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x69, 0x6e, 0x73, 0x74, 0x61,
	0x6e, 0x63, 0x65, 0x49, 0x64, 0x22, 0xbb, 0x01, 0x0a, 0x05, 0x4b, 0x61, 0x66, 0x6b, 0x61, 0x12,
	0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x74, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x1e, 0x0a, 0x04, 0x73, 0x61, 0x73, 0x6c, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x2e, 0x73, 0x61, 0x73, 0x6c, 0x52,
	0x04, 0x73, 0x61, 0x73, 0x6c, 0x12, 0x2b, 0x0a, 0x11, 0x62, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72,
	0x61, 0x70, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x10, 0x62, 0x6f, 0x6f, 0x74, 0x73, 0x74, 0x72, 0x61, 0x70, 0x53, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x73, 0x12, 0x2a, 0x0a, 0x11, 0x63, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x72, 0x5f, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x63,
	0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x72, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x64, 0x12, 0x23,
	0x0a, 0x0d, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6e, 0x75, 0x6d, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e,
	0x4e, 0x75, 0x6d, 0x22, 0x72, 0x0a, 0x03, 0x4e, 0x73, 0x71, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f,
	0x70, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63,
	0x12, 0x18, 0x0a, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x21, 0x0a, 0x0c, 0x6e, 0x73,
	0x71, 0x64, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x0b, 0x6e, 0x73, 0x71, 0x64, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a,
	0x07, 0x6c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x64, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07,
	0x6c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x64, 0x22, 0x20, 0x0a, 0x04, 0x53, 0x6f, 0x63, 0x6b, 0x12,
	0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x25, 0x0a, 0x09, 0x50, 0x79, 0x72,
	0x6f, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x22, 0xf5, 0x01, 0x0a, 0x09, 0x54, 0x72, 0x61, 0x63, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x12, 0x2a,
	0x0a, 0x10, 0x73, 0x6c, 0x73, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x48, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x73, 0x6c, 0x73, 0x50, 0x72, 0x6f,
	0x6a, 0x65, 0x63, 0x74, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x30, 0x0a, 0x13, 0x73, 0x6c,
	0x73, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x49, 0x44, 0x48, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x13, 0x73, 0x6c, 0x73, 0x49, 0x6e, 0x73, 0x74,
	0x61, 0x6e, 0x63, 0x65, 0x49, 0x44, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x32, 0x0a, 0x14,
	0x73, 0x6c, 0x73, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x4b, 0x65, 0x79, 0x49, 0x44, 0x48, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x14, 0x73, 0x6c, 0x73, 0x41,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x4b, 0x65, 0x79, 0x49, 0x44, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x12, 0x3a, 0x0a, 0x18, 0x73, 0x6c, 0x73, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x4b, 0x65, 0x79,
	0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x18, 0x73, 0x6c, 0x73, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x4b, 0x65, 0x79,
	0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08,
	0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x2a, 0x34, 0x0a, 0x09, 0x51, 0x75, 0x65, 0x75,
	0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x09, 0x0a, 0x05, 0x4b, 0x41, 0x46, 0x4b, 0x41, 0x10, 0x00,
	0x12, 0x07, 0x0a, 0x03, 0x4e, 0x53, 0x51, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x52, 0x45, 0x44,
	0x49, 0x53, 0x10, 0x02, 0x12, 0x08, 0x0a, 0x04, 0x53, 0x4f, 0x43, 0x4b, 0x10, 0x03, 0x42, 0x23,
	0x5a, 0x21, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x32, 0x70, 0x67,
	0x63, 0x6e, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x69, 0x6d, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x3b, 0x63,
	0x6f, 0x6e, 0x66, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_proto_rawDescOnce sync.Once
	file_common_proto_rawDescData = file_common_proto_rawDesc
)

func file_common_proto_rawDescGZIP() []byte {
	file_common_proto_rawDescOnce.Do(func() {
		file_common_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_proto_rawDescData)
	})
	return file_common_proto_rawDescData
}

var file_common_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_common_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_common_proto_goTypes = []interface{}{
	(QueueType)(0),              // 0: conf.QueueType
	(*QueueMsg)(nil),            // 1: conf.queueMsg
	(*Queue)(nil),               // 2: conf.Queue
	(*Redis)(nil),               // 3: conf.Redis
	(*RedisQueue)(nil),          // 4: conf.RedisQueue
	(*Sasl)(nil),                // 5: conf.sasl
	(*Kafka)(nil),               // 6: conf.Kafka
	(*Nsq)(nil),                 // 7: conf.Nsq
	(*Sock)(nil),                // 8: conf.Sock
	(*Pyroscope)(nil),           // 9: conf.Pyroscope
	(*TraceConf)(nil),           // 10: conf.TraceConf
	(*durationpb.Duration)(nil), // 11: google.protobuf.Duration
}
var file_common_proto_depIdxs = []int32{
	0,  // 0: conf.queueMsg.type:type_name -> conf.QueueType
	2,  // 1: conf.queueMsg.queue:type_name -> conf.Queue
	3,  // 2: conf.queueMsg.redis:type_name -> conf.Redis
	7,  // 3: conf.queueMsg.nsq:type_name -> conf.Nsq
	6,  // 4: conf.queueMsg.kafka:type_name -> conf.Kafka
	8,  // 5: conf.queueMsg.sock:type_name -> conf.Sock
	11, // 6: conf.Redis.read_timeout:type_name -> google.protobuf.Duration
	11, // 7: conf.Redis.write_timeout:type_name -> google.protobuf.Duration
	11, // 8: conf.Redis.keep_live:type_name -> google.protobuf.Duration
	3,  // 9: conf.RedisQueue.rdscfg:type_name -> conf.Redis
	2,  // 10: conf.RedisQueue.queue:type_name -> conf.Queue
	5,  // 11: conf.Kafka.sasl:type_name -> conf.sasl
	12, // [12:12] is the sub-list for method output_type
	12, // [12:12] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_common_proto_init() }
func file_common_proto_init() {
	if File_common_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_common_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueueMsg); i {
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
		file_common_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Queue); i {
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
		file_common_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Redis); i {
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
		file_common_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RedisQueue); i {
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
		file_common_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Sasl); i {
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
		file_common_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Kafka); i {
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
		file_common_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Nsq); i {
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
		file_common_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Sock); i {
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
		file_common_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Pyroscope); i {
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
		file_common_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TraceConf); i {
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
			RawDescriptor: file_common_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_proto_goTypes,
		DependencyIndexes: file_common_proto_depIdxs,
		EnumInfos:         file_common_proto_enumTypes,
		MessageInfos:      file_common_proto_msgTypes,
	}.Build()
	File_common_proto = out.File
	file_common_proto_rawDesc = nil
	file_common_proto_goTypes = nil
	file_common_proto_depIdxs = nil
}
