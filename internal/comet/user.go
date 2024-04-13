package comet

import (
	"bufio"
	"context"
	"errors"
	"github.com/2pgcn/gameim/api/protocol"
	"github.com/2pgcn/gameim/pkg/event"
	"github.com/2pgcn/gameim/pkg/gamelog"
	"github.com/2pgcn/gameim/pkg/safe"
	"net"
	"sync"
	"sync/atomic"
)

type userId string

type User struct {
	ctx         context.Context
	log         gamelog.GameLog
	Uid         userId
	Room        *Room
	Next        *User
	Prev        *User
	AppId       string
	RoomId      roomId
	lock        sync.RWMutex
	ReadBuf     *bufio.Reader
	WriteBuf    *bufio.Writer
	conn        *net.TCPConn
	pool        *safe.GoPool
	msgQueueLen int64
	msgQueue    chan event.Event
	allRecvMsg  int64
}

func NewUser(ctx context.Context, conn *net.TCPConn, log gamelog.GameLog) *User {
	pool := safe.NewGoPool(ctx, "gameim-comet-user")
	return &User{
		ctx:      ctx,
		msgQueue: make(chan event.Event, 128),
		conn:     conn,
		log:      log.AppendPrefix("user"),
		pool:     pool,
	}
}

func (u *User) Push(ctx context.Context, m event.Event) (err error) {
	u.msgQueue <- m
	atomic.AddInt64(&u.msgQueueLen, 1)
	return nil
}

func (u *User) Pops(ctx context.Context) chan event.Event {
	return u.msgQueue

}

func (u *User) Start() {
	//todo,accept fin to exit
	u.pool.GoCtx(func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case msgEvent := <-u.Pops(ctx):
				msgEvents := []event.Event{msgEvent}
				l := atomic.LoadInt64(&u.msgQueueLen)
				for i := 1; i < int(l); i++ {
					msgEvents = append(msgEvents, <-u.Pops(ctx))
				}
				atomic.AddInt64(&u.msgQueueLen, int64(len(msgEvents))*-1)
				for _, v := range msgEvents {
					writeProto, err := v.ToProtocol()
					if err != nil {
						u.log.Info("writeProto err: %+v,%+v", writeProto, v)
						continue
					}
					if err = writeProto.WriteTcpNotFlush(u.WriteBuf); err != nil {
						u.log.Infof("writeProto.EncodeTo(user.WriteBuf) error(%v)", err)
						continue
					}
					if v.GetQueueMsg().Data.Type == protocol.Type_CLOSE {
						//close
						event.PutQueueMsg(v.GetQueueMsg())
						u.log.Infof("recv msg close:%v", v.GetQueueMsg())
						u.WriteBuf.Flush()
						return
					}
				}
				u.WriteBuf.Flush()

			}
		}
	})
}

// Close the channel.
func (u *User) Close() {
	msg := event.GetQueueMsg()
	msg.Data.Type = protocol.Type_CLOSE
	err := u.Push(context.Background(), msg)
	if err != nil {
		u.log.Errorf("close u.push error %s", err.Error())
	}
	u.lock.Lock()
	defer u.lock.Unlock()
	err = u.conn.CloseRead()
	if err != nil && !errors.Is(err, net.ErrClosed) {
		u.log.Debugf("close error %s", err.Error())
	}
}

func (u *User) GetConn() net.Conn {
	return u.conn
}
