package main

import (
	"context"
	"fmt"
	"github.com/php403/gameim/config"
	"github.com/php403/gameim/internal/comet"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"log"
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
	var CfgFile string
	rootCmd.PersistentFlags().StringVar(&CfgFile, "config", "./config/comet.yaml", "config file (default is $HOME/.cobra.yaml)")
	cometConfig := config.InitCometConfig(CfgFile)
	fmt.Println(cometConfig)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
	devLog, _ := zap.NewDevelopment()
	sugaredLogger := devLog.Sugar()
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
