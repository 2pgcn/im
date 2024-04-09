package main

import (
	"context"
	"fmt"
	"github.com/2pgcn/gameim/api/client"
	"github.com/2pgcn/gameim/api/logic"
	"github.com/2pgcn/gameim/api/protocol"
	"github.com/2pgcn/gameim/conf"
	"github.com/2pgcn/gameim/pkg/event"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func benchLogic(ctx context.Context, address string, num int) {
	//kDis, err := client.NewK8sDiscovery("")
	//if err != nil {
	//	panic(err)
	//}
	err := StartTestSockRecv("../sock_queue.sock")
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 5)
	ctx, cancel := context.WithCancel(context.Background())
	//默认k8s注册
	cc, err := client.NewGrpcClient(ctx, address, nil)
	if err != nil {
		panic(err)
	}
	gclient := logic.NewLogicClient(cc)
	var lasterr error

	gopool.GoCtx(func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
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

// 模拟sock接收,方便压测logic
func StartTestSockRecv(address string) error {
	rcvQueue, err := event.NewSockReceiver(&conf.Sock{
		Address: address,
	})
	if err != nil {
		return err
	}
	e, err := rcvQueue.Receive(context.Background())
	for _, v1 := range e {
		v := v1
		gopool.GoCtx(func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				case _ = <-v:
					//atomic.AddInt64(&countDown, 1)
					//event.PutQueueMsg(msg.GetQueueMsg())
					//err = consumer.Commit(ctx, msg)
					//if err != nil {
					//	fmt.Println(err)
					//}
				}
			}
		})
	}
	return nil
}
