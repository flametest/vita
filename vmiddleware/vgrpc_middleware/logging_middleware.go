package vgrpc_middleware

import (
	"context"
	"time"

	log "github.com/flametest/vita/vlog"
	"google.golang.org/grpc"
)

func LoggingMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()
	m, err := handler(ctx, req)
	if err != nil {
		log.WithCtx(ctx).Error().Msgf("[grpc-server] call handler error: %v", err)
	}
	log.WithCtx(ctx).Info().Msgf("[grpc-server] latency=%vms, method=%v, from=%v", time.Since(start).Milliseconds(), info.FullMethod, ctx.Value("from"))
	return m, err
}
