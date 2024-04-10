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
)

type User struct {
	ctx      context.Context
	log      gamelog.GameLog
	Uid      userId
	Room     *Room
	Next     *User
	Prev     *User
	AppId    string
	RoomId   roomId
	lock     sync.RWMutex
	msgQueue *event.Channel
	ReadBuf  *bufio.Reader
	WriteBuf *bufio.Writer
	conn     *net.TCPConn
	pool     *safe.GoPool
}

func NewUser(ctx context.Context, conn *net.TCPConn, log gamelog.GameLog) *User {
	pool := safe.NewGoPool(ctx, "gameim-comet-user")
	return &User{
		ctx:      ctx,
		msgQueue: event.NewChannel(128),
		conn:     conn,
		log:      log.AppendPrefix("user"),
		pool:     pool,
	}
}

func (u *User) Push(ctx context.Context, m event.Event) (err error) {
	return u.msgQueue.Send(ctx, m)
}

func (u *User) Pop(ctx context.Context) (res []chan event.Event) {
	var err error
	res, err = u.msgQueue.Receive(ctx)
	if err != nil {
		u.log.Errorf("pop msg error:%s", err)
		return
	}
	return

}

func (u *User) Start() {
	chans := u.Pop(u.ctx)
	for _, v1 := range chans {
		v := v1
		//todo,accept fin to exit
		u.pool.GoCtx(func(ctx context.Context) {
			for {
				select {
				case <-u.ctx.Done():
					return
				case msgEvent := <-v:
					gamelog.GetGlobalog().Debugf("user recv msg:%v", msgEvent)
					writeProto, err := msgEvent.ToProtocol()
					if err != nil {
						u.log.Errorf("writeProto err: %+v", writeProto)
						continue
					}
					if err = writeProto.WriteTcp(u.WriteBuf); err != nil {
						u.log.Errorf("writeProto.EncodeTo(user.WriteBuf) error(%v)", err)
						continue
					}
					if msgEvent.GetQueueMsg().Data.Type == protocol.Type_CLOSE {
						//close
						event.PutQueueMsg(msgEvent.GetQueueMsg())
						u.log.Debugf("recv msg close:%v", msgEvent.GetQueueMsg())
						return
					}
				}
			}
		})
	}
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
		u.log.Errorf("close error %s", err.Error())
	}
}

func (u *User) GetConn() net.Conn {
	return u.conn
}
