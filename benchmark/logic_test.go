package main

import (
	"context"
	"github.com/2pgcn/gameim/api/comet"
	"github.com/2pgcn/gameim/api/logic"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	grpc2 "google.golang.org/grpc"
	"testing"
	"time"
)

func getLogicGrpcCLient(ctx context.Context, endpoint string) (*grpc2.ClientConn, error) {
	conn, err := grpc.DialInsecure(ctx,
		grpc.WithEndpoint(endpoint),
		grpc.WithTimeout(time.Second*10),
	)
	return conn, err
}

func TestKafkaClinet(t *testing.T) {
	cc, err := getLogicGrpcCLient(context.Background(), "127.0.0.1:9001")
	cli := logic.NewLogicClient(cc)
	if err != nil {
		t.Errorf("get logic grpc client error:%s", err.Error())
	}
	ctx, canal := context.WithTimeout(context.Background(), time.Second*10)
	defer canal()
	_, err = cli.OnMessage(ctx, &logic.MessageReq{
		Type:     comet.Type_PUSH,
		CometKey: "test_key",
		ToId:     0,
		SendId:   0,
		Msg:      []byte("hello world"),
	})
	if err != nil {
		t.Logf("logic send message error %s", err.Error())
	}
}

func BenchmarkTestLogicMessage(b *testing.B) {
	cc, err := getLogicGrpcCLient(context.Background(), "127.0.0.1:9001")
	cli := logic.NewLogicClient(cc)
	if err != nil {
		b.Errorf("get logic grpc client error:%s", err.Error())
	}
	ctx, canal := context.WithTimeout(context.Background(), time.Second*10)
	defer canal()
	b.Log(b.N)
	for i := 0; i < b.N; i++ {
		_, err = cli.OnMessage(ctx, &logic.MessageReq{
			Type:     comet.Type_PUSH,
			CometKey: "test_key",
			ToId:     0,
			SendId:   0,
			Msg:      []byte("hello world"),
		})
		if err != nil {
			b.Logf("logic send message error %s", err.Error())
		}
	}
}
