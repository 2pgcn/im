package comet

import (
	"context"
	"fmt"
	"github.com/2pgcn/gameim/api/comet"
	"github.com/2pgcn/gameim/config"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"sync"
)

type App struct {
	ctx       context.Context
	Appid     string
	config    *config.CometConfigQueueMsg
	Buckets   []*Bucket //app bucket
	bucketIdx uint32
	lock      sync.RWMutex
	log       *zap.SugaredLogger
}

func (a *App) GetBucketIndex(userid uint64) int {
	return int(userid % uint64(a.bucketIdx))
}

// NewApp todo 暂时写死app 需要改为从存储中获取 config改成app config
func NewApp(ctx context.Context, appid string, config *config.CometConfigQueueMsg, log *zap.SugaredLogger) (*App, error) {
	app := &App{
		ctx:       ctx,
		Appid:     appid,
		Buckets:   make([]*Bucket, 1),
		bucketIdx: 1, //todo 改为配置
		log:       log,
		config:    config,
	}
	for i := 0; i < 1; i++ {
		app.Buckets[i] = NewBucket()
	}
	//consumer, err := NewConsumerGroupHandler(config)
	//if err != nil {
	//	return nil, err
	//}

	go app.queueHandle()
	//handle := &consumerGroupHandler{
	//	AppId:        appid,
	//	ctx:          ctx,
	//	consumerFunc: app.queueHandle,
	//	log:          app.log,
	//}
	//go func() {
	//	defer func() {
	//		fmt.Println("exit consumer")
	//		err := app.queues.Close()
	//		if err != nil {
	//			handle.log.Errorf("%v app.queues.Close error %v", app, err)
	//		}
	//	}()
	//	for {
	//		err := app.queues.Consume(app.ctx, app.config.Topic, handle)
	//		if err != nil {
	//			handle.log.Errorf("%v group.Consume error %v", app, err)
	//		}
	//	}
	//}()
	//go func() {
	//	select {
	//	case err := <-app.queues.Errors():
	//		app.log.Errorf("app %s .Queues.Errors() %s", appid, err)
	//	case <-ctx.Done():
	//		return
	//	}
	//}()
	return app, nil
}

func (a *App) GetBucket(userId userId) *Bucket {
	return a.Buckets[a.GetBucketIndex(uint64(userId))]
}

type consumerGroupHandler struct {
	AppId        string
	ctx          context.Context
	consumerFunc func(msg *sarama.ConsumerMessage) error
	log          *zap.SugaredLogger
}

func (consumer consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (consumer consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (consumer consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	//todo 改成消费->后聚合排序comit最后一个消息
	fmt.Println(<-claim.Messages())
	for msg := range claim.Messages() {
		sess.MarkMessage(msg, "")
		//手动提交offset 消费完后聚合comit
		err := consumer.consumerFunc(msg)
		if err != nil {
			consumer.log.Errorf("consumerFunc error:%s", err.Error())
			return nil
		}

	}
	return nil
}

func (a *App) queueHandle() (err error) {
	//cometMsg := &comet.Msg{}
	//err := proto.Unmarshal(msg.Value, cometMsg)
	//if err != nil {
	//	a.log.Errorf("proto.Unmarshal error:%s", err.Error())
	//	return err
	//}
	//a.log.Debugf("%+v", cometMsg)
	for {
		select {
		case <-a.ctx.Done():
			//todo,得判断队列是否还有消息,有消息消费完
			qLen := GameQueue.Len(a.Appid)
			for i := 0; i >= qLen; i++ {
				a.log.Errorf("close app error,queue len:%d", qLen)
				a.PopQueueHandle()
			}
			return nil
		default:
			a.PopQueueHandle()
		}
		a.log.Debugf("queue qLen is %d", GameQueue.Len(a.Appid))
	}

	return nil
}
func (a *App) PopQueueHandle() {
	cometMsg, err := GameQueue.Pop(a.Appid)
	if err != nil {
		a.log.Errorf("pop message error msg:%v,err:%v", cometMsg, err)
	}
	switch cometMsg.Type {
	case comet.Type_PUSH:
		bucket := a.GetBucket(userId(cometMsg.ToId))
		if user, ok := bucket.users[userId(cometMsg.ToId)]; ok {
			err = user.Push(cometMsg)
			if err != nil {
				a.log.Errorf("user.Push error:%s", err.Error())
			}
		} else {
			a.log.Errorf("user not exist:%d", cometMsg.ToId)
		}
	case comet.Type_ROOM, comet.Type_AREA:
		a.broadcast(cometMsg)
	case comet.Type_CLOSE:
		//todo 关闭连接
		//if user, ok := bucket.users[cometMsg.ToId]; ok {
		//	_ = user.Close()
		//}
	default:
		a.log.Errorf("unknown msg type:%d", cometMsg.Type)
	}
}

// 广播工会消息
func (a *App) broadcast(c *comet.Msg) {
	a.log.Debugf("broadcast:%+v", c)
	//bucket 不涉及动态扩容 不需加锁
	for _, v := range a.Buckets {
		v.broadcast(c)
	}
}
