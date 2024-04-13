package comet

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"github.com/2pgcn/gameim/api/client"
	"github.com/2pgcn/gameim/api/gerr"
	"github.com/2pgcn/gameim/api/logic"
	"github.com/2pgcn/gameim/api/protocol"
	"github.com/2pgcn/gameim/conf"
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

	//rcvQueue, err := event.NewSockReceiver(c.Queue.GetSock())
	if err != nil {
		return server, gerr.ErrorServerError("NewServer").WithMetadata(gerr.GetStack()).WithCause(err)
	}
	//启动所有app,目前从配置里读,改成从数据中心取,可类似traefik config模式改成动态配置启动
	for _, v := range c.GetAppConfig() {
		tmpApp, err := NewApp(ctx, v, c.Queue, log, server.gopool)
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
	var err error
	if addr, err = net.ResolveTCPAddr("tcp", host); err != nil {
		s.GetLog().Errorf("net.ResolveTCPAddr(tcp, %s) error(%v)", host, err)
		return
	}
	if listener, err = net.ListenTCP("tcp", addr); err != nil {
		s.GetLog().Errorf("net.ListenTCP(tcp, %s) error(%v)", addr, err)
		return
	}
	s.Lock.Lock()
	s.listens = append(s.listens, listener)
	s.Lock.Unlock()
	for {
		conn, err := listener.AcceptTCP()
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
		_ = conn.Close()
		return gerr.ErrorServerError("protocol.DecodeFromBytes error(%s),userData", err).WithCause(err).WithMetadata(gerr.GetStack())
	}
	if p.Op != protocol.OpAuth {
		p.SetErrReply(gerr.ErrorAuthError("protocol.Op != protocol.OpLogin"))
		connExit(conn, bw, p)
		return nil
	}
	AuthToken := &protocol.Auth{}
	if err := json.Unmarshal(p.Data, AuthToken); err != nil {
		p.SetErrReply(gerr.ErrorAuthError("Unmarshal token data error").WithCause(err))
		connExit(conn, bw, p)
		return nil
	}
	//仅在支持在线配置时需要,提前写上
	s.Lock.RLock()
	app, ok := s.apps[AuthToken.Appid]
	s.Lock.RUnlock()
	if !ok {
		p.SetErrReply(gerr.ErrorServerError("s.logic.OnAuth error"))
		connExit(conn, bw, p)
		return gerr.ErrorServerError("s.logic.OnAuth error").WithCause(err).WithMetadata(gerr.GetStack())
	}
	//进入app逻辑后
	app.AddUser(AuthToken.Token, conn, br, bw)
	return nil
}

func connExit(conn *net.TCPConn, bw *bufio.Writer, p *protocol.Proto) {
	_ = p.WriteTcp(bw)
	_ = conn.Close()
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

// todo,cc未清理
func newLogicClient(ctx context.Context, c *conf.LogicClientGrpc) logic.LogicClient {
	cc, err := client.NewGrpcClient(ctx, c.Addr, nil)
	if err != nil {
		panic(err)
	}
	return logic.NewLogicClient(cc)
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
