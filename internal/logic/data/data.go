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

// Data .
type Data struct {
	redisClient *redis.Client
	mysqlClient *gorm.DB
	nsqClient   event.Sender
	log         log.Logger
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	//redisClient, err := NewRedis(context.Background(), c.Redis)
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
	nsqClient, err := event.NewNsqSender(c.GetNsq())
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

		if err = nsqClient.Close(); err != nil {
			log.NewHelper(logger).Errorf("closing the nsq client err:%v", err)
		}
	}
	return &Data{nsqClient: nsqClient}, cleanup, nil
}

func NewRedis(ctx context.Context, data *conf.Data_Redis) (client *redis.Client, err error) {
	client = redis.NewClient(&redis.Options{
		Addr:         data.Addr,
		Password:     "",              // no password set
		DB:           int(data.UseDb), // use default DB
		PoolSize:     int(data.PoolSize),
		ReadTimeout:  data.ReadTimeout.AsDuration(),
		WriteTimeout: data.WriteTimeout.AsDuration(),
	})
	client.Ping(ctx)
	return
}
