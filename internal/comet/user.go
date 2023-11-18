package comet

import (
	"bufio"
	"context"
	"github.com/2pgcn/gameim/api/comet"
	"github.com/2pgcn/gameim/pkg/event"
	"github.com/2pgcn/gameim/pkg/gamelog"
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
}

func NewUser(ctx context.Context, conn *net.TCPConn, log gamelog.GameLog) *User {
	return &User{
		ctx:      ctx,
		msgQueue: event.NewChannel(128),
		conn:     conn,
		log:      log,
	}
}

func (u *User) Push(m event.Event) (err error) {
	return u.msgQueue.Send(context.Background(), m)
}

// Ready check the channel ready or close?
func (u *User) Pop() (event.Event, error) {
	return u.msgQueue.Receive(context.Background())
}

// Close the channel.
func (u *User) Close() {
	err := u.Push(&event.Msg{
		Data: &comet.MsgData{
			Type: comet.Type_CLOSE,
		}})

	if err != nil {
		u.log.Errorf("close u.push error %s", err.Error())
	}
	u.lock.Lock()
	defer u.lock.Unlock()
	err = u.conn.CloseRead()
	if err != nil {
		u.log.Errorf("close error %s", err.Error())
	}
}

func (u *User) GetConn() net.Conn {
	return u.conn
}
