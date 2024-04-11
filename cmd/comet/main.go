package main

import (
	"context"
	"github.com/2pgcn/gameim/api/protocol"
	"github.com/2pgcn/gameim/conf"
	"github.com/2pgcn/gameim/internal/comet"
	"github.com/2pgcn/gameim/pkg/gamelog"
	"github.com/2pgcn/gameim/pkg/pprof"
	"github.com/2pgcn/gameim/pkg/safe"
	"github.com/2pgcn/gameim/pkg/trace_conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/grafana/pyroscope-go"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap/zapcore"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
)

var (
	Name    = "comet-0"
	Version = "1.0.0"
)

var rootCmd = &cobra.Command{
	Use:     "game im",
	Short:   "game im comet",
	Long:    `implementing game im comet in go`,
	Version: strconv.Itoa(int(protocol.Version)),
	Run: func(cmd *cobra.Command, args []string) {
		//trace_conf.SetTraceConfig(cometConfig.UpData.TraceConf)
		//if err := startTrace(); err != nil {
		//	panic(err)
		//}
		//if err := startPyroscope(Name, Version, cometConfig.UpData.Pyroscope.Address, gamelog.GetGlobalog()); err != nil {
		//	panic(err)
		//
		//}
		cometConfig := conf.InitCometConfig(CfgFile)
		//todo add to conf
		l := gamelog.GetZapLog(zapcore.InfoLevel, 2)
		zlog := gamelog.NewHelper(l)
		ctx, cancel := context.WithCancel(context.Background())
		gopool := safe.NewGoPool(ctx, "gameim-comet-main")

		s, err := comet.NewServer(ctx, cometConfig, gopool, zlog)
		if err != nil {
			panic(err)
		}

		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, os.Kill)
		select {
		case <-signals:
			zlog.Debug("stop gameim comet")
			cancel()
			zlog.Debug("send ctx done success")
			s.Close()
			zlog.Debug("server closing")
			gopool.Stop()
		}
	},
}
var CfgFile string

func main() {
	rootCmd.PersistentFlags().StringVar(&CfgFile, "conf", "", "conf file (default is $HOME/.cobra.yaml)")
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func startTrace() error {
	tp, err := trace_conf.GetTracerProvider()
	if err != nil {
		return err
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}))
	return nil
}

func startPyroscope(appname, version, endpoint string, logger pyroscope.Logger) error {
	return pprof.InitPyroscope(appname, version, endpoint, logger)
}
