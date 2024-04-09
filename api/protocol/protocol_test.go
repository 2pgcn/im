package protocol

import (
	"bufio"
	"github.com/2pgcn/gameim/pkg/safe"
	"google.golang.org/protobuf/proto"
	"net"
	"testing"
)

func BenchmarkTest(b *testing.B) {
	pool := safe.NewPool(func() any {
		return make([]byte, HeaderLen+16)
	})
	pool.Grow(10240)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := &Proto{
			Version:  1,
			Op:       OpHeartbeatReply,
			Checksum: 1,
			Seq:      1,
			Data:     []byte("hello world"),
		}
		data := pool.Get().([]byte)
		err := p.SerializeTo(data)
		if err != nil {
			b.Fatal(err)
		}
		//pool.Put(data)
	}
}

func BenchmarkSerializeTo(b *testing.B) {
	server, client := net.Pipe()
	cliBuf := bufio.NewWriter(client)
	serBuf := bufio.NewReader(server)
	hcProto := ProtoPool.Get()
	hcProto.Op = OpSendAreaMsg
	hcProto.Seq = 1
	hcProto.Data, _ = proto.Marshal(&Msg{
		Type:   Type_ROOM,
		ToId:   "1",
		SendId: "1",
		Msg:    []byte("hello world gameim"),
	})

	step := make(chan bool)
	hsProto := ProtoPool.Get()
	go func() {
		for i := 0; i < b.N; i++ {
			err := hcProto.writeTcp(cliBuf)
			if err != nil {
				b.Error(err)
			}
		}
	}()
	go func() {
		for i := 0; i < b.N; i++ {
			err := hsProto.DecodeFromBytes(serBuf)
			if err != nil {
				b.Error(err)
			}
			if hcProto.Op != hsProto.Op || hcProto.Seq != hsProto.Seq || len(hcProto.Data) != len(hsProto.Data) {
				b.Errorf("client write(%v):server read error,data(%v)", hcProto, hsProto)
			}
		}
		step <- true
	}()
	<-step
}
