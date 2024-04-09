package main

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

// 展示数据使用
var (
	countDown  int64
	countSend  int64
	aliveCount int64
)

func addCountDown(n int64) {
	atomic.AddInt64(&countDown, n)
}

func addCountSend(n int64) {
	atomic.AddInt64(&countSend, n)
}
func addAliveCount(n int64) {
	atomic.AddInt64(&aliveCount, n)
}

func result(ctx context.Context) {
	var (
		lastDown int64
		lastSend int64
		interval = int64(3)
	)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			nowDownCount := atomic.LoadInt64(&countDown)
			nowSendCount := atomic.LoadInt64(&countSend)
			nowAlive := atomic.LoadInt64(&aliveCount)
			diffDown := nowDownCount - lastDown
			diffSend := nowSendCount - lastSend
			lastDown = nowDownCount
			lastSend = nowSendCount
			fmt.Println(fmt.Sprintf("%s alive:%d down:%d down/s:%d", time.Now().Format("2006-01-02 15:04:05"), nowAlive, nowDownCount, diffDown/interval))
			fmt.Println(fmt.Sprintf("%s alive:%d send:%d send/s:%d", time.Now().Format("2006-01-02 15:04:05"), nowAlive, nowSendCount, diffSend/interval))
			time.Sleep(time.Second * time.Duration(interval))
		}

	}
}
