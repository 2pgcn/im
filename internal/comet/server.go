package comet

import (
	"bufio"
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/php403/gameim/api/comet"
	"github.com/php403/gameim/api/logic"
	"github.com/php403/gameim/api/protocol"
	"github.com/php403/gameim/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"sync"
	"time"
)

const appid = "app01"

type Server struct {
	ctx   context.Context
	conf  *config.CometConfig
	Name  string
	log   *zap.SugaredLogger
	Apps  map[string]*App
	logic logic.LogicClient
	Lock  sync.RWMutex
}

func NewServer(ctx context.Context, log *zap.SugaredLogger, config *config.CometConfig) {
	server := &Server{
		ctx:  ctx,
		Name: config.Server.Name,
		log:  log,
		Apps: map[string]*App{},
		conf: config,
		//todo 绑到app下
		logic: newLogicClient(ctx, config.LogicClientGrpc),
	}
	//启动消费队列消费
	for _, v := range config.Server.Addrs {
		addr := v
		go func() {
			select {
			case <-ctx.Done():
				return
			default:
				server.bindConn(ctx, addr)
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

	for {
		if conn, err = listener.AcceptTCP(); err != nil {
			s.log.Errorf("accept comet error %v", zap.Error(err))
			return
		}
		go s.handleComet(ctx, conn)
	}
}

func (s *Server) handleComet(ctx context.Context, conn net.Conn) {

	user := NewUser()
	user.ReadBuf = bufio.NewReader(conn)
	user.WriteBuf = bufio.NewWriter(conn)
	//todo 用户登录后加入最小堆(nginx采用红黑树event)
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
	grpcCtx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()
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
	user.Uid = authReply.Uid
	user.AreaId = authReply.AreaId
	user.RoomId = authReply.RoomId
	var app *App
	var ok bool
	s.Lock.Lock()
	if app, ok = s.Apps[appid]; !ok {
		app, err = NewApp(ctx, appid, s.conf.Queue, s.log)
		if err != nil {
			s.log.Errorf("NewApp(ctx, %s, %s) error(%v)", appid, s.conf.Queue, err)
			return
		}
		s.Apps[appid] = app
	}
	s.log.Debugf("test appid %s", appid)
	s.Lock.Unlock()
	//加入区服,加入房间
	bucket := app.GetBucket(user.Uid)
	//bucket.users
	bucket.PutUser(user)
	//发送用户登录成功消息
	var writeProto *protocol.Proto
	writeProto = &protocol.Proto{
		Version:  1,
		Op:       protocol.OpAuthReply,
		Checksum: 0,
		Seq:      0, //todo
	}
	s.log.Debugf("auth success  : %v", authReply)
	writeProto.Data, _ = proto.Marshal(authReply)
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
			case msg := <-user.msgQueue:
				data, err := proto.Marshal(msg)
				s.log.Debugf("recv msg : %v", msg)
				if err != nil {
					s.log.Errorf("proto.Marshal(msg) error(%v)", err)
					break
				}
				writeProto = &protocol.Proto{
					Version:  1,
					Op:       protocol.OpSendAreaMsg,
					Checksum: 0,
					Seq:      0, //todo
					Data:     data,
				}
				s.log.Debugf("writeProto : %v", writeProto)
				if err = writeProto.WriteTcp(user.WriteBuf); err != nil {
					s.log.Errorf("writeProto.EncodeTo(user.WriteBuf) error(%v)", err)
					break
				}
			}
		}
	}()
	//不断读消息
	var sendType comet.Type
	for {
		err = p.DecodeFromBytes(user.ReadBuf)
		if err != nil {
			continue
		}
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
				sendType = comet.Type_AREA
			}
			if p.Op == protocol.OpSendRoomMsg {
				p.Op = protocol.OpSendMsgRoomReply
				sendType = comet.Type_ROOM
			}
			if p.Op == protocol.OpSendMsg {
				p.Op = protocol.OpSendMsgReply
				sendType = comet.Type_PUSH
			}
			_, err = s.logic.OnMessage(ctx, &logic.MessageReq{
				Type:     sendType,
				SendId:   user.Uid,
				ToId:     user.AreaId,
				Msg:      p.Data,
				CometKey: s.Name,
			})
			if err != nil {
				s.log.Errorf("s.logic.OnMessage(ctx, &logic.MessageReq{Type:comet.Type_AREA,SendId:%d,ToId:%d,Msg:%s,CometKey:%s}) error(%v)", user.Uid, user.AreaId, string(p.Data), s.Name, err)
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
		//grpc.WithKeepaliveParams(keepalive.ClientParameters{
		//	Time:                grpcKeepAliveTime,
		//	Timeout:             grpcKeepAliveTimeout,
		//	PermitWithoutStream: true,
		//}),
		grpc.WithTransportCredentials(insecure.NewCredentials()), //todo 安全验证
	)
	if err != nil {
		panic(err)
	}
	return logic.NewLogicClient(conn)
}
