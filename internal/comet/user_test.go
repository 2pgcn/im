package comet

import (
	"context"
	"github.com/2pgcn/gameim/api/comet"
	"github.com/2pgcn/gameim/pkg/event"
	"go.uber.org/zap"
	"net"
	"reflect"
	"testing"
)

var tLog, _ = zap.NewProduction()

func newListener(t testing.TB, network string) net.Listener {
	var lc *net.ListenConfig
	ln, err := lc.Listen(context.Background(), network, "127.0.0.1:0")
	if err != nil {
		t.Helper()
		t.Fatal(err)
	}
	return ln
}

func newAccept(t testing.TB, listener *net.TCPListener) *net.TCPConn {
	conn, err := listener.AcceptTCP()
	if err != nil {
		t.Errorf("accept error")
	}
	return conn
}

func newUser(t testing.TB) *User {
	listener := newListener(t, "tcp")
	l, ok := listener.(*net.TCPListener)
	if !ok {
		t.Errorf("newUser want to TCPListener not %s", reflect.TypeOf(listener))
	}
	accept := newAccept(t, l)
	return NewUser(context.Background(), accept, tLog.Sugar())
}

func newUserNotListen(t testing.TB) *User {
	return NewUser(context.Background(), nil, tLog.Sugar())
}

func TestNewUserSendAndRecv(t *testing.T) {
	//user := newUserNotListen(t)
	//msg := &comet.MsgData{
	//	Type:   comet.Type_CLOSE,
	//	ToId:   0,
	//	SendId: 0,
	//	Msg:    []byte("test"),
	//}
	//err := user.Push(msg)
	//if err != nil {
	//	t.Error(err)
	//}
	//newMsg := user.Pop()
	//if newMsg.GetType() != msg.GetType() && string(newMsg.GetMsg()) != string(newMsg.GetMsg()) {
	//	t.Errorf("TestNewUser new user error:msg error")
	//}
}

func TestUserClose(t *testing.T) {
	user := newUser(t)
	user.Close()
	msgE := user.Pop()
	msg := msgE.(*event.Msg)
	if msg.GetData().Type != comet.Type_CLOSE {
		t.Errorf("TestUserClose user msg type error")
	}
	var b [1]byte
	n, err := user.GetConn().Read(b[:])
	if n != 0 || err == nil {
		t.Fatalf("TestUserClose read got (%d, %v); want (0, error)", n, err)
	}
}
