package vserver

import (
	"time"

	"github.com/flametest/vita/vmiddleware/vgrpc_middleware"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func NewGrpcServer() *grpc.Server {
	return grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				vgrpc_middleware.LoggingMiddleware,
				vgrpc_middleware.RecoverMiddleware,
			)),
		grpc.KeepaliveEnforcementPolicy(
			keepalive.EnforcementPolicy{
				MinTime:             10 * time.Second,
				PermitWithoutStream: true,
			}),
	)
}
