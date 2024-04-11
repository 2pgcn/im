package conf

import (
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/spf13/viper"
	_ "google.golang.org/protobuf/reflect/protoreflect"
)

var ServerName = "gameim"

func InitCometConfig(CfgFile string) *CometConfig {
	c := config.New(
		config.WithSource(
			file.NewSource(CfgFile),
		),
	)
	defer c.Close()
	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc CometConfig
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	return &bc
}
func initCometConfig(CfgFile string) (cometConfig *CometConfig) {
	if len(CfgFile) == 0 {
		panic(fmt.Errorf("config file %s not found", CfgFile))
	}
	viper.AddConfigPath(CfgFile)
	viper.SetConfigName("comet.yaml")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	data, err := json.Marshal(viper.AllSettings()["comet"])
	if err != nil {
		panic(err)
	}
	cometConfig = &CometConfig{}
	if err := json.Unmarshal(data, cometConfig); err != nil {
		panic(err)
	}
	return
}
