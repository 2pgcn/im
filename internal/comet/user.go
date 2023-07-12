package comet

import (
	"bufio"
	"github.com/2pgcn/gameim/api/comet"
	"github.com/2pgcn/gameim/pkg/errors"
	"sync"
)

type User struct {
	Uid      userId
	Area     *Area
	Room     *Room
	Next     *User
	Prev     *User
	AppId    uint64
	AreaId   areaId
	RoomId   roomId
	lock     sync.RWMutex
	msgQueue chan *comet.Msg
	ReadBuf  *bufio.Reader
	WriteBuf *bufio.Writer
}

func NewUser() *User {
	return &User{
		msgQueue: make(chan *comet.Msg, 128),
	}
}

func (u *User) Push(m *comet.Msg) (err error) {
	select {
	case u.msgQueue <- m:
	default:
		err = errors.ErrMsgQueue
	}
	return
}

// Ready check the channel ready or close?
func (u *User) Ready() *comet.Msg {
	return <-u.msgQueue
}

// Close the channel.
func (u *User) Close() {
	u.msgQueue <- &comet.Msg{
		Type: comet.Type_CLOSE,
	}
	close(u.msgQueue)
}
