package main

import (
	"context"
	"github.com/2pgcn/gameim/pkg/gamelog"
	"github.com/2pgcn/gameim/pkg/safe"
	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
)

var gopool *safe.GoPool
var c *BenchConf

var address string
var num int

func initLog() {
	l := gamelog.GetZapLog(zapcore.DebugLevel, 2)
	_ = gamelog.NewHelper(l)
}

func main() {
	go func() {
		_ = http.ListenAndServe("0.0.0.0:9999", nil)
	}()
	ctx, cannel := context.WithCancel(context.Background())
	gopool = safe.NewGoPool(ctx, "main")
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
			case "sock":
				benchSock(cmd.Context())
				break
			case "logic":
				benchLogic(cmd.Context(), address, num)
				break
			case "comet":
				benchComet(cmd.Context(), address, num)
			case "redis":
				benchRedis(cmd.Context(), c.GetRedis())
			}
		},
	}
}
