package server

import (
	v1 "github.com/2pgcn/gameim/api/logic"
	"github.com/2pgcn/gameim/conf"
	"github.com/2pgcn/gameim/internal/logic/service"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	grpc2 "google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"time"
)

const (
	// grpc options
	grpcInitialWindowSize     = 1 << 24
	grpcInitialConnWindowSize = 1 << 24
	grpcMaxSendMsgSize        = 1 << 24
	grpcMaxCallMsgSize        = 1 << 24
	grpcKeepAliveTime         = time.Second * 60
	grpcKeepAliveTimeout      = time.Second * 30
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, auth *service.AuthService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			//tracing.Server(
			//	tracing.WithTracerProvider(otel.GetTracerProvider()),
			//	tracing.WithTracerName("gameim"),
			//),
			recovery.Recovery(),
			ratelimit.Server(),
		),

		grpc.Options(
			grpc2.KeepaliveParams(keepalive.ServerParameters{
				Time:    grpcKeepAliveTime,
				Timeout: grpcKeepAliveTimeout,
			}),
			grpc2.MaxRecvMsgSize(grpcMaxCallMsgSize),
			grpc2.MaxSendMsgSize(grpcMaxSendMsgSize),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterLogicServer(srv, auth)
	return srv
}
