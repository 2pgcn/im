package main

import (
	"context"
	"github.com/2pgcn/gameim/config"
	"github.com/2pgcn/gameim/internal/comet"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
)

var rootCmd = &cobra.Command{
	Use:     "goim",
	Short:   "game im comet",
	Long:    `implementing game im comet in go`,
	Version: "0.0.1",
	Run: func(cmd *cobra.Command, args []string) {
		//server := comet.NewServer(context.Background(), config.GetConfig())
	},
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:8888", nil))
	}()
	var CfgFile string
	rootCmd.PersistentFlags().StringVar(&CfgFile, "config", "./config/comet.yaml", "config file (default is $HOME/.cobra.yaml)")
	devLog, _ := zap.NewProduction()
	sugaredLogger := devLog.Sugar()
	cometConfig := config.InitCometConfig(CfgFile)
	sugaredLogger.Debug(cometConfig)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)
	comet.NewServer(ctx, sugaredLogger, cometConfig)
	for {
		select {
		case <-signals:
			cancel()
			return
		}
	}
}
