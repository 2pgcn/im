package event

import (
	"context"
	"math/rand"
	"sync"
)

var queueChannelNum = 12

// todo 抽象出来为了方便后续链路追踪信息注入
type Channel struct {
	queueNum int
	ch       []chan Event
	one      sync.Once
}

func NewChannel(size int) *Channel {
	chs := make([]chan Event, queueChannelNum)
	for i := 0; i < queueChannelNum; i++ {
		chs[i] = make(chan Event, size)
	}
	return &Channel{
		queueNum: queueChannelNum,
		ch:       chs,
		one:      sync.Once{},
	}

}

// Send todo 统一queue格式 类似nsq queue实现
func (ch *Channel) Send(ctx context.Context, msg Event) error {
	ch.ch[rand.Intn(ch.queueNum)] <- msg
	return nil
}

// 只发送端关闭,接收端不关闭
func (ch *Channel) Close() error {
	ch.one.Do(func() {
		for _, v := range ch.ch {
			close(v)
		}
	})
	return nil
}

func (ch *Channel) Receive(ctx context.Context) (e []chan Event, err error) {
	return ch.ch, nil
}

func (ch *Channel) Commit(ctx context.Context, event Event) error {
	return nil
}
