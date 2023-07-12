package comet

import (
	"context"
	"github.com/2pgcn/gameim/api/comet"
	"github.com/2pgcn/gameim/config"
	"github.com/go-redis/redis/v8"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
)

var GameQueue Queue

// todo 代完善,目前仅为测试
type Queue interface {
	Pop(topic string) (*comet.Msg, error)
	Push(topic string, msg *comet.Msg) error
	Close() error
	Len(topic string) int
}

type RedisQueue struct {
	appid  string
	ctx    context.Context
	client *redis.Client
	log    *zap.SugaredLogger
}

func (c *RedisQueue) Pop(topic string) (*comet.Msg, error) {
	result, err := c.client.BRPop(c.ctx, 0, topic).Result()
	if err != nil {
		return nil, err
	}
	msg := &comet.Msg{}
	err = proto.Unmarshal([]byte(result[1]), msg)
	return msg, err
}

func (c *RedisQueue) Push(topic string, msg *comet.Msg) error {
	msgS, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	c.client.LPush(c.ctx, topic, msgS)
	return nil
}

func (c *RedisQueue) Len(topic string) int {
	return int(c.client.LLen(c.ctx, topic).Val())
}

func (c *RedisQueue) Close() error {
	return c.client.Close()
}

func InitChanQueue(ctx context.Context, log *zap.SugaredLogger, c *config.CometConfigQueueMsg) Queue {
	client := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       int(c.Redis.Db),
	})
	GameQueue = &RedisQueue{
		ctx:    ctx,
		client: client,
		log:    log,
	}
	res := client.Ping(ctx)
	if res.Err() != nil {
		panic(res.Err())
	}
	return GameQueue
}

//// NewConsumerGroupHandler NewKafkaConsumer todo  NewKafkaConsumer改成接口 后续扩展
//func NewConsumerGroupHandler(c *config.CometConfigQueueMsg) (consumer sarama.ConsumerGroup, err error) {
//	kconfig := sarama.NewConfig()
//	kconfig.Version = sarama.V2_7_1_0 //
//	//kconfig.Consumer.Return.Errors = true
//	kconfig.Consumer.Offsets.AutoCommit.Enable = true              // 自动提交
//	kconfig.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second // 间隔
//	kconfig.Consumer.Offsets.Initial = sarama.OffsetNewest
//	kconfig.Consumer.Offsets.Retry.Max = 3
//	kconfig.Consumer.MaxWaitTime = 1 * time.Second
//	return sarama.NewConsumerGroup(c.Brokers, c.Group, kconfig)
//}
