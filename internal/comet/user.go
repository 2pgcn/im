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
	msgQueue chan event.Event
	ReadBuf  *bufio.Reader
	WriteBuf *bufio.Writer
	conn     *net.TCPConn
}

func NewUser(ctx context.Context, conn *net.TCPConn, log gamelog.GameLog) *User {
	return &User{
		msgQueue: make(chan event.Event, 128),
		conn:     conn,
		log:      log,
	}
}

func (u *User) Push(m event.Event) (err error) {
	select {
	case u.msgQueue <- m:
	}
	return
}

// Ready check the channel ready or close?
func (u *User) Pop() event.Event {
	return <-u.msgQueue
}

// Close the channel.
func (u *User) Close() {
	u.msgQueue <- &event.Msg{
		Data: &comet.MsgData{
			Type: comet.Type_CLOSE,
		},
	}
	u.lock.Lock()
	defer u.lock.Unlock()
	err := u.conn.CloseRead()
	if err != nil {
		u.log.Errorf("close error %s", err.Error())
	}
}

func (u *User) GetConn() net.Conn {
	return u.conn
}
