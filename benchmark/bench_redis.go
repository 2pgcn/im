package main

import (
	"context"
	"github.com/2pgcn/gameim/api/protocol"
	"github.com/2pgcn/gameim/conf"
	"github.com/2pgcn/gameim/pkg/event"
	"github.com/2pgcn/gameim/pkg/gamelog"
	"google.golang.org/protobuf/types/known/durationpb"
	"strconv"
	"sync/atomic"
	"time"
)

var redisConf = &conf.Queue{
	Topic:        "benchRedis",
	Channel:      "benchRedis",
	RecvQueueLen: 10240,
	RecvQueueNum: 8,
	SendQueueLen: 10240,
	SendQueueNum: 8,
}

func benchRedis(ctx context.Context, c *Redis) {
	var Rdscfg = &conf.Redis{
		Addr:         c.RedisAddress,
		Passwd:       "",
		ReadTimeout:  durationpb.New(time.Second * 3),
		WriteTimeout: durationpb.New(time.Second * 3),
		KeepLive:     durationpb.New(time.Minute * 2),
		PoolSize:     5,
		UseDb:        0,
	}
	redisConf.Topic = c.Topic
	redisConf.Channel = c.Channel
	producer := event.NewRedisSender(&conf.RedisQueue{
		Rdscfg: Rdscfg,
		Queue:  redisConf,
	})
	err := producer.Start()
	if err != nil {
		panic(err)
	}
	consumer, err := event.NewRedisReceiver(&conf.RedisQueue{
		Rdscfg: Rdscfg,
		Queue:  redisConf,
	})
	if err != nil {
		panic(err)
	}
	if err = consumer.Start(); err != nil {
		panic(err)
	}
	e, err := consumer.Receive(ctx)
	if err != nil {
		panic(err)
	}
	for _, v1 := range e {
		v := v1
		gopool.GoCtx(func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				case _ = <-v:
					atomic.AddInt64(&countDown, 1)
				}
			}
		})
	}
	var msgId = 0
	gopool.GoCtx(func(ctx context.Context) {
		for {
			for i := 0; i <= 100000; i++ {
				queueMsg := &event.QueueMsg{
					H:    make(map[string]string, 8),
					Data: &protocol.Msg{},
				}
				queueMsg.Data.Msg = []byte("hello world")
				queueMsg.Data.Type = protocol.Type_APP
				queueMsg.Data.ToId = "0"
				queueMsg.Data.SendId = "0"
				msgId++
				queueMsg.SetId(strconv.Itoa(msgId))
				select {
				case <-ctx.Done():
					gamelog.GetGlobalog().Debug("exit producer start stop")
					producer.Close()
					return
				default:
					err = producer.Send(ctx, queueMsg)
					if err != nil {
						gamelog.GetGlobalog().Error(err)
						continue
					}
					atomic.AddInt64(&countSend, 1)
				}
			}
			time.Sleep(time.Second * 1)
		}
	})
	gopool.GoCtx(func(ctx context.Context) {
		select {
		case <-ctx.Done():
			consumer.Close()
		}
	})
	gopool.GoCtx(func(ctx context.Context) {
		result(ctx)
	})
}
