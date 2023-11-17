package service

import (
	"context"
	"github.com/2pgcn/gameim/api/comet"
	error2 "github.com/2pgcn/gameim/api/error"
	pb "github.com/2pgcn/gameim/api/logic"
	"github.com/2pgcn/gameim/conf"
	"github.com/2pgcn/gameim/internal/logic/data"
	"github.com/2pgcn/gameim/pkg/event"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
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
	uidInt, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		return nil, error2.AuthError
	}
	return &pb.AuthReply{Appid: "app001", Uid: uint64(uidInt), RoomId: strconv.Itoa(int(uidInt % 100))}, err
}

func (s *AuthService) OnConnect(ctx context.Context, req *pb.ConnectReq) (*emptypb.Empty, error) {
	return nil, nil
}

// todo 添加到kafka
func (s *AuthService) OnMessage(ctx context.Context, req *pb.MessageReq) (*emptypb.Empty, error) {
	ctx, span := otel.Tracer(conf.ServerName).Start(ctx, "OnMessageServices", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()
	err := s.user.WriteKafkaMessage(ctx, &event.Msg{
		H: map[string]any{},
		K: []byte(req.CometKey),
		Data: &comet.MsgData{
			Type:   req.Type,
			ToId:   req.ToId,
			SendId: req.SendId,
			Msg:    req.Msg,
		},
	})
	return &emptypb.Empty{}, err
}

func (s *AuthService) OnClose(context.Context, *pb.CloseReq) (*emptypb.Empty, error) {
	return nil, nil
}
