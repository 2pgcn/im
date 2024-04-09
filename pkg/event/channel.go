package event

import (
	"context"
	"math/rand"
	"sync"
)

var queueChannelNum = 1

// todo 抽象出来为了方便后续链路追踪信息注入
type Channel struct {
	queueNum int
	ch       []chan Event
	one      sync.Once
}

func NewChannel(size int) *Channel {
	return &Channel{
		queueNum: queueChannelNum,
		ch:       []chan Event{make(chan Event, size)},
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
			v1 := v
			close(v1)
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
