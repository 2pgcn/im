package event

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/2pgcn/gameim/api/protocol"
	"github.com/2pgcn/gameim/conf"
	"github.com/2pgcn/gameim/pkg/gamelog"
	"github.com/2pgcn/gameim/pkg/safe"
	"github.com/nsqio/go-nsq"
	"github.com/pkg/errors"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var FullError = fmt.Errorf("queue is full")

// NsqSender todo 增加连接池,查看源码后发现为单连接....
type NsqSender struct {
	ctx                context.Context
	q                  *nsq.Producer
	topic              string
	senderQueue        []*SendMsgChans
	nsqAsyncReturnChan chan *nsq.ProducerTransaction
	gopool             *safe.GoPool
}

// NewNsqSender todo 仅测试,produce与nsqd一对一,需要修改代码添加容错
func NewNsqSender(c *conf.Nsq) (*NsqSender, error) {
	//todo 添加到配置中
	//defaultQueueNum = runtime.NumCPU()
	config := nsq.NewConfig()
	config.MaxBackoffDuration = time.Second * 10
	p, err := nsq.NewProducer(c.GetNsqdAddress()[0], config)
	if err != nil {
		return nil, err
	}
	if err = p.Ping(); err != nil {
		return nil, err
	}

	p.SetLogger(gamelog.GetGlobalNsqLog(), nsq.LogLevelDebug)
	var sendQueue []*SendMsgChans
	gopool := safe.NewGoPool(context.Background(), "NsqSender")
	nsqAsyncReturnChan := make(chan *nsq.ProducerTransaction, defaultQueueLen/10)
	for i := 0; i < defaultQueueNum; i++ {
		sendQueue = append(sendQueue, &SendMsgChans{
			chs:  make(chan Event, defaultQueueLen),
			msgs: make(map[string]Event, defaultQueueLen),
			lens: 0,
			lock: sync.RWMutex{},
		})
	}
	for _, v1 := range sendQueue {
		v := v1
		gopool.GoCtx(func(ctx context.Context) {
			for {
				select {
				case msg := <-v.chs:
					if msg.IsClose() {
						gamelog.GetGlobalog().Debug("recv close msg,close")
						return
					}
					atomic.AddInt64(&v.lens, -1)
					err = p.PublishAsync(c.Topic, msg.Value(), nsqAsyncReturnChan, msg.GetId())
					if err != nil {
						gamelog.GetGlobalog().Errorf("send PublishAsync error:%s", err)
						return
					}

					////有消息判断长度,最大为len(v.chs)/20
					//maxLen := atomic.LoadInt64(&v.lens)
					//benchMsgLen := defaultQueueLen / 10
					//benMsgs := [][]byte{msg.Value()}
					//args := []string{msg.GetId()}
					//i := 0
					//if maxLen <= 0 {
					//	break
					//}
					//for i < benchMsgLen && int64(i) < maxLen {
					//	m := <-v.chs
					//	if m.IsClose() {
					//		break
					//	}
					//	benMsgs = append(benMsgs, m.Value())
					//	args = append(args, fmt.Sprintf("i%si", m.GetId()))
					//	i++
					//}
					//atomic.AddInt64(&v.lens, -1*int64(max(benchMsgLen, int(maxLen))))
					//if len(benMsgs) > 0 {
					//	err = p.MultiPublishAsync(c.Topic, benMsgs, nsqAsyncReturnChan, args)
					//	if err != nil {
					//		gamelog.GetGlobalog().Errorf("nsq conn err:%s", err)
					//		return
					//	}
					//}

				case _ = <-nsqAsyncReturnChan:
					//if res.Error != nil {
					//	gamelog.GetGlobalog().Error(res)
					//}
					//for _, vArg := range res.Args {
					//	id := vArg.(string)
					//	v.lock.Lock()
					//	delete(v.msgs, id)
					//	v.lock.Unlock()
					//
					//}
				}
			}
		})
	}
	return &NsqSender{
		topic:              c.Topic,
		q:                  p,
		gopool:             gopool,
		senderQueue:        sendQueue,
		nsqAsyncReturnChan: nsqAsyncReturnChan,
	}, nil
}

