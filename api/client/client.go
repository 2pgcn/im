package client

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware/circuitbreaker"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/registry"
	trgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"time"
)

const (
	// grpc options
	grpcInitialWindowSize     = 1 << 24
	grpcInitialConnWindowSize = 1 << 24
	grpcMaxSendMsgSize        = 1 << 24
	grpcMaxCallMsgSize        = 1 << 24
	grpcKeepAliveTime         = time.Second * 10
	grpcKeepAliveTimeout      = time.Second * 3
)

// NewGrpcClient todo 连接池 默认单连接qps10万以内,超过后可多次注册d增加连接池连接
func NewGrpcClient(ctx context.Context, addr string, d registry.Discovery, options ...trgrpc.ClientOption) (*grpc.ClientConn, error) {
	ops := []trgrpc.ClientOption{
		trgrpc.WithEndpoint(addr),
		trgrpc.WithTimeout(time.Second * 3),
		trgrpc.WithMiddleware(
			recovery.Recovery(),
			circuitbreaker.Client(),
		),
		trgrpc.WithOptions(
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithInitialWindowSize(grpcInitialWindowSize),
			grpc.WithInitialConnWindowSize(grpcInitialConnWindowSize),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcMaxCallMsgSize)),
			grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(grpcMaxSendMsgSize)),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                grpcKeepAliveTime,
				Timeout:             grpcKeepAliveTimeout,
				PermitWithoutStream: true,
			}),
			grpc.WithConnectParams(grpc.ConnectParams{
				Backoff: backoff.Config{
					BaseDelay:  time.Second * 1,
					Multiplier: 1.6,
					Jitter:     0.2,
					MaxDelay:   5 * time.Second,
				}}),
		),
	}
	if d != nil {
		ops = append(ops, trgrpc.WithDiscovery(d))
	}
	ops = append(ops, options...)
	return trgrpc.DialInsecure(
		ctx,
		ops...,
	)
}
