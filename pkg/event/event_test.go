package event

import (
	"github.com/2pgcn/gameim/api/comet"
	"github.com/golang/protobuf/proto"
	"testing"
)

func TestEventMsg(t *testing.T) {
	m := &Msg{
		H: map[string]any{"111": 1},
		K: []byte("1111"),
		Data: &comet.MsgData{
			Type:   1,
			ToId:   "1",
			SendId: 1,
			Msg:    []byte("test"),
		},
	}
	testV := m.Value()
	testH := m.Key()
	testD := m.Value()
	t.Log(testV, testH, testD)

	testMsgD, err := m.GetData().Marshal()
	t.Log(testMsgD, err)
	var MsgD comet.MsgData
	_ = proto.Unmarshal(testMsgD, &MsgD)
	t.Logf("%+v", MsgD)
}
