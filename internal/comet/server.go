package comet

import (
	"bufio"
	"context"
	"errors"
	"github.com/2pgcn/gameim/api/comet"
	error2 "github.com/2pgcn/gameim/api/error"
	"github.com/2pgcn/gameim/api/logic"
	"github.com/2pgcn/gameim/api/protocol"
	"github.com/2pgcn/gameim/conf"
	"github.com/2pgcn/gameim/pkg/event"
	"github.com/2pgcn/gameim/pkg/gamelog"
	"github.com/2pgcn/gameim/pkg/trace_conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/circuitbreaker"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	//"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials/insecure"
	//"google.golang.org/grpc/keepalive"
	trgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"net"
	"sync"
	"time"
)

type Server struct {
	ctx     context.Context
	wg      *sync.WaitGroup
	listens []*net.TCPListener
	conf    *conf.CometConfig
	Name    string
	log     gamelog.GameLog
	apps    map[string]*App
	logic   logic.LogicClient
	Lock    sync.RWMutex
}

func (s *Server) logHelper() gamelog.GameLog {
	return s.log
}

func NewServer(ctx context.Context, c *conf.CometConfig, wg *sync.WaitGroup, log gamelog.GameLog) (server *Server, err error) {
	server = &Server{
		wg:   wg,
		ctx:  ctx,
		Name: c.Server.Name,
		apps: map[string]*App{},
		conf: c,
	}
	server.log = log.ReplacePrefix(gamelog.DefaultMessageKey + ":server")
	rcvQueue, err := event.NewKafkaReceiver(ctx, c.Queue.GetKafka())
	if err != nil {
		return server, err
	}
	//启动app
	for _, v := range c.GetAppConfig() {
		tmpApp, err := NewApp(ctx, v, rcvQueue, log)
		if err != nil {
			return server, err
		}
		server.apps[v.Appid] = tmpApp
	}
	server.logic = newLogicClient(ctx, c.LogicClientGrpc)
	//启动消费队列消费
	for _, v := range c.Server.Addrs {
		addr := v
		go server.bindConn(ctx, addr)
	}
	return server, err
}

func (s *Server) bindConn(ctx context.Context, host string) {
	var addr *net.TCPAddr
	var listener *net.TCPListener
	var conn *net.TCPConn
	var err error
	if addr, err = net.ResolveTCPAddr("tcp", host); err != nil {
		s.log.Errorf("net.ResolveTCPAddr(tcp, %s) error(%v)", host, err)
		return
	}

	if listener, err = net.ListenTCP("tcp", addr); err != nil {
		s.log.Errorf("net.ListenTCP(tcp, %s) error(%v)", addr, err)
		return
	}
	s.wg.Add(1)
	defer s.wg.Done()
	s.listens = append(s.listens, listener)
	for {
		conn, err = listener.AcceptTCP()
		if err != nil {
			if !errors.Is(err, net.ErrClosed) {
				s.log.Infof("listener is close: accept:%s", err.Error())
			}
			return
		}
		go s.handleComet(ctx, conn)
	}

}

