package event

import (
	"encoding/json"
	"fmt"
	"github.com/2pgcn/gameim/api/gerr"
	"github.com/2pgcn/gameim/api/protocol"
	"github.com/2pgcn/gameim/pkg/gamelog"
	"github.com/2pgcn/gameim/pkg/safe"
	"github.com/golang/protobuf/proto"
	"github.com/segmentio/kafka-go"
	"strconv"
	"sync"
)

var queueMsgPool *safe.Pool[*queueMsg]
var msgPool *safe.Pool[*protocol.Msg]
var defaultPoolSize = 4096
var once sync.Once

func init() {
	once.Do(func() {
		msgPool = safe.NewPool(func() *protocol.Msg {
			return &protocol.Msg{}
		})
		msgPool.Grow(defaultPoolSize)
		queueMsgPool = safe.NewPool(func() *queueMsg {
			return &queueMsg{
				H:    make(map[string]any, 8),
				Data: msgPool.Get(),
			}
		})
		queueMsgPool.Grow(defaultPoolSize)
	})
}

type queueMsg struct {
	traceName string
	H         EventHeader
	Data      *protocol.Msg
}

func GetQueueMsg() *queueMsg {
	return queueMsgPool.Get()
}
func PutQueueMsg(m *queueMsg) {
	if m.Data != nil {
		msgPool.Put(m.Data)
		m.Data = nil
	}
	queueMsgPool.Put(m)
}

func GetMsg() *protocol.Msg {
	return msgPool.Get()
}
func PutMsg(m *protocol.Msg) {
	msgPool.Put(m)
}

func (m *queueMsg) StartTrace(traceName string) {
	m.traceName = traceName
}

func (m *queueMsg) Header() EventHeader {
	return m.H
}
func (m *queueMsg) GetQueueMsg() *queueMsg {
	return m
}
func (m *queueMsg) ToProtocol() (p *protocol.Proto, err error) {
	p = protocol.ProtoPool.Get()
	p.Version = protocol.Version
	p.Op = m.Data.Type.ToOp()
	reply, err := proto.Marshal(&protocol.Reply{
		Code: 0,
		Msg:  m.Data,
	})
	if err != nil {
		return p, gerr.ErrorServerError("protocol error,queuedata:%+v", m).WithMetadata(gerr.GetStack()).WithCause(err)
	}
	p.Data = reply
	return
}

func (m *queueMsg) GetKafkaCommitMsg() (kmsg kafka.Message, err error) {
	var (
		partitionStrInt int64
		offset          int64
	)
	partitionStrInt, err = strconv.ParseInt(m.Header()["partition"].(string), 10, 64)
	if err != nil {
		return
	}
	offset, err = strconv.ParseInt(m.Header()["offset"].(string), 10, 64)
	kmsg.Partition = int(partitionStrInt)
	kmsg.Offset = offset
	return kmsg, err
}

func (m *queueMsg) Value() (res []byte) {
	var err error
	if res, err = json.Marshal(m); err != nil {
		gamelog.Errorf("queue msg error:%+v", m)
		return []byte{}
	}
	return res
}

func (m *queueMsg) String() string {
	return fmt.Sprintf("head:%+v,data:%+v", m.Header(), m.Data)
}
