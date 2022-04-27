package config

import (
	"fmt"
	"github.com/spf13/viper"
	_ "google.golang.org/protobuf/reflect/protoreflect"
)

func InitCometConfig(CfgFile string) (cometConfig *CometConfig) {
	if len(CfgFile) == 0 {
		panic(fmt.Errorf("config file %s not found", CfgFile))
	}

	viper.AddConfigPath("/Users/pg/work/go/gopath/src/github.com/gameim/config")
	viper.SetConfigName("comet")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.UnmarshalKey("Comet", &cometConfig); err != nil {
		panic(err)
	}
	return
}
