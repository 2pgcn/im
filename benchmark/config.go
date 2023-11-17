package main

import (
	"flag"
	"github.com/2pgcn/gameim/conf"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
)

var flagconf string

func getLogicConfig() *conf.Bootstrap {
	flag.Parse()
	flag.StringVar(&flagconf, "conf", "../conf/logic.yaml", "config path, eg: -conf config.yaml")
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()
	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	return &bc
}
