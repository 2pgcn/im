package comet

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/php403/gameim/api/comet"
	"github.com/php403/gameim/config"
	"go.uber.org/zap"
	"sync"
	"time"
)

type App struct {
	ctx       context.Context
	Appid     string
	config    *config.CometConfigQueueMsg
	Buckets   []*Bucket //app bucket
	bucketIdx uint32
	lock      sync.RWMutex
	queues    sarama.ConsumerGroup
	log       *zap.SugaredLogger
}

func (a *App) GetBucketIndex(userid uint64) int {
	return int(userid % uint64(a.bucketIdx))
}

// NewApp todo 暂时写死app 需要改为从存储中获取
func NewApp(ctx context.Context, appid string) (*App, error) {
	app := &App{
		ctx:       ctx,
		Appid:     appid,
		Buckets:   make([]*Bucket, 1024),
		bucketIdx: 1024, //todo 改为配置
	}
	for i := 0; i < 1024; i++ {
		app.Buckets[i] = NewBucket()
	}
	consumer, err := NewConsumerGroupHandler(app.config)
	if err != nil {
		return nil, err
	}
	app.queues = consumer
	handle := &consumerGroupHandler{
		AppId:        appid,
		ctx:          ctx,
		consumerFunc: app.queueHandle,
	}
	go func() {
		defer func() {
			err := app.queues.Close()
			if err != nil {

			}
		}()
		select {
		case <-ctx.Done():
			app.log.Info("app context done")
			return
		default:
			for {
				err := app.queues.Consume(app.ctx, app.config.Topic, handle)
				if err != nil {
					handle.log.Errorf("%v group.Consume error %v", app, err)
				}
				time.Sleep(time.Second * 2)
			}
		}
	}()
	go func(ctx context.Context) {
		select {
		case err := <-app.queues.Errors():
			app.log.Errorf("app %s .Queues.Errors() %s", appid, err.Error())
		case <-ctx.Done():
			return
		}
	}(ctx)
	return app, nil
}

func (a *App) GetBucket(userId uint64) *Bucket {
	return a.Buckets[a.GetBucketIndex(userId)]
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
	for msg := range claim.Messages() {
		//手动提交offset 消费完后聚合comit
		err := consumer.consumerFunc(msg)
		if err != nil {
			consumer.log.Errorf("consumerFunc error:%s", err.Error())
			return err
		}
	}
	return nil
}

func (a *App) queueHandle(msg *sarama.ConsumerMessage) error {
	cometMsg := &comet.Msg{}
	err := proto.Unmarshal(msg.Value, cometMsg)
	if err != nil {
		a.log.Errorf("proto.Unmarshal error:%s", err.Error())
		return err
	}
	bucket := a.GetBucket(cometMsg.ToId)
	switch cometMsg.Type {
	case comet.Type_PUSH:
		if user, ok := bucket.users[cometMsg.ToId]; ok {
			return user.Push(cometMsg)
		}
	case comet.Type_ROOM:
		if room := bucket.Room(cometMsg.ToId); room != nil {
			return room.Push(cometMsg)
		}
	case comet.Type_AREA:
		if area := bucket.Area(cometMsg.ToId); area != nil {
			return area.Push(cometMsg)
		}
	case comet.Type_CLOSE:
		//todo 关闭连接
		//if user, ok := bucket.users[cometMsg.ToId]; ok {
		//	_ = user.Close()
		//}
	}
	return nil
}
