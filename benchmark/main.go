package main

import (
	"context"
	"errors"
	"github.com/2pgcn/gameim/api/protocol"
	"github.com/2pgcn/gameim/conf"
	"github.com/2pgcn/gameim/pkg/event"
	"github.com/2pgcn/gameim/pkg/gamelog"
	"github.com/2pgcn/gameim/pkg/safe"
	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"time"
)

var gopool *safe.GoPool
var c *BenchConf

var address string
var num int

func initLog() {
	l := gamelog.GetZapLog(zapcore.DebugLevel, 0)
	_ = gamelog.NewHelper(l)
}

func main() {
	go func() {
		_ = http.ListenAndServe("0.0.0.0:9999", nil)
	}()
	ctx, cannel := context.WithCancel(context.Background())
	gopool = safe.NewGoPool(ctx)
	initLog()
	cmd := NewServerArgs()
	cmd.SetContext(ctx)
	cmd.PersistentFlags().StringVar(&address, "address", "", "address eg:127.0.0.1:9000")
	cmd.PersistentFlags().StringVar(&flagconf, "conf", "/Users/pg/work/go/src/github.com/2pgcn/gameim/benchmark/conf.yaml", "address eg:10")
	cmd.PersistentFlags().IntVar(&num, "num", 1, "address eg:10")
	c = getBenchConfig()
	//pprof.InitPyroscope("benchmark", "1.0.0", c.GetPyroscope().GetAddress(), gamelog.GetGlobalog())

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Kill, os.Interrupt)
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
	select {
	case <-exitChan:
		cannel()
		gopool.Stop()
	}

}

func NewServerArgs() *cobra.Command {
	return &cobra.Command{
		Short: "cmd",
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "nsq":
				benchNsq(cmd.Context())
				break
			case "logic":
				benchLogic(cmd.Context(), address, num)
				break
			case "comet":
				benchComet(cmd.Context(), address, num)
			}
		},
	}
}

func benchNsq(ctx context.Context) {
	consumer, err := event.NewNsqReceiver(&conf.QueueMsg_Nsq{
		Topic:       c.GetNsq().GetTopic(),
		Channel:     c.GetNsq().GetChannel(),
		NsqdAddress: c.GetNsq().GetNsqdAddress(),
	})
	if err != nil {
		panic(err)
	}
	e, err := consumer.Receive(ctx)
	if err != nil {
		panic(err)
	}
	producer, err := event.NewNsqSender(&conf.Data_Nsq{
		Topic:   c.GetNsq().GetTopic(),
		Channel: c.GetNsq().GetChannel(),
		Address: c.GetNsq().GetNsqdAddress(),
	})
	if err != nil {
		panic(err)
	}
	for _, v1 := range e {
		v := v1
		gopool.GoCtx(func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				case msg := <-v:
					atomic.AddInt64(&countDown, 1)
					event.PutQueueMsg(msg.GetQueueMsg())
					//err = consumer.Commit(ctx, msg)
					//if err != nil {
					//	fmt.Println(err)
					//}
				}
			}
		})
	}
	var msgId = 0
	gopool.GoCtx(func(ctx context.Context) {
		for {
			for i := 0; i <= 10000; i++ {
				queueMsg := &event.QueueMsg{
					H:    make(map[string]string, 8),
					Data: &protocol.Msg{},
				}
				queueMsg.Data.Msg = []byte("hello world")
				queueMsg.Data.Type = protocol.Type_APP
				queueMsg.Data.ToId = "0"
				queueMsg.Data.SendId = "0"
				msgId++
				queueMsg.SetId(strconv.Itoa(msgId))
				select {
				case <-ctx.Done():
					gamelog.Debug("exit producer start stop")
					producer.Close()
					return
				default:
					err = producer.Send(ctx, queueMsg)
					if err != nil {
						if errors.Is(err, event.FullError) {
							time.Sleep(time.Second * 1)
							continue
						}
						gamelog.GetGlobalog().Error(err)
						continue
					}
					atomic.AddInt64(&countSend, 1)
				}
			}
			time.Sleep(time.Second * 1)
		}
	})
	gopool.GoCtx(func(ctx context.Context) {
		select {
		case <-ctx.Done():
			consumer.Close()
		}
	})
	gopool.GoCtx(func(ctx context.Context) {
		result(ctx)
	})
}
