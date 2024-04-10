package comet

import (
	"bufio"
	"context"
	"errors"
	"github.com/2pgcn/gameim/api/client"
	"github.com/2pgcn/gameim/api/gerr"
	"github.com/2pgcn/gameim/api/logic"
	"github.com/2pgcn/gameim/api/protocol"
	"github.com/2pgcn/gameim/conf"
	"github.com/2pgcn/gameim/pkg/event"
	"github.com/2pgcn/gameim/pkg/gamelog"
	"github.com/2pgcn/gameim/pkg/safe"

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

func (s *Server) GetLog() gamelog.GameLog {
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
	server.log = log.AppendPrefix("server")
	//todo 工厂模式适配多个队列
	//rcvQueue, err := event.NewNsqReceiver(c.Queue.GetNsq())
	rcvQueue, err := event.NewSockReceiver(c.Queue.GetSock())
	if err != nil {
		return server, gerr.ErrorServerError("NewServer").WithMetadata(gerr.GetStack()).WithCause(err)
	}
	//启动所有app,目前从配置里读,改成从数据中心取,可类似traefik config模式改成动态配置启动
	for _, v := range c.GetAppConfig() {
		tmpApp, err := NewApp(ctx, v, rcvQueue, log, server.gopool)
		if err != nil {
			return server, gerr.ErrorServerError("NewServer").WithMetadata(gerr.GetStack()).WithCause(err)
		}
		server.apps[v.Appid] = tmpApp
		err = tmpApp.Start()
		if err != nil {
			err = nil
			gamelog.GetGlobalog().Errorf("app(%s) start error:%s", v.Appid, err)
			continue
		}
		server.log.Debugf("step1:appid is started:%+v", v.Appid)
	}

	server.logic = newLogicClient(ctx, c.LogicClientGrpc)
	for _, v := range c.Server.Addrs {
		addr := v
		server.gopool.GoCtx(func(ctx context.Context) {
			server.bindConn(ctx, addr, c.Tcp)
		})
	}
	return server, nil
}

func (s *Server) bindConn(ctx context.Context, host string, c *conf.TcpMsg) {
	var addr *net.TCPAddr
	var listener *net.TCPListener
	var conn *net.TCPConn
	var err error
	if addr, err = net.ResolveTCPAddr("tcp", host); err != nil {
		s.GetLog().Errorf("net.ResolveTCPAddr(tcp, %s) error(%v)", host, err)
		return
	}

	if listener, err = net.ListenTCP("tcp", addr); err != nil {
		s.GetLog().Errorf("net.ListenTCP(tcp, %s) error(%v)", addr, err)
		return
	}
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
			err = s.handleComet(ctx, conn, c)
			if err != nil {
				s.GetLog().Info(err)
			}
		})
	}

}

