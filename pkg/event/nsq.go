package event

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/2pgcn/gameim/conf"
	"github.com/2pgcn/gameim/pkg/gamelog"
	"github.com/2pgcn/gameim/pkg/safe"
	"github.com/nsqio/go-nsq"
	"math/rand"
)

type NsqSender struct {
	topic string
	q     *nsq.Producer
}

// NewNsqSender todo 仅测试,produce与nsqd一对一,需要修改代码添加容错
func NewNsqSender(c *conf.Data_Nsq) (*NsqSender, error) {
	config := nsq.NewConfig()
	p, err := nsq.NewProducer(c.GetAddress()[0], config)
	if err != nil {
		return nil, err
	}
	p.SetLogger(gamelog.GetGlobalNsqLog(), nsq.LogLevelWarning)
	return &NsqSender{
		topic: c.Topic,
		q:     p,
	}, nil
}

func (s *NsqSender) Send(ctx context.Context, msg Event) error {
	return s.q.PublishAsync(s.topic, msg.Value(), make(chan *nsq.ProducerTransaction, 16))
}

// 阻塞到完成
func (s *NsqSender) Close() error {
	s.q.Stop()
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

func NewNsqReceiver(c *conf.QueueMsg_Nsq) (r *NsqReceiver, err error) {
	consumerConfig := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(c.GetTopic(), c.GetChannel(), consumerConfig)
	consumer.SetLogger(gamelog.GetGlobalNsqLog(), nsq.LogLevelWarning)
	if err != nil {
		return nil, err
	}
	consumer.ChangeMaxInFlight(10240)
	r = &NsqReceiver{}
	r.gopool = safe.NewGoPool(context.Background())

	r.msgHandel = &MessageHandler{
		queueNum:      8,
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
	msg := GetQueueMsg()
	if err := json.Unmarshal(message.Body, &msg); err != nil {
		return err
	}
	//message.DisableAutoResponse()
	//msg.H["id"] = message.ID
	index := rand.Intn(r.queueNum)
	r.receiverQueue[index] <- msg
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
		return fmt.Errorf("error:chan is full:%d", r.queueLen)
	}
	r.ack[rand.Intn(r.msgHandel.queueNum)] <- event
	return nil
}

func (r *NsqReceiver) Close() error {
	r.q.Stop()
	for _, v := range r.msgHandel.receiverQueue {
		close(v)
	}
	for _, v := range r.ack {
		close(v)
	}
	return nil
}
