package protocol

import (
	"encoding/json"
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

func (p *Proto) String() string {
	var msg *Msg
	_ = json.Unmarshal(p.Data, msg)
	return fmt.Sprintf("version:%d,op:%d,checksum:%d,seq:%d,data:%+v", p.Version, p.Op, p.Checksum, p.Seq, msg)
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
	default:
		return OpHeartbeatReply
	}
}

//func MarshalErr(e *error2.Error) (res []byte) {
//	//res, _ = proto.Marshal(&Msg{
//	//	Code:    e.Code,
//	//	Message: e.Message,
//	//})
//	return res
//}
