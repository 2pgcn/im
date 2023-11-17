package comet

import (
	"github.com/2pgcn/gameim/api/protocol"
	"github.com/golang/protobuf/proto"
)

func (md *MsgData) Marshal() ([]byte, error) {
	return proto.Marshal(md)
}

func (md *MsgData) Unmarshal(b []byte) error {
	return proto.Unmarshal(b, md)
}

func (t Type) ToOp() uint16 {
	switch t {
	case Type_APP:
		return protocol.OpSendAreaMsgReply
	case Type_ROOM:
		return protocol.OpSendMsgRoomReply
	case Type_PUSH:
		return protocol.OpSendMsgReply
	default:
		return protocol.OpHeartbeatReply
	}
}
