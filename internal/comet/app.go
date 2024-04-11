package comet

import (
	"bufio"
	"context"
	"github.com/2pgcn/gameim/api/gerr"
	"github.com/2pgcn/gameim/api/logic"
	"github.com/2pgcn/gameim/api/protocol"
	"github.com/2pgcn/gameim/conf"
	"github.com/2pgcn/gameim/pkg/event"
	"github.com/2pgcn/gameim/pkg/gamelog"
	"github.com/2pgcn/gameim/pkg/safe"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const defaultLogicLens = 12
const defaultMessageReqLens = 4096

type App struct {
	ctx      context.Context
	Appid    string
	conf     *conf.AppConfig
	lock     sync.RWMutex
	log      gamelog.GameLog
	receiver event.Receiver
	gopool   *safe.GoPool
	Buckets  []*Bucket //app bucket
	//len(logicMsgs)==len(logicClients)=defaultLogiLens
	lenNums      int
	logicMsgs    []*logicMsg
	logicClients []logic.LogicClient
}

type logicMsg struct {
	msgReq chan *logic.MessageReq
	len    int64
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
		conf:     c,
		Appid:    c.Appid,
		Buckets:  make([]*Bucket, c.BucketNum),
		log:      l.AppendPrefix("app"),
		receiver: receiver,
		gopool:   gopool,
	}

	for i := 0; i < int(c.BucketNum); i++ {
		app.Buckets[i] = NewBucket(ctx, l)
	}
	for i := 0; i < defaultLogicLens; i++ {
		app.logicMsgs = append(app.logicMsgs, &logicMsg{
			msgReq: make(chan *logic.MessageReq, defaultMessageReqLens),
			len:    0,
		})
		app.logicClients = append(app.logicClients, newLogicClient(ctx, c.LogicClientGrpc))
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
	//读写logic msg
	for i := 0; i < defaultLogicLens; i++ {
		index := i
		a.gopool.GoCtx(func(ctx context.Context) {
			//等待grpc自动重连
			for {
				select {
				case <-ctx.Done():
					return
					//一次发送1/3
				case req := <-a.logicMsgs[index].msgReq:
					stream, err := a.logicClients[index].OnMessage(ctx)
					if err != nil {
						a.logicMsgs[index].msgReq <- req
						a.GetLog().Errorf("stream :OnMessage init is error:%s", err)
						err = nil
						time.Sleep(time.Second * 1)
					}
					atomic.AddInt64(&a.logicMsgs[index].len, -1)
					num := min((atomic.LoadInt64(&a.logicMsgs[index].len)), int64(defaultMessageReqLens/8))
					//todo 与for合并代码
					err = stream.Send(req)
					if err != nil {
						a.GetLog().Errorf("ListStr get stream err: %v", err)
						continue
					}
					for i := 0; i < int(num); i++ {
						err = stream.Send(req)
						if err != nil {
							a.GetLog().Errorf("ListStr get stream err: %v", err)
							continue
						}
						//失败暂不重新入队,成功则会收到ack并发送给客户端
					}
					resp, err := stream.CloseAndRecv()
					for _, v := range resp.Msgs {
						msg := event.GetQueueMsg()
						msg.Data.Type = protocol.Type_ACK
						msg.SetId(v.SendId)
						user, ok := a.GetBucket(userId(v.SendId)).users[userId(v.SendId)]
						if !ok {
							continue
						}
						user.Push(ctx, msg)
					}
					time.Sleep(time.Millisecond * 50)
				}
			}
		})
	}

	return nil
}

func (a *App) AddUser(token string, conn *net.TCPConn, br *bufio.Reader, bw *bufio.Writer) {
	a.gopool.GoCtx(func(ctx context.Context) {
		p := protocol.ProtoPool.Get()
		defer func() {
			protocol.ProtoPool.Put(p)
		}()
		//todo 配置超时时间
		grpcCtx, _ := context.WithTimeout(ctx, a.conf.LogicClientGrpc.Timeout.AsDuration())
		authReply, err := a.logicClients[rand.Int63n(int64(a.conf.LogicClientGrpc.ClientNum))].OnAuth(grpcCtx, &logic.AuthReq{
			Token: token,
		})
		gamelog.GetGlobalog().Info(authReply)

		if err != nil {
			gamelog.GetGlobalog().Errorf("req msg error%s", err)
			return
		}
		user := NewUser(ctx, conn, a.log)
		user.ReadBuf = br
		user.WriteBuf = bw
		user.Uid = userId(authReply.Uid)
		user.RoomId = roomId(authReply.RoomId)
		bucket := a.GetBucket(user.Uid)
		bucket.PutUser(user)
		defer bucket.DeleteUser(user.Uid)
		p.Reset()
		p.Op = protocol.OpAuthReply
		if err = p.WriteTcp(user.WriteBuf); err != nil {
			a.GetLog().Debugf("write proto err:%s", err)
			return
		}
		//启动用户获取自己msg
		a.gopool.GoCtx(func(ctx context.Context) {
			user.Start()
		})
		//不断读消息
		var sendType protocol.Type
		for {
			err = p.DecodeFromBytes(user.ReadBuf)
			if err != nil {
				gamelog.GetGlobalog().Debugf("client DecodeFromBytes error,%s", err)
				return
			}
			//msgCtx, span := trace_conf.SetTrace(context.Background(), trace_conf.COMET_RECV_CIENT_MSG,
			//	trace.WithSpanKind(trace.SpanKindInternal), trace.WithAttributes(
			//		attribute.Int("type", int(p.Op)),
			//		attribute.String("data", string(p.Data)),
			//	))
			//span.End()
			msgP := event.GetMsg()
			err = msgP.Unmarshal(p.Data)
			if err != nil {
				replyErr := gerr.ErrorMsgFormatError("p.DecodeFromBytes error")
				p.SetErrReply(replyErr)
				_ = p.WriteTcp(user.WriteBuf)
				return
			}
			gamelog.GetGlobalog().Debugf("recv client msg:%v,:srcMsg:%v", msgP, p)
			switch p.Op {
			case protocol.OpHeartbeat:
				// 更新最小堆
				bucket.lock.RLock()
				bucket.UpHeartTime(user.Uid, time.Now().Add(a.conf.KeepaliveTimeout.AsDuration()))
				bucket.lock.Unlock()
				continue
			case protocol.OpDisconnect:
				//心跳堆无需删除,在取出的时候兼容
				bucket.DeleteUser(user.Uid)
				return
			case protocol.OpSendAreaMsg, protocol.OpSendRoomMsg, protocol.OpSendMsg:
				if p.Op == protocol.OpSendAreaMsg {
					sendType = protocol.Type_APP
				}
				if p.Op == protocol.OpSendRoomMsg {
					sendType = protocol.Type_ROOM
				}
				if p.Op == protocol.OpSendMsg {
					sendType = protocol.Type_PUSH
				}
				a.logicMsgs[a.GetBucketIndex(user.Uid)].msgReq <- &logic.MessageReq{
					Type:     sendType,
					SendId:   string(user.Uid),
					MsgId:    strconv.Itoa(int(p.Seq)),
					ToId:     msgP.ToId,
					Msg:      p.Data,
					CometKey: a.Appid,
				}
			}
		}
	})
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
	for i := 0; i < len(a.logicClients); i++ {
		close(a.logicMsgs[i].msgReq)
	}
	a.GetLog().Debug("app receiver close success")

}
