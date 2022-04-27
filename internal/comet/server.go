package comet

import (
	"bufio"
	"context"
	"github.com/php403/gameim/api/comet"
	"github.com/php403/gameim/api/logic"
	"github.com/php403/gameim/api/protocol"
	"github.com/php403/gameim/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"net"
	"sync"
	"time"
)

const appid = "app01"

type Server struct {
	ctx   context.Context
	conf  config.CometConfig
	Name  string
	log   *zap.SugaredLogger
	Apps  map[string]*App
	logic logic.LogicClient
	Lock  sync.RWMutex
}

func NewServer(ctx context.Context, log *zap.SugaredLogger, config *config.CometConfig) {
	server := &Server{
		ctx:   ctx,
		Name:  config.Server.Name,
		log:   log,
		Apps:  map[string]*App{},
		logic: newLogicClient(ctx, config.LogicClientGrpc),
	}
	//启动消费队列消费
	for _, v := range config.Server.Addrs {
		addr := v
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					server.bindConn(ctx, addr)
				}
			}
		}()
	}

}

func (s *Server) bindConn(ctx context.Context, host string) {
	var addr *net.TCPAddr
	var listener *net.TCPListener
	var conn net.Conn
	var err error
	if addr, err = net.ResolveTCPAddr("tcp", host); err != nil {
		s.log.Errorf("net.ResolveTCPAddr(tcp, %s) error(%v)", host, err)
		return
	}

	if listener, err = net.ListenTCP("tcp", addr); err != nil {
		s.log.Errorf("net.ListenTCP(tcp, %s) error(%v)", addr, err)
		return
	}

	if conn, err = listener.AcceptTCP(); err != nil {
		s.log.Errorf("accept comet error %v", zap.Error(err))
		return
	}
	go s.handleComet(ctx, conn)
	return
}

func (s *Server) handleComet(ctx context.Context, conn net.Conn) {

	user := NewUser()
	user.ReadBuf = bufio.NewReader(conn)
	user.WriteBuf = bufio.NewWriter(conn)
	//todo 用户登录后加入最小堆(nginx采用红黑树event)
	//获取用户消息
	proto := &protocol.Proto{}
	if err := proto.DecodeFromBytes(user.ReadBuf); err != nil {
		s.log.Errorf("protocol.DecodeFromBytes(user.ReadBuf, &user.Msg) error(%v)", err)
		return
	}
	if proto.Op != protocol.OpAuth {
		//todo return client err msg
		err := conn.Close()
		if err != nil {
			s.log.Errorf("conn.Close() error(%v)", err)
		}
		s.log.Errorf("protocol.Op != protocol.OpLogin")
		return
	}
	grpcCtx, cancal := context.WithTimeout(s.ctx, time.Duration(s.conf.LogicClientGrpc.Timeout)*time.Second)
	defer cancal()
	authReply, err := s.logic.OnAuth(grpcCtx, &logic.AuthReq{
		Token: string(proto.Data),
	})
	if err != nil {
		err = conn.Close()
		if err != nil {
			s.log.Errorf("conn.Close() error(%v)", err)
			return
		}
		s.log.Errorf("conn.Close() error(%v)", err)
		return
	}
	//appid 先写死,后面通过proto.data里解码获取
	//userInfo authReply
	user.Uid = authReply.Uid
	user.AreaId = authReply.AreaId
	user.RoomId = authReply.RoomId
	var app *App
	var ok bool
	s.Lock.RLock()
	if app, ok = s.Apps[appid]; !ok {
		app, err = NewApp(ctx, appid)
		if err != nil {
			return
		}
		s.Apps[appid] = app
	}
	s.Lock.RUnlock()
	//加入区服,加入房间
	bucket := app.GetBucket(user.Uid)
	bucket.Area(user.AreaId).JoinArea(user)
	bucket.Room(user.RoomId).JoinRoom(user)
	//不断读消息
	var sendType comet.Type
	for err = proto.DecodeFromBytes(user.ReadBuf); err == nil; {
		switch proto.Op {
		case protocol.OpHeartbeat:
			//todo 添加到最小堆
			continue
		case protocol.OpDisconnect:
			//删除用户对应信息
			return
		case protocol.OpSendAreaMsg | protocol.OpSendRoomMsg | protocol.OpSendMsg:
			if proto.Op == protocol.OpSendAreaMsg {
				sendType = comet.Type_AREA
			}
			if proto.Op == protocol.OpSendRoomMsg {
				sendType = comet.Type_ROOM
			}
			if proto.Op == protocol.OpSendMsg {
				sendType = comet.Type_PUSH
			}
			_, err = s.logic.OnMessage(ctx, &logic.MessageReq{
				Type:     sendType,
				SendId:   user.Uid,
				ToId:     user.AreaId,
				Msg:      proto.Data,
				CometKey: s.Name,
			})
			if err != nil {
				s.log.Errorf("s.logic.OnMessage(ctx, &logic.MessageReq{Type:comet.Type_AREA,SendId:%d,ToId:%d,Msg:%s,CometKey:%s}) error(%v)", user.Uid, user.AreaId, string(proto.Data), s.Name, err)
				return
			}
		default:
			//todo send msg to user
			s.log.Errorf("protocol.Op error(%v)", proto.Op)
		}

	}
}

const (
	// grpc options
	grpcInitialWindowSize     = 1 << 24
	grpcInitialConnWindowSize = 1 << 24
	grpcMaxSendMsgSize        = 1 << 24
	grpcMaxCallMsgSize        = 1 << 24
	grpcKeepAliveTime         = time.Second * 10
	grpcKeepAliveTimeout      = time.Second * 3
)

func newLogicClient(ctx context.Context, c *config.CometConfigLogicClientGrpc) logic.LogicClient {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(c.Timeout))
	defer cancel()
	//todo 改成服务发现,接口获取 抽象etcd等
	//grpc.WithInitialWindowSize(grpcInitialWindowSize),
	//			grpc.WithInitialConnWindowSize(grpcInitialConnWindowSize),
	//			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcMaxCallMsgSize)),
	//			grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(grpcMaxSendMsgSize)),
	//			grpc.WithKeepaliveParams(keepalive.ClientParameters{
	//				Time:                grpcKeepAliveTime,
	//				Timeout:             grpcKeepAliveTimeout,
	//				PermitWithoutStream: true,
	//			}),
	conn, err := grpc.Dial(c.Addr,
		grpc.WithInitialWindowSize(grpcInitialWindowSize),
		grpc.WithInitialConnWindowSize(grpcInitialConnWindowSize),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcMaxCallMsgSize)),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(grpcMaxSendMsgSize)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                grpcKeepAliveTime,
			Timeout:             grpcKeepAliveTimeout,
			PermitWithoutStream: true,
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()), //todo 安全验证
	)
	if err != nil {
		panic(err)
	}
	return logic.NewLogicClient(conn)
}
