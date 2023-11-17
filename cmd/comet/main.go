package main

import (
	"context"
	"fmt"
	"github.com/2pgcn/gameim/conf"
	"github.com/2pgcn/gameim/internal/comet"
	"github.com/2pgcn/gameim/pkg/gamelog"
	"github.com/2pgcn/gameim/pkg/trace_conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	_ "go.uber.org/automaxprocs"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
)

var rootCmd = &cobra.Command{
	Use:     "game im",
	Short:   "game im comet",
	Long:    `implementing game im comet in go`,
	Version: "0.0.1",
	Run: func(cmd *cobra.Command, args []string) {
		cometConfig := conf.InitCometConfig(CfgFile)
		trace_conf.SetTraceConfig(cometConfig.TraceConf)
		if port := os.Getenv("ILOGTAIL_PROFILE_PORT"); len(port) > 0 {
			startPprof(fmt.Sprintf(":", port))
		}
		if err := startTrace(); err != nil {
			panic(err)
		}

		wg := &sync.WaitGroup{}
		ctx, cancel := context.WithCancel(context.Background())
		l := gamelog.GetZapLog()
		zlog := gamelog.NewHelper(l)
		go func() {
			zlog.Debug(http.ListenAndServe("0.0.0.0:8888", nil))
		}()
		s, err := comet.NewServer(ctx, cometConfig, wg, zlog)
		if err != nil {
			panic(err)
		}
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, os.Kill)
		select {
		case <-signals:
			cancel()
			s.Close()
		}
		wg.Wait()
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

func startPprof(port string) {
	go func() {
		_ = http.ListenAndServe(port, nil)
	}()
}
