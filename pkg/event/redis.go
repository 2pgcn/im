package event

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/2pgcn/gameim/api/gerr"
	"github.com/2pgcn/gameim/api/protocol"
	"github.com/2pgcn/gameim/conf"
	"github.com/2pgcn/gameim/pkg/gamelog"
	"github.com/2pgcn/gameim/pkg/safe"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"sync/atomic"
	"time"
)

var defaultPoolName = "RedisSender"
var closeInt = 1

type RedisSender struct {
	qlen        int
	ctx         context.Context
	topic       string
	client      *redis.Client
	q           redis.Pipeliner
	senderQueue []*SendMsgChans
	gopool      *safe.GoPool
}

// NewRedisSender 仅测试使用,生产环境不应该使用redis作为消息队列,消费端出问题内存会oom,以及缺少很多队列功能
func NewRedisSender(c *conf.RedisQueue) *RedisSender {
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.GetRdscfg().GetAddr(),
		Password:     c.GetRdscfg().GetPasswd(),
		DB:           0,
		PoolSize:     int(c.GetRdscfg().GetPoolSize()),
		ReadTimeout:  c.GetRdscfg().GetReadTimeout().AsDuration(),
		WriteTimeout: c.GetRdscfg().GetWriteTimeout().AsDuration(),
	})
	res := rdb.Ping(context.Background())
	if err := res.Err(); err != nil {
		gamelog.GetGlobalog().Error("redis ping err", err)
		return nil
	}
	r := &RedisSender{
		ctx:         context.Background(),
		topic:       fmt.Sprintf("%s-%s", c.GetQueue().GetTopic(), c.GetQueue().GetChannel()),
		q:           rdb.Pipeline(),
		client:      rdb,
		senderQueue: make([]*SendMsgChans, 0),
		gopool:      safe.NewGoPool(context.Background(), defaultPoolName),
	}
	queueLen := max(int(c.Queue.SendQueueLen), defaultQueueLen)
	queueNum := max(int(c.Queue.SendQueueNum), defaultQueueNum)
	r.qlen = queueLen
	for i := 0; i < queueNum; i++ {
		r.senderQueue = append(r.senderQueue, NewSendMsgChans(queueLen, queueLen/10))
	}
	return r
}

func (r *RedisSender) Start() error {
	if r.q == nil {
		return gerr.ErrorServerError("RedisSender init redis client error").WithMetadata(gerr.GetStack())
	}
	for _, v1 := range r.senderQueue {
		v := v1
		r.gopool.GoCtx(func(ctx context.Context) {
			for {
				select {
				case msg := <-v.chs:
					if msg.IsClose() {
						gamelog.GetGlobalog().Debug("recv close msg,close")
						return
					}
					atomic.AddInt64(&v.lens, -1)
					r.q.LPush(ctx, r.topic, string(msg.Value()))
					n := min(r.qlen/5, int(atomic.LoadInt64(&v.lens))) / 3
					//msgs := make([]any, n+1)
					//msgs[0] = string(msg.Value())
					for i := 1; i < n+1; i++ {
						//批量发送queue
						msg = <-v.chs
						r.q.LPush(ctx, r.topic, string(msg.Value()))
						//msgs[i] = string(msg.Value())
						atomic.AddInt64(&v.lens, -1)
					}
					//res := r.q.LPushX(r.ctx, r.topic, msgs...)
					//if res.Err() != nil {
					//	gamelog.GetGlobalog().Error("redis sender send msg error", res.Err())
					//	continue
					//}
					ctx, _ := context.WithTimeout(ctx, time.Second*5)
					_, err := r.q.Exec(ctx)
					if err != nil {
						//todo 发送失败,重新入队,仅测试未处理异常
						gamelog.GetGlobalog().Error(err)
						continue
					}
					//for _, v11 := range t {
					//	gamelog.GetGlobalog().Debug(v11)
					//}
				}
			}
		})
	}
	return nil
}
func (s *RedisSender) Send(ctx context.Context, msg Event) error {
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

func (s *RedisSender) Close() error {
	//先关闭写
	for _, v1 := range s.senderQueue {
		v := v1
		v.chs <- GetCloseMsg()
		close(v.chs)
	}
	s.gopool.Stop()
	gamelog.GetGlobalog().Debugf("RedisSender.Close() s.q.Stop() start close")
	err := s.client.Close()
	if err != nil {
		return gerr.ErrorServerError("RedisSender close error").WithCause(err).WithMetadata(gerr.GetStack())
	}
	gamelog.GetGlobalog().Infof("RedisSender.Close() is closed")
	return nil
}

type RedisReceiver struct {
	qLen  int
	q     *redis.Client
	topic string
	//减少锁力度
	ack       []chan Event
	msgHandel *MessageHandler
	gopool    *safe.GoPool
	isClose   int32
}

func NewRedisReceiver(c *conf.RedisQueue) (r *RedisReceiver, err error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.GetRdscfg().GetAddr(),
		Password:     c.GetRdscfg().GetPasswd(),
		DB:           int(c.GetRdscfg().GetUseDb()),
		PoolSize:     int(c.GetRdscfg().GetPoolSize()),
		ReadTimeout:  c.GetRdscfg().GetReadTimeout().AsDuration(),
		WriteTimeout: c.GetRdscfg().GetWriteTimeout().AsDuration(),
	})
	r = &RedisReceiver{
		topic: fmt.Sprintf("%s-%s", c.GetQueue().GetTopic(), c.GetQueue().GetChannel()),
	}
	r.msgHandel = &MessageHandler{
		queueNum:      defaultQueueNum,
		receiverQueue: []chan Event{},
	}
	r.gopool = safe.NewGoPool(context.Background(), "nsq-receiver")

	//todo 加到配置里
	r.qLen = 1024
	for i := 0; i < r.msgHandel.queueNum; i++ {
		r.msgHandel.receiverQueue = append(r.msgHandel.receiverQueue, make(chan Event, r.qLen))
		r.ack = append(r.ack, make(chan Event, r.qLen))
	}
	r.q = rdb
	err = r.Start()
	return r, err
}

