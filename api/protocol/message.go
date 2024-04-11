package protocol

import (
	"fmt"
	"github.com/golang/protobuf/proto"
)

// Proto client->comet
type Proto struct {
	Version  uint16 //2 byte 16 bit
	Op       uint16
	Checksum uint16
	Seq      uint16
	Data     []byte
}
type Auth struct {
	Appid string `json:"appid"`
	Token string `json:"token"`
}

func (p *Proto) SetErrReply(err error) {
	p.Op = OpErrReply
	p.Data = []byte(err.Error())
}

func (p *Proto) String() string {
	return fmt.Sprintf("version:%d,op:%d,checksum:%d,seq:%d,data:%+v", p.Version, p.Op, p.Checksum, p.Seq, string(p.Data))
}

//data=comet->logic->queue->comet

// comet->logic
func (md *Msg) Marshal() ([]byte, error) {
	return proto.Marshal(md)
}

func (md *Msg) Unmarshal(b []byte) error {
	return proto.Unmarshal(b, md)
}

func (t Type) ToOp() uint16 {
	switch t {
	case Type_APP:
		return OpSendAreaMsgReply
	case Type_ROOM:
		return OpSendMsgRoomReply
	case Type_PUSH:
		return OpSendMsgReply
	case Type_ACK:
		return OpAck
	default:
		return OpHeartbeatReply
	}
}
