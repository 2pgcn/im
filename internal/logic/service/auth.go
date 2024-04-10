package service

import (
	"context"
	pb "github.com/2pgcn/gameim/api/logic"
	"github.com/2pgcn/gameim/internal/logic/data"
	"github.com/2pgcn/gameim/pkg/event"
	"github.com/go-kratos/kratos/v2/log"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	"math/rand"
	"strconv"
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
func (s *AuthService) OnMessage(ctx context.Context, req *pb.MessageReq) (*emptypb.Empty, error) {
	//ctx, span := otel.Tracer(conf.ServerName).Start(ctx, "OnMessageServices", trace.WithSpanKind(trace.SpanKindInternal))
	//defer span.End()
	data := event.GetQueueMsg()
	data.Data.Type = req.Type
	data.Data.ToId = req.ToId
	data.Data.SendId = req.SendId
	data.Data.Msg = req.Msg
	//todo,根据userid 获取appid,->sendto/topic+appid
	err := s.user.WriteMessage(ctx, data)
	return &emptypb.Empty{}, err
}

func (s *AuthService) OnClose(context.Context, *pb.CloseReq) (*emptypb.Empty, error) {
	return nil, nil
}
