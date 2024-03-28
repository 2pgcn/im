package main

import (
	"context"
	"fmt"
	"github.com/2pgcn/gameim/api/logic"
	"github.com/2pgcn/gameim/api/protocol"
	"github.com/go-kratos/kratos/v2/middleware/circuitbreaker"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	trgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func benchLogic(ctx context.Context, address string, num int) {
	ctx, cancel := context.WithCancel(context.Background())
	client := newLogicClient(ctx, address)
	var lasterr error
	ticker := time.Tick(time.Second * 1)
	gopool.GoCtx(func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker:
				for i := 0; i < num; i++ {
					addCountSend(1)
					_, err := client.OnMessage(ctx, &logic.MessageReq{
						Type:     protocol.Type_PUSH,
						CometKey: "test",
						ToId:     "1",
						SendId:   strconv.Itoa(i),
						Msg:      []byte("hello world"),
					})
					if err != nil {
						lasterr = err
					}
					addCountDown(1)
				}
				if lasterr != nil {
					fmt.Printf("lasterr:%s", lasterr)
				}
			}
		}
	})

	gopool.GoCtx(func(ctx context.Context) {
		result(ctx)
	})
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Kill, os.Interrupt)
	go func() {
		for {
			select {
			case <-ctx.Done():
				//fmt.Println("exitChan:", sig)
				cancel()
				return
			}
		}
	}()
}

const (
	// grpc options
	grpcInitialWindowSize     = 1 << 24
	grpcInitialConnWindowSize = 1 << 24
	grpcMaxSendMsgSize        = 1 << 24
	grpcMaxCallMsgSize        = 1 << 24
	grpcKeepAliveTime         = time.Second * 10
	grpcKeepAliveTimeout      = time.Second * 3
)

func newLogicClient(ctx context.Context, addr string) logic.LogicClient {
	conn, err := trgrpc.DialInsecure(
		ctx,
		trgrpc.WithEndpoint(addr),
		trgrpc.WithTimeout(time.Second*5),
		trgrpc.WithMiddleware(
			recovery.Recovery(),
			circuitbreaker.Client(),
		),
		trgrpc.WithOptions(
			grpc.WithInitialWindowSize(grpcInitialWindowSize),
			grpc.WithInitialWindowSize(grpcInitialWindowSize),
			grpc.WithInitialConnWindowSize(grpcInitialConnWindowSize),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcMaxCallMsgSize)),
			grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(grpcMaxSendMsgSize)),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                grpcKeepAliveTime,
				Timeout:             grpcKeepAliveTimeout,
				PermitWithoutStream: true,
			}),
		),
		//grpc.WithTransportCredentials(insecure.NewCredentials()), //todo 安全验证
	)
	if err != nil {
		panic(err)
	}
	return logic.NewLogicClient(conn)
}
