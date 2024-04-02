package main

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
)

var flagconf string

func getBenchConfig() *BenchConf {
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()
	if err := c.Load(); err != nil {
		panic(err)
	}
	var bc BenchConf
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	return &bc
}
