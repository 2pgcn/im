package event

import (
	"context"
	"github.com/2pgcn/gameim/conf"
	"github.com/2pgcn/gameim/pkg/safe"
	"github.com/nsqio/go-nsq"
	"sync"
)

type SendMsgChans struct {
	chs chan Event
	//记录chs里有多少消息
	lens int64
	msgs map[string]Event
	lock sync.RWMutex
}

type QueueProducer struct {
	ctx                context.Context
	producer           Sender
	topic              string
	producerChans      []*SendMsgChans
	nsqAsyncReturnChan chan *nsq.ProducerTransaction
	gopool             *safe.GoPool
}

type Consumer struct {
	ctx      context.Context
	consumer Receiver
	//todo 如kafka,每次从最小的message id查找,找到当前最大的message id,改成最小堆实现
	//例如nsq,则可以无序,
	msgHeap map[string]Event
}

func NewQueueProducer(ctx context.Context, c *conf.Queue) (r *QueueProducer, err error) {
	//var producer Sender
	/*	switch c.q {

		}*/
	r = &QueueProducer{
		ctx:                ctx,
		producer:           nil,
		topic:              "",
		producerChans:      nil,
		nsqAsyncReturnChan: nil,
		gopool:             nil,
	}
	return nil, err
}
