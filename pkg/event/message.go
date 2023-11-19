package event

import (
	"fmt"
	"github.com/2pgcn/gameim/api/comet"
	"github.com/segmentio/kafka-go"
	"strconv"
)

//var (
//	_ event.Sender   = (*kafkaSender)(nil)
//	_ event.Receiver = (*kafkaReceiver)(nil)
//	_ event.Event    = (*Message)(nil)
//)

type Msg struct {
	traceName string
	H         EventHeader
	K         []byte
	Data      *comet.MsgData
}

func (m *Msg) StartTrace(traceName string) {
	m.traceName = traceName
}

func (m *Msg) GetData() *comet.MsgData {
	return m.Data
}

func (m *Msg) Header() EventHeader {
	return m.H
}

func (m *Msg) GetKafkaCommitMsg() (kmsg kafka.Message, err error) {
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
func (m *Msg) Key() []byte {
	return m.K
}
func (m *Msg) Value() (res []byte) {
	var err error
	if res, err = m.GetData().Marshal(); err != nil {
		return nil
	}
	return res
}

func (m *Msg) RawValue() any {
	return m.Data
}

func (m *Msg) String() string {
	return fmt.Sprintf("head:%+v:key:%s,data:%+v", m.Header(), string(m.Key()), m.RawValue())
}
