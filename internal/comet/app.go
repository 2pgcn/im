package comet

import (
	"context"
	"github.com/2pgcn/gameim/api/protocol"
	"github.com/2pgcn/gameim/conf"
	"github.com/2pgcn/gameim/pkg/event"
	"github.com/2pgcn/gameim/pkg/gamelog"
	"github.com/2pgcn/gameim/pkg/safe"
	"strconv"
	"sync"
)

type App struct {
	ctx      context.Context
	Appid    string
	Buckets  []*Bucket //app bucket
	lock     sync.RWMutex
	log      gamelog.GameLog
	receiver event.Receiver
	gopool   *safe.GoPool
}

func (a *App) GetLog() gamelog.GameLog {
	//todo set log prefix
	return a.log
}

// todo 修改下hash
func (a *App) GetBucketIndex(userid userId) int {
	idx, _ := strconv.Atoi(string(userid))
	return idx % len(a.Buckets)
}

// NewApp todo 暂时写死app 需要改为从存储中获取 config改成app config
func NewApp(ctx context.Context, c *conf.AppConfig, receiver event.Receiver, l gamelog.GameLog, gopool *safe.GoPool) (*App, error) {
	app := &App{
		ctx:      ctx,
		Appid:    c.Appid,
		Buckets:  make([]*Bucket, c.BucketNum),
		log:      l.AppendPrefix("app"),
		receiver: receiver,
		gopool:   gopool,
	}

	for i := 0; i < int(c.BucketNum); i++ {
		app.Buckets[i] = NewBucket(ctx, l)
	}
	return app, nil
}

func (a *App) Start() error {
	a.gopool.GoCtx(func(ctx context.Context) {
		err := a.queueHandle()
		if err != nil {
			a.GetLog().Errorf("app start queueHandle error:%s", err)
		}
	})
	return nil
}

func (a *App) GetBucket(userId userId) *Bucket {
	return a.Buckets[a.GetBucketIndex(userId)]
}

func (a *App) queueHandle() (err error) {
	receiverQueues, err := a.receiver.Receive(a.ctx)
	if err != nil {
		a.GetLog().Errorf("queueHandle error:%s", err)
		return err
	}
	for _, q1 := range receiverQueues {
		q := q1
		a.gopool.GoCtx(func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				case m := <-q:
					msg := m.GetQueueMsg()
					switch msg.Data.Type {
					case protocol.Type_PUSH:
						bucket := a.GetBucket(userId(msg.Data.GetToId()))
						if user, ok := bucket.users[userId(msg.Data.GetToId())]; ok {
							err = user.Push(a.ctx, m)
							if err != nil {
								a.GetLog().Errorf("user.Push error:%s", err.Error())
							}
						} else {
							a.GetLog().Debugf("user not exist:%d", msg.Data.GetToId())
						}
					case protocol.Type_ROOM, protocol.Type_APP:
						a.broadcast(m)

					}
				}
			}
		})
	}
	return nil
}

// 广播工会消息
func (a *App) broadcast(c event.Event) {
	//bucket 不涉及动态扩容 不需加锁
	for _, v := range a.Buckets {
		v.broadcast(c)
	}

}

func (a *App) Close() {
	msg := event.GetQueueMsg()
	msg.Data.Type = protocol.Type_CLOSE
	a.broadcast(msg)
	a.GetLog().Debug("app broadcast close success")
	for _, v := range a.Buckets {
		v.Close()
	}
	err := a.receiver.Close()
	if err != nil {
		a.GetLog().Errorf("%s:", err.Error())
	}
	a.GetLog().Debug("app receiver close success")

}