func (s *Server) handleComet(ctx context.Context, conn *net.TCPConn, c *conf.TcpMsg) error {
	//defer func() {
	//	err := conn.Close()
	//	if err != nil {
	//		s.GetLog().Error("conn close err:%s", err)
	//	}
	//}()
	if err := conn.SetReadBuffer(int(c.SendBuf)); err != nil {
		return gerr.ErrorServerError("conn.SetReadBuffer err").WithCause(err).WithMetadata(gerr.GetStack())
	}
	if err := conn.SetWriteBuffer(int(c.SendBuf)); err != nil {
		return gerr.ErrorServerError("conn.SetWriteBuffer err").WithCause(err).WithMetadata(gerr.GetStack())
	}
	//todo 可以改成从自定义byte pool里拿,减少gc
	br := bufio.NewReader(conn)
	bw := bufio.NewWriter(conn)
	p := protocol.ProtoPool.Get()
	defer func() {
		p.Reset()
		protocol.ProtoPool.Put(p)
	}()
	err := p.DecodeFromBytes(br)
	if err != nil {
		return gerr.ErrorServerError("protocol.DecodeFromBytes error(%s),userData", err).WithCause(err).WithMetadata(gerr.GetStack())
	}
	if p.Op != protocol.OpAuth {
		p.SetErrReply(gerr.ErrorAuthError("protocol.Op != protocol.OpLogin"))
		return p.WriteTcp(bw)
	}
	//todo 配置超时时间
	//grpcCtx, cancel := context.WithTimeout(ctx, s.conf.LogicClientGrpc.Timeout.AsDuration())
	grpcCtx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	authReply, err := s.logic.OnAuth(grpcCtx, &logic.AuthReq{
		Token: string(p.Data),
	})
	if err != nil {
		p.SetErrReply(gerr.ErrorServerError("s.logic.OnAuth error"))
		_ = p.WriteTcp(bw)
		return gerr.ErrorServerError("s.logic.OnAuth error").WithCause(err).WithMetadata(gerr.GetStack())
	}
	//appid 先写死,后面通过proto.data里解码获取
	//userInfo authReply
	var app *App
	var ok bool
	s.Lock.RLock()
	app, ok = s.apps[authReply.Appid]
	s.Lock.RUnlock()
	if !ok {
		p.SetErrReply(gerr.ErrorAuthAppidError("auth error:app_id:%s not found ", authReply.Appid))
		return p.WriteTcp(bw)
	}
	user := NewUser(ctx, conn, s.log)
	user.ReadBuf = br
	user.WriteBuf = bw
	user.Uid = userId(authReply.Uid)
	user.RoomId = roomId(authReply.RoomId)
	bucket := app.GetBucket(user.Uid)
	bucket.PutUser(user)
	var writeProto *protocol.Proto
	writeProto, err = protocol.NewProtoMsg(protocol.OpAuthReply, authReply)
	if err = writeProto.WriteTcp(user.WriteBuf); err != nil {
		return err
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
			gamelog.GetGlobalog().Debugf("client DecodeFromBytes error,%s", err)
			return nil
		}
		//msgCtx, span := trace_conf.SetTrace(context.Background(), trace_conf.COMET_RECV_CIENT_MSG,
		//	trace.WithSpanKind(trace.SpanKindInternal), trace.WithAttributes(
		//		attribute.Int("type", int(p.Op)),
		//		attribute.String("data", string(p.Data)),
		//	))
		//span.End()
		msgCtx, _ := context.WithTimeout(context.Background(), time.Second*10)
		msgP := event.GetMsg()
		err = msgP.Unmarshal(p.Data)
		if err != nil {
			replyErr := gerr.ErrorMsgFormatError("p.DecodeFromBytes error")
			p.SetErrReply(replyErr)
			err1 := p.WriteTcp(user.WriteBuf)
			return replyErr.WithCause(err1).WithMetadata(gerr.GetStack())
		}
		gamelog.GetGlobalog().Debugf("recv client msg:%v,:srcMsg:%v", msgP, p)
		switch p.Op {
		case protocol.OpHeartbeat:
			// 更新最小堆
			bucket.lock.RLock()
			bucket.UpHeartTime(user.Uid, time.Now().Add(c.KeepaliveTimeout.AsDuration()))
			bucket.lock.Unlock()
			continue
		case protocol.OpDisconnect:
			//心跳堆无需删除,在取出的时候兼容
			bucket.DeleteUser(user.Uid)
			//删除用户对应信息
			//bucket.DeleteUser(user.Uid)
			return nil
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

			req := &logic.MessageReq{
				Type:     sendType,
				SendId:   string(user.Uid),
				ToId:     msgP.ToId, //p.Data
				Msg:      p.Data,
				CometKey: s.Name,
			}
			_, err = s.logic.OnMessage(msgCtx, req)
			if err != nil {
				return gerr.ErrorServerError("s.logic.OnMessage err,data(%v)", req).WithCause(err).WithMetadata(gerr.GetStack())
			}
		}
	}
	return nil
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
	cc, err := client.NewGrpcClient(ctx, c.Addr, nil)
	if err != nil {
		panic(err)
	}
	return logic.NewLogicClient(cc)
	//provider, err := trace_conf.GetTracerProvider()
	//if err != nil {
	//	panic(err)
	//}
	//conn, err := trgrpc.DialInsecure(
	//	ctx,
	//	trgrpc.WithEndpoint(c.Addr),
	//	trgrpc.WithTimeout(time.Second*10),
	//	trgrpc.WithMiddleware(
	//		//tracing.Client(tracing.WithTracerProvider(provider), tracing.WithTracerName("gameim")),
	//		recovery.Recovery(),
	//		//circuitbreaker.Client(),
	//	),
	//	trgrpc.WithOptions(
	//		grpc.WithInitialWindowSize(grpcInitialWindowSize),
	//		grpc.WithInitialWindowSize(grpcInitialWindowSize),
	//		grpc.WithInitialConnWindowSize(grpcInitialConnWindowSize),
	//		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcMaxCallMsgSize)),
	//		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(grpcMaxSendMsgSize)),
	//		grpc.WithKeepaliveParams(keepalive.ClientParameters{
	//			Time:                grpcKeepAliveTime,
	//			Timeout:             grpcKeepAliveTimeout,
	//			PermitWithoutStream: true,
	//		}),
	//		grpc.WithConnectParams(grpc.ConnectParams{
	//			Backoff: backoff.Config{
	//				BaseDelay:  time.Second * 1,
	//				Multiplier: 1.6,
	//				Jitter:     0.2,
	//				MaxDelay:   5 * time.Second,
	//			}}),
	//	),
	//	//grpc.WithTransportCredentials(insecure.NewCredentials()), //todo 安全验证
	//)
	//if err != nil {
	//	panic(err)
	//}
	//gamelog.GetGlobalog().Debugf("start grpc client:%s", c.Addr)

}

func (s *Server) Close() {
	s.GetLog().Debug("listen is stop")
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	for _, v := range s.apps {
		v.Close()
	}
	//关闭listen
	for i := 0; i < len(s.listens); i++ {
		err := s.listens[i].Close()
		if err != nil {
			s.log.Errorf("Server Close err:%s", err.Error())
		}
	}
}