func (r *RedisReceiver) Start() error {
	//每次取min 1/5
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	if r.q == nil {
		return gerr.ErrorServerError("start RedisReceiver error")
	}
	if res := r.q.Ping(ctx); res.Err() != nil {
		return gerr.ErrorServerError("start RedisReceiver error").WithCause(res.Err())
	}
	r.gopool.GoCtx(func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if atomic.LoadInt32(&r.isClose) == int32(closeInt) {
					return
				}
				ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
				vals, err := r.q.LPopCount(ctx, r.topic, r.qLen/5).Result()
				if err != nil {
					if errors.Is(err, redis.Nil) {
						time.Sleep(time.Second * 1)
						continue
					}

					gamelog.GetGlobalog().Error(err)
					continue
				}
				for _, tem := range vals {
					msg := &QueueMsg{H: make(map[string]string, 8), Data: &protocol.Msg{}}
					if err := json.Unmarshal([]byte(tem), &msg); err != nil {
						gamelog.GetGlobalog().Infof("data:%s,err:%s", string(tem), err)
						continue
					}
					r.msgHandel.receiverQueue[rand.Intn(r.msgHandel.queueNum)] <- msg
				}

			}

		}
	})
	return nil
}

func (r *RedisReceiver) Receive(ctx context.Context) (e []chan Event, err error) {
	return r.msgHandel.receiverQueue, nil
}

func (r *RedisReceiver) RecvCommit(ctx context.Context) (e []chan Event, err error) {
	return r.ack, nil
}

func (r *RedisReceiver) Commit(ctx context.Context, event Event) error {
	if len(r.ack) == r.qLen {
		return errors.Wrap(FullError, fmt.Sprintf("error:chan is full:%d", r.qLen))
	}
	r.ack[rand.Intn(r.msgHandel.queueNum)] <- event
	return nil
}

func (r *RedisReceiver) Close() error {
	atomic.SwapInt32(&r.isClose, int32(closeInt))
	for _, v1 := range r.msgHandel.receiverQueue {
		v := v1
		v <- GetCloseMsg()
		close(v)
	}
	err := r.q.Close()
	if err != nil {
		return gerr.ErrorServerError("RedisReceiver close error").WithCause(err).WithMetadata(gerr.GetStack())
	}
	for _, v := range r.ack {
		close(v)
	}
	return nil
}