func (s *Server) handleComet(ctx context.Context, conn *net.TCPConn) {
	user := NewUser(ctx, conn, s.log)
	//todo 可以改成从自定义pool里拿,规避gc
	user.ReadBuf = bufio.NewReader(conn)
	user.WriteBuf = bufio.NewWriter(conn)

	//获取用户消息
	p := &protocol.Proto{}
	if err := p.DecodeFromBytes(user.ReadBuf); err != nil {
		s.log.Errorf("protocol.DecodeFromBytes(user.ReadBuf, &user.Msg) error(%v)", err)
		return
	}
	if p.Op != protocol.OpAuth {
		//todo return client err msg
		err := conn.Close()
		if err != nil {
			s.log.Errorf("conn.Close() error(%v)", err)
		}
		s.log.Errorf("protocol.Op != protocol.OpLogin")
		return
	}
	//todo 配置超时时间
	s.log.Debug("protocol.OpAuth pass")
	grpcCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	//开始记录链路信息
	authReply, err := s.logic.OnAuth(grpcCtx, &logic.AuthReq{
		Token: string(p.Data),
	})
	if err != nil {
		s.log.Errorf("s.logic.OnAuth error(%v)", err)
		err = conn.Close()
		if err != nil {
			s.log.Errorf("conn.Close() error(%v)", err)
			return
		}
		return
	}
	//appid 先写死,后面通过proto.data里解码获取
	//userInfo authReply
	user.Uid = userId(authReply.Uid)
	user.RoomId = roomId(authReply.RoomId)
	var app *App
	var ok bool
	s.Lock.Lock()
	log.Debugf("auth success:reply message:%+v", authReply)
	if app, ok = s.apps[authReply.Appid]; !ok {
		p.SendError(error2.AuthAppIdError, protocol.OpAuthReply)
		err = p.WriteTcp(user.WriteBuf)
		if err != nil {
			user.log.Errorf("auth error:app_id:%s not found ", user.AppId)
			return
		}
	}
	s.log.Debugf("test appid %s", authReply.Appid)
	s.Lock.Unlock()
	bucket := app.GetBucket(user.Uid)
	bucket.PutUser(user)
	var writeProto *protocol.Proto
	writeProto, err = protocol.NewProtoMsg(protocol.OpAuthReply, authReply)
	if err = writeProto.WriteTcp(user.WriteBuf); err != nil {
		s.log.Errorf("writeProto.EncodeTo(user.WriteBuf) error(%v)", err)
		return
	}
	//发送消息
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msgEvent := <-user.msgQueue:
				msg := msgEvent.(*event.Msg)
				writeProto, err = protocol.NewProtoMsg(msg.GetData().GetType().ToOp(), msg.GetData())
				if err != nil {
					s.log.Errorf("writeProto err: %+v", writeProto)
					continue
				}
				if err = writeProto.WriteTcp(user.WriteBuf); err != nil {
					s.log.Errorf("writeProto.EncodeTo(user.WriteBuf) error(%v)", err)
					continue
				}
			}
		}
	}()
	//不断读消息
	//todo
	var sendType comet.Type
	for {
		err = p.DecodeFromBytes(user.ReadBuf)
		if err != nil {
			if errors.Is(err, protocol.ErrInvalidBuffer) {
				continue
			}
			s.log.Errorf("")
			break
		}
		msgCtx, span := trace_conf.SetTrace(context.Background(), trace_conf.COMET_RECV_CIENT_MSG,
			trace.WithSpanKind(trace.SpanKindInternal), trace.WithAttributes(
				attribute.Int("type", int(p.Op)),
				attribute.String("data", string(p.Data)),
			))
		span.End()
		gamelog.Debugf("test:%s+%+v", "test server", p)
		switch p.Op {
		case protocol.OpHeartbeat:
			//todo 添加到最小堆
			continue
		case protocol.OpDisconnect:
			//删除用户对应信息
			return
		case protocol.OpSendAreaMsg, protocol.OpSendRoomMsg, protocol.OpSendMsg:
			if p.Op == protocol.OpSendAreaMsg {
				p.Op = protocol.OpSendAreaMsgReply
				sendType = comet.Type_APP
			}
			if p.Op == protocol.OpSendRoomMsg {
				p.Op = protocol.OpSendMsgRoomReply
				sendType = comet.Type_ROOM
			}
			if p.Op == protocol.OpSendMsg {
				p.Op = protocol.OpSendMsgReply
				sendType = comet.Type_PUSH
			}

			_, err = s.logic.OnMessage(msgCtx, &logic.MessageReq{
				Type:     sendType,
				SendId:   uint64(user.Uid),
				ToId:     string(1), //p.Data
				Msg:      p.Data,
				CometKey: s.Name,
			})
			if err != nil {
				s.log.Errorf("s.logic.OnMessage(ctx, &logic.MessageReq{Type:comet.Type_AREA,SendId:%d,ToId:%d,Msg:%s,CometKey:%s}) error(%v)", user.Uid, string(p.Data), s.Name, err)
				return
			}
		default:
			//todo send msg to user
			s.log.Errorf("protocol.Op error(proto:%v)", p)
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

func newLogicClient(ctx context.Context, c *conf.LogicClientGrpc) logic.LogicClient {
	provider, err := trace_conf.GetTracerProvider()
	if err != nil {
		panic(err)
	}
	conn, err := trgrpc.DialInsecure(
		ctx,
		trgrpc.WithEndpoint(c.Addr),
		trgrpc.WithMiddleware(
			tracing.Client(tracing.WithTracerProvider(provider), tracing.WithTracerName("gameim")),
			recovery.Recovery(),
			circuitbreaker.Client(),
		),
		trgrpc.WithOptions(
			grpc.WithInitialWindowSize(grpcInitialWindowSize),
			grpc.WithInitialWindowSize(grpcInitialWindowSize),
			grpc.WithInitialConnWindowSize(grpcInitialConnWindowSize),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcMaxCallMsgSize)),
			grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(grpcMaxSendMsgSize)),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                grpcKeepAliveTime,
				Timeout:             grpcKeepAliveTimeout,
				PermitWithoutStream: true,
			}),
		),
		//grpc.WithTransportCredentials(insecure.NewCredentials()), //todo 安全验证
	)
	if err != nil {
		panic(err)
	}
	return logic.NewLogicClient(conn)
}

func (s *Server) Close() {
	//关闭listen
	for i := 0; i < len(s.listens); i++ {
		err := s.listens[i].Close()
		if err != nil {
			s.log.Errorf("Server Close err:%s", err.Error())
		}
	}
	s.Lock.Lock()
	defer s.Lock.Unlock()
	for _, v := range s.apps {
		v.Close()
	}

}
