package data

import (
	"context"
	"github.com/2pgcn/gameim/conf"
	"github.com/2pgcn/gameim/pkg/event"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData)

var defaultMsgSize = 10240

// Data .
type Data struct {
	redisClient *redis.Client
	mysqlClient *gorm.DB
	//nsqClient   event.Sender
	producer event.Sender
	log      log.Logger
	//解决queue client写入chan满了后会无限阻塞造成程序hang死问题
	msgs chan event.Event
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	//redisClient, err := event.NewRedisReceiver(c.RedisQueue)
	//if err != nil {
	//	panic(err)
	//}
	//mysqlClient, err := NewMysql(c.Mysql)
	//if err != nil {
	//	panic(err)
	//}
	//kafkaClient, err := event.NewKafkaSender(c.GetKafka())
	//if err != nil {
	//	panic(err)
	//}
	//nsqClient, err := event.NewNsqSender(c.GetNsq())
	//if err != nil {
	//	panic(err)
	//}
	sockClient, err := event.NewSockRender(c.Queue.Sock.Address)
	if err != nil {
		panic(err)
	}
	cleanup := func() {
		//if err := redisClient.Close(); err != nil {
		//	log.NewHelper(logger).Errorf("closing the redis client err:%v", err)
		//}
		//sqlDb, err := mysqlClient.DB()
		//if err != nil {
		//	log.NewHelper(logger).Errorf("closing the mysql client err:%v", err)
		//}
		//if err = sqlDb.Close(); err != nil {
		//	log.NewHelper(logger).Errorf("closing the sqlDb client err:%v", err)
		//}

		//log.NewHelper(logger).Infof("start cleanup")
		//if err = nsqClient.Close(); err != nil {
		//	log.NewHelper(logger).Errorf("closing the nsq client err:%v", err)
		//}
		if err = sockClient.Close(); err != nil {
			log.NewHelper(logger).Errorf("closing the nsq client err:%v", err)
		}
	}
	return &Data{producer: sockClient}, cleanup, nil
}

func NewRedis(ctx context.Context, data *conf.RedisQueue) (client *redis.Client, err error) {
	client = redis.NewClient(&redis.Options{
		Addr:         data.Rdscfg.Addr,
		Password:     data.Rdscfg.GetPasswd(), // no password set
		DB:           int(data.Rdscfg.UseDb),  // use default DB
		PoolSize:     int(data.Rdscfg.PoolSize),
		ReadTimeout:  data.Rdscfg.ReadTimeout.AsDuration(),
		WriteTimeout: data.Rdscfg.WriteTimeout.AsDuration(),
	})
	err = client.Ping(ctx).Err()
	return
}
