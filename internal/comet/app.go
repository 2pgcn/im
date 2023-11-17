package comet

import (
	"context"
	"github.com/2pgcn/gameim/api/comet"
	"github.com/2pgcn/gameim/conf"
	"github.com/2pgcn/gameim/pkg/event"
	"github.com/2pgcn/gameim/pkg/gamelog"
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
}

func (a *App) GetLog() gamelog.GameLog {
	//todo set log prefix
	return a.log
}

func (a *App) GetBucketIndex(userid userId) int {
	return int(uint64(userid) % uint64(len(a.Buckets)))
}

// NewApp todo 暂时写死app 需要改为从存储中获取 config改成app config
func NewApp(ctx context.Context, c *conf.AppConfig, receiver event.Receiver, l gamelog.GameLog) (*App, error) {
	app := &App{
		ctx:      ctx,
		Appid:    c.Appid,
		Buckets:  make([]*Bucket, c.BucketNum),
		log:      l,
		receiver: receiver,
	}

	for i := 0; i < int(c.BucketNum); i++ {
		app.Buckets[i] = NewBucket(ctx, l)
	}
	go func() {
		err := app.queueHandle()
		if err != nil {

		}
	}()
	return app, nil
}

func (a *App) GetBucket(userId userId) *Bucket {
	return a.Buckets[a.GetBucketIndex(userId)]
}

func (a *App) queueHandle() (err error) {
	for {
		select {
		case <-a.ctx.Done():
			return a.receiver.Close()
		default:
			msg, err := a.receiver.Receive(a.ctx)
			a.GetLog().Debug(msg)
			if err != nil {
				a.GetLog().Errorf("message error:%s", err)
				continue
			}
			msgData, ok := msg.RawValue().(*comet.MsgData)
			if !ok {
				continue
			}
			switch msgData.GetType() {
			case comet.Type_PUSH:
				uidInt, err := strconv.ParseUint(msgData.ToId, 10, 64)
				if err != nil {
					a.GetLog().Errorf("uidInt parse error:%s", err)
				}
				bucket := a.GetBucket(userId(uidInt))
				if user, ok := bucket.users[userId(uidInt)]; ok {
					err = user.Push(msg)
					if err != nil {
						a.GetLog().Errorf("user.Push error:%s", err.Error())
					}
				} else {
					a.GetLog().Errorf("user not exist:%d", msgData.ToId)
				}
			case comet.Type_ROOM, comet.Type_APP:
				a.broadcast(msg)
			// todo add func
			//case comet.Type_CLOSE:
			//	if user, ok := bucket.users[cometMsg.ToId]; ok {
			//		_ = user.Close()
			//	}
			default:
				a.GetLog().Errorf("unknown msg type:%d", msgData.GetType())
			}
			if err != nil {
				a.GetLog().Errorf("queueHandle %s", err)
			}
		}
	}
}

// 广播工会消息
func (a *App) broadcast(c event.Event) {
	a.GetLog().Debugf("broadcast:%+v", c)
	//bucket 不涉及动态扩容 不需加锁
	for _, v := range a.Buckets {
		v.broadcast(c)
	}
}

func (a *App) Close() {
	a.broadcast(&event.Msg{
		Data: &comet.MsgData{
			Type:   comet.Type_CLOSE,
			ToId:   "",
			SendId: 0,
			Msg:    []byte("server close"),
		},
	})
	for _, v := range a.Buckets {
		v.Close()
	}
	err := a.receiver.Close()
	if err != nil {
		a.GetLog().Errorf("%s:", err.Error())
	}

}
