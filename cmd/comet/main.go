package main

import (
	"context"
	"fmt"
	"github.com/php403/gameim/config"
	"github.com/php403/gameim/internal/comet"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"log"
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
	sugaredLogger := zap.NewExample().Sugar()
	comet.NewServer(context.Background(), sugaredLogger, cometConfig)
}
