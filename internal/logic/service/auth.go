package service

import (
	"context"
	pb "github.com/2pgcn/gameim/api/logic"
	"github.com/2pgcn/gameim/internal/logic/data"
	"github.com/2pgcn/gameim/pkg/event"
	"github.com/2pgcn/gameim/pkg/gamelog"
	"github.com/go-kratos/kratos/v2/log"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	"io"
	"math/rand"
	"strconv"
	"time"
)

type AuthService struct {
	pb.UnimplementedLogicServer
	log  log.Logger
	user *data.Data
}

func NewAuthService(log log.Logger, d *data.Data) *AuthService {
	return &AuthService{
		log:  log,
		user: d,
	}
}

func (s *AuthService) OnAuth(ctx context.Context, req *pb.AuthReq) (*pb.AuthReply, error) {
	//uid, err := s.user.GetIdByToken(ctx, req.Token)
	//if uid <= 0 {
	//	return nil, error2.AuthError
	//}
	uid := req.GetToken()
	r := rand.Intn(100)
	return &pb.AuthReply{Appid: "app001", Uid: uid, RoomId: strconv.Itoa(r % 100)}, nil
}

func (s *AuthService) OnConnect(ctx context.Context, req *pb.ConnectReq) (*emptypb.Empty, error) {
	return nil, nil
}

// todo 添加到kafka
func (s *AuthService) OnMessage(stream pb.Logic_OnMessageServer) error {
	//ctx, span := otel.Tracer(conf.ServerName).Start(ctx, "OnMessageServices", trace.WithSpanKind(trace.SpanKindInternal))
	//defer span.End()
	var msgsReply = pb.MessageReply{}
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&msgsReply)
		}
		if err != nil {
			//todo,设置logic错误,方便调用方处理
			return err
		}
		ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
		qMsg := event.GetQueueMsg()
		qMsg.Data.Type = r.Type
		qMsg.Data.ToId = r.ToId
		qMsg.Data.SendId = r.SendId
		qMsg.Data.Msg = r.Msg
		//todo,根据userid 获取appid,->sendto/topic+appid
		err = s.user.WriteMessage(ctx, qMsg)
		if err != nil {
			gamelog.GetGlobalog().Errorf("logic writeMessage error:%s", err)
			continue
		}
		msgsReply.Msgs = append(msgsReply.Msgs, &pb.MsgReply{
			SendId: r.SendId,
			MsgId:  r.MsgId,
		})
	}

}

func (s *AuthService) OnClose(context.Context, *pb.CloseReq) (*emptypb.Empty, error) {
	return nil, nil
}