// todo 可更改负载均衡配置,若该chan满了,其余chan负载也挺高
func (s *NsqSender) Send(ctx context.Context, msg Event) error {
	index := rand.Intn(defaultQueueNum)
	if atomic.LoadInt64(&s.senderQueue[index].lens) > int64(defaultQueueLen) {
		return errors.Wrap(FullError, fmt.Sprintf("error:chan is full:%d", defaultQueueLen))
	}
	s.senderQueue[index].chs <- msg
	s.senderQueue[index].lock.Lock()
	s.senderQueue[index].msgs[msg.GetId()] = msg
	s.senderQueue[index].lock.Unlock()
	atomic.AddInt64(&s.senderQueue[index].lens, 1)
	return nil
}

// 阻塞到完成
func (s *NsqSender) Close() error {
	//先关闭写
	for _, v1 := range s.senderQueue {
		v := v1
		v.chs <- GetCloseMsg()
		close(v.chs)
	}
	gamelog.GetGlobalog().Infof("NsqSender.Close() s.q.Stop() start close")
	s.q.Stop()
	gamelog.GetGlobalog().Infof("NsqSender.Close() start close")
	s.gopool.Stop()
	gamelog.GetGlobalog().Infof("NsqSender.Close() is closed")
	return nil
}

type NsqReceiver struct {
	queueLen int
	q        *nsq.Consumer
	//减少锁力度
	ack       []chan Event
	msgHandel *MessageHandler
	gopool    *safe.GoPool
}

func NewNsqReceiver(c *conf.Nsq) (r *NsqReceiver, err error) {
	consumerConfig := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(c.GetTopic(), c.GetChannel(), consumerConfig)
	//	consumer.SetLogger(gamelog.GetGlobalNsqLog(), nsq.LogLevelWarning)
	if err != nil {
		return nil, err
	}
	consumer.ChangeMaxInFlight(2048)
	r = &NsqReceiver{}
	r.gopool = safe.NewGoPool(context.Background(), "nsq-receiver")

	r.msgHandel = &MessageHandler{
		queueNum:      defaultQueueNum,
		receiverQueue: []chan Event{},
	}
	//todo 加到配置里
	r.queueLen = 1024
	for i := 0; i < r.msgHandel.queueNum; i++ {
		r.msgHandel.receiverQueue = append(r.msgHandel.receiverQueue, make(chan Event, r.queueLen))
		r.ack = append(r.ack, make(chan Event, r.queueLen))
	}
	consumer.AddHandler(r.msgHandel)
	if len(c.GetLookupd()) > 0 {
		err = consumer.ConnectToNSQLookupds(c.GetLookupd())
		if err != nil {
			return nil, err
		}
	} else {
		err = consumer.ConnectToNSQDs(c.GetNsqdAddress())
		if err != nil {
			return nil, err
		}
	}
	r.q = consumer
	return r, nil
}

type MessageHandler struct {
	queueNum      int
	receiverQueue []chan Event
}

func (r *MessageHandler) HandleMessage(message *nsq.Message) error {
	msg := &QueueMsg{H: make(map[string]string, 8), Data: &protocol.Msg{}}
	if err := json.Unmarshal(message.Body, &msg); err != nil {
		gamelog.GetGlobalog().Error(err)
		return err
	}
	//message.DisableAutoResponse()
	//msg.H["id"] = message.ID
	index := rand.Intn(r.queueNum)
	r.receiverQueue[index] <- msg
	//todo 参考生产者
	return nil
}

func (r *NsqReceiver) Receive(ctx context.Context) (e []chan Event, err error) {
	return r.msgHandel.receiverQueue, nil
}

func (r *NsqReceiver) RecvCommit(ctx context.Context) (e []chan Event, err error) {
	return r.ack, nil
}

func (r *NsqReceiver) Commit(ctx context.Context, event Event) error {
	if len(r.ack) == r.queueLen {
		return errors.Wrap(FullError, fmt.Sprintf("error:chan is full:%d", r.queueLen))
	}
	r.ack[rand.Intn(r.msgHandel.queueNum)] <- event
	return nil
}

func (r *NsqReceiver) Close() error {
	r.q.Stop()

	<-r.q.StopChan
	for _, v := range r.msgHandel.receiverQueue {
		close(v)
	}
	for _, v := range r.ack {
		close(v)
	}
	return nil
}
