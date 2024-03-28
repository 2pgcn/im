package comet

import (
	"bufio"
	"context"
	"errors"
	"github.com/2pgcn/gameim/api/gerr"
	"github.com/2pgcn/gameim/api/logic"
	"github.com/2pgcn/gameim/api/protocol"
	"github.com/2pgcn/gameim/conf"
	"github.com/2pgcn/gameim/pkg/event"
	"github.com/2pgcn/gameim/pkg/gamelog"
	"github.com/2pgcn/gameim/pkg/safe"
	"github.com/2pgcn/gameim/pkg/trace_conf"
	"github.com/go-kratos/kratos/v2/middleware/circuitbreaker"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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
	gopool  *safe.GoPool
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

func NewServer(ctx context.Context, c *conf.CometConfig, gopool *safe.GoPool, log gamelog.GameLog) (server *Server, err error) {
	server = &Server{
		gopool: gopool,
		ctx:    ctx,
		Name:   c.Server.Name,
		apps:   map[string]*App{},
		conf:   c,
	}
	server.log = log.ReplacePrefix(gamelog.DefaultMessageKey + ":server")
	rcvQueue, err := event.NewNsqReceiver(c.Queue.GetNsq())
	if err != nil {
		return server, gerr.ErrorServerError("NewServer").WithMetadata(gerr.GetStack()).WithCause(err)

	}
	//启动所有app,目前从配置里读,改成从数据中心取
	for _, v := range c.GetAppConfig() {
		tmpApp, err := NewApp(ctx, v, rcvQueue, log, server.gopool)
		if err != nil {
			return server, gerr.ErrorServerError("NewServer").WithMetadata(gerr.GetStack()).WithCause(err)
		}
		server.apps[v.Appid] = tmpApp
		server.log.Debugf("appid is starting:%+v", v.Appid)
	}
	server.logic = newLogicClient(ctx, c.LogicClientGrpc)
	for _, v := range c.Server.Addrs {
		addr := v
		server.gopool.GoCtx(func(ctx context.Context) {
			server.bindConn(ctx, addr)
		})
	}
	return server, nil
}

func (s *Server) bindConn(ctx context.Context, host string) {
	var addr *net.TCPAddr
	var listener *net.TCPListener
	var conn *net.TCPConn
	var err error
	if addr, err = net.ResolveTCPAddr("tcp", host); err != nil {
		s.logHelper().Errorf("net.ResolveTCPAddr(tcp, %s) error(%v)", host, err)
		return
	}

	if listener, err = net.ListenTCP("tcp", addr); err != nil {
		s.logHelper().Errorf("net.ListenTCP(tcp, %s) error(%v)", addr, err)
		return
	}
	s.logHelper().Debugf("server listen:%s", addr)
	s.listens = append(s.listens, listener)
	for {
		conn, err = listener.AcceptTCP()
		if err != nil {
			if !errors.Is(err, net.ErrClosed) {
				s.log.Infof("listener is close: accept:%s", err.Error())
			}
			return
		}
		s.gopool.GoCtx(func(ctx context.Context) {
			s.handleComet(ctx, conn)
		})
	}

}

func (s *Server) handleComet(ctx context.Context, conn *net.TCPConn) {
	user := NewUser(ctx, conn, s.log)
	//todo 可以改成从自定义pool里拿,减少gc
	user.ReadBuf = bufio.NewReader(conn)
	user.WriteBuf = bufio.NewWriter(conn)

	p := protocol.ProtoPool.Get()
	defer func() {
		p.Reset()
		protocol.ProtoPool.Put(p)
	}()
	//todo pool里获取获取用户消息
	//var operation = func() error {
	//	return p.DecodeFromBytes(user.ReadBuf)
	//}
	//err := backoff.Retry(safe.OperationWithRecover(operation), backoff.NewExponentialBackOff())
	err := p.DecodeFromBytes(user.ReadBuf)
	if err != nil {
		s.log.Errorf("protocol.DecodeFromBytes(user.ReadBuf, &user.Msg) error(%v)", err)
		user.Close()
		return
	}
	s.logHelper().Debugf("server msg recv:%+v", p)
	if p.Op != protocol.OpAuth {
		//todo return client err msg
		err = conn.Close()
		if err != nil {
			s.log.Errorf("conn.Close() error(%v)", err)
		}
		s.log.Error("protocol.Op != protocol.OpLogin,msg(%)", p)
		return
	}
	//todo 配置超时时间
	//grpcCtx, cancel := context.WithTimeout(ctx, s.conf.LogicClientGrpc.Timeout.AsDuration())
	grpcCtx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	authReply, err := s.logic.OnAuth(grpcCtx, &logic.AuthReq{
		Token: string(p.Data),
	})
	s.log.Debug("protocol.OpAuth pass")
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
	s.Lock.RLock()
	app, ok = s.apps[authReply.Appid]
	s.Lock.RUnlock()
	s.log.Debugf("auth success:reply message:%+v", authReply)
	if !ok {
		// todo
		//p.SetError(errors2.TypeStatusCode_AUTH_ERROR, protocol.OpAuthReply)
		err = p.WriteTcp(user.WriteBuf)
		if err != nil {
			user.log.Errorf("auth error:app_id:%s not found ", user.AppId)
			return
		}
	}
	s.log.Debugf("test appid %s", authReply.Appid)
	bucket := app.GetBucket(user.Uid)
	bucket.PutUser(user)
	var writeProto *protocol.Proto
	writeProto, err = protocol.NewProtoMsg(protocol.OpAuthReply, authReply)
	if err = writeProto.WriteTcp(user.WriteBuf); err != nil {
		s.log.Errorf("writeProto.EncodeTo(user.WriteBuf) error(%v)", err)
		return
	}
	//启动用户获取自己msg
	s.gopool.GoCtx(func(ctx context.Context) {
		user.Start()
	})

	//不断读消息
	var sendType protocol.Type
	for {
		err = p.DecodeFromBytes(user.ReadBuf)
		if err != nil {
			if errors.Is(err, protocol.ErrInvalidBuffer) || errors.Is(err, protocol.ErrEOFData) {
				continue
			}
			s.log.Error(err)
			break
		}
		msgCtx, span := trace_conf.SetTrace(context.Background(), trace_conf.COMET_RECV_CIENT_MSG,
			trace.WithSpanKind(trace.SpanKindInternal), trace.WithAttributes(
				attribute.Int("type", int(p.Op)),
				attribute.String("data", string(p.Data)),
			))
		span.End()
		switch p.Op {
		case protocol.OpHeartbeat:
			// 更新最小堆
			bucket.lock.RLock()
			bucket.UpHeartTime(user.Uid, time.Now().Add(time.Second*60*3))
			bucket.lock.Unlock()
			continue
		case protocol.OpDisconnect:
			//心跳堆无需删除,在取出的时候兼容
			bucket.DeleteUser(user.Uid)
			//删除用户对应信息
			return
		case protocol.OpSendAreaMsg, protocol.OpSendRoomMsg, protocol.OpSendMsg:
			if p.Op == protocol.OpSendAreaMsg {
				p.Op = protocol.OpSendAreaMsgReply
				sendType = protocol.Type_APP
			}
			if p.Op == protocol.OpSendRoomMsg {
				p.Op = protocol.OpSendMsgRoomReply
				sendType = protocol.Type_ROOM
			}
			if p.Op == protocol.OpSendMsg {
				p.Op = protocol.OpSendMsgReply
				sendType = protocol.Type_PUSH
			}

			_, err = s.logic.OnMessage(msgCtx, &logic.MessageReq{
				Type:     sendType,
				SendId:   string(user.Uid),
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
		trgrpc.WithTimeout(time.Second*3),
		trgrpc.WithMiddleware(
			tracing.Client(tracing.WithTracerProvider(provider), tracing.WithTracerName("gameim")),
			recovery.Recovery(),
			circuitbreaker.Client(),
		),
		//trgrpc.WithOptions(
		//	grpc.WithInitialWindowSize(grpcInitialWindowSize),
		//	grpc.WithInitialWindowSize(grpcInitialWindowSize),
		//	grpc.WithInitialConnWindowSize(grpcInitialConnWindowSize),
		//	grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcMaxCallMsgSize)),
		//	grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(grpcMaxSendMsgSize)),
		//	grpc.WithKeepaliveParams(keepalive.ClientParameters{
		//		Time:                grpcKeepAliveTime,
		//		Timeout:             grpcKeepAliveTimeout,
		//		PermitWithoutStream: true,
		//	}),
		//),
		//grpc.WithTransportCredentials(insecure.NewCredentials()), //todo 安全验证
	)
	if err != nil {
		panic(err)
	}
	gamelog.Debugf("start grpc client:%s", c.Addr)
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
	s.logHelper().Debug("listen is stop")
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	for _, v := range s.apps {
		v.Close()
	}
	s.logHelper().Debug("apps is stop")
}
