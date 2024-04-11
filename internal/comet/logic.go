package comet

import (
	"context"
	"github.com/2pgcn/gameim/api/logic"
	"github.com/2pgcn/gameim/conf"
	"github.com/2pgcn/gameim/pkg/event"
	"github.com/2pgcn/gameim/pkg/gamelog"
	"strconv"
)

type LogicInterface interface {
	//可能单条可批量发送
	OnMessage(ctx context.Context, e []event.Event) (failEvent []event.Event, err error)
	OnClose(ctx context.Context, uid userId) error
	OnAuth(ctx context.Context, token string) (reply *logic.AuthReply, err error)
	OnConnect(ctx context.Context, uid userId) error
}

type LogicTest struct {
	ctx      context.Context
	producer event.Sender
}

func NewLogicClientTest(ctx context.Context, c *conf.Sock) (LogicInterface, error) {
	s, err := event.NewSockRender(c.Address)
	return &LogicTest{
		ctx:      ctx,
		producer: s,
	}, err
}

func (t *LogicTest) OnMessage(ctx context.Context, e []event.Event) (failEvent []event.Event, err error) {
	gamelog.GetGlobalog().Debug(len(e))
	for _, v := range e {
		err = t.producer.Send(ctx, v)
		failEvent = append(failEvent, v)
		return
	}
	return
}

func (t *LogicTest) OnClose(ctx context.Context, uid userId) error {
	msg := event.GetCloseMsg()
	return t.producer.Send(ctx, msg)
}

func (t *LogicTest) OnAuth(ctx context.Context, token string) (reply *logic.AuthReply, err error) {
	uidInt, _ := strconv.Atoi(token)
	return &logic.AuthReply{
		Uid:    token,
		Appid:  "app001",
		RoomId: strconv.Itoa(uidInt / 20),
	}, nil
}

func (t *LogicTest) OnConnect(ctx context.Context, uid userId) error {
	return nil
}
