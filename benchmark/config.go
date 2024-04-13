package main

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
)

func getBenchConfig() *BenchConf {
	c := config.New(
		config.WithSource(
			file.NewSource(cpath),
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
