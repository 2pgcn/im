package main

import (
	"context"
	"fmt"
	"github.com/2pgcn/gameim/api/client"
	"github.com/2pgcn/gameim/api/logic"
	"github.com/2pgcn/gameim/api/protocol"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func benchLogic(ctx context.Context, address string, num int) {
	kDis, err := client.NewK8sDiscovery("")
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	//默认k8s注册
	cc, err := client.NewGrpcClient(ctx, address, kDis)
	if err != nil {
		panic(err)
	}
	gclient := logic.NewLogicClient(cc)
	var lasterr error
	ticker := time.Tick(time.Second * 1)
	gopool.GoCtx(func(ctx context.Context) {
		t := time.Tick(time.Second * 3)
		for {
			select {
			case <-ctx.Done():
				return
			case <-t:
				fmt.Println(cc.GetState())
			}
		}
	})
	gopool.GoCtx(func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker:
				for i := 0; i < num; i++ {
					addCountSend(1)
					_, err := gclient.OnMessage(ctx, &logic.MessageReq{
						Type:     protocol.Type_PUSH,
						CometKey: "test",
						ToId:     "1",
						SendId:   strconv.Itoa(i),
						Msg:      []byte("hello world"),
					})
					if err != nil {
						fmt.Println(cc.GetState())
						fmt.Printf("lasterr:%s", err)
						continue
					}
					addCountDown(1)
				}
				if lasterr != nil {

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
