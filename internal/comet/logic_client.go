package comet

import (
	"context"
	"github.com/2pgcn/gameim/api/logic"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"strconv"
	"sync/atomic"
)

var online uint64 = 0

type logicClientTest struct {
	cc grpc.ClientConnInterface
}

func (c *logicClientTest) OnAuth(ctx context.Context, in *logic.AuthReq, opts ...grpc.CallOption) (*logic.AuthReply, error) {
	atomic.AddUint64(&online, 1)
	uid := atomic.LoadUint64(&online)
	return &logic.AuthReply{
		Uid: strconv.Itoa(int(uid)),
		//AreaId: uint64(1),
		RoomId: strconv.FormatUint(uid%10000, 10),
	}, nil
}

func (c *logicClientTest) OnConnect(ctx context.Context, in *logic.ConnectReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, nil
}

func (c *logicClientTest) OnMessage(ctx context.Context, in *logic.MessageReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	//发送chan
	return nil, nil
}

func (c *logicClientTest) OnClose(ctx context.Context, in *logic.CloseReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, nil
}
