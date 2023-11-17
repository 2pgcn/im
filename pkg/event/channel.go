package event

import (
	"context"
	"sync"
)

// todo 抽象出来为了方便后续链路追踪信息注入
type channel struct {
	ch  chan Event
	one sync.Once
}

func NewChannel(size int) *channel {
	return &channel{
		ch:  make(chan Event, size),
		one: sync.Once{},
	}

}

func (ch *channel) Send(ctx context.Context, msg Event) error {
	ch.ch <- msg
	return nil
}

// 只发送端关闭,接收端不关闭
func (ch *channel) Close() error {
	ch.one.Do(func() {
		close(ch.ch)
	})
	return nil
}

func (ch *channel) Receive(ctx context.Context) (e Event, err error) {
	e = <-ch.ch
	return
}

func (ch *channel) Commit(ctx context.Context, event Event) error {
	return nil
}
