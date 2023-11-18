package event

import (
	"context"
	"sync"
)

// todo 抽象出来为了方便后续链路追踪信息注入
type Channel struct {
	ch  chan Event
	one sync.Once
}

func NewChannel(size int) *Channel {
	return &Channel{
		ch:  make(chan Event, size),
		one: sync.Once{},
	}

}

func (ch *Channel) Send(ctx context.Context, msg Event) error {
	ch.ch <- msg
	return nil
}

// 只发送端关闭,接收端不关闭
func (ch *Channel) Close() error {
	ch.one.Do(func() {
		close(ch.ch)
	})
	return nil
}

func (ch *Channel) Receive(ctx context.Context) (e Event, err error) {
	e = <-ch.ch
	return
}

func (ch *Channel) Commit(ctx context.Context, event Event) error {
	return nil
}
