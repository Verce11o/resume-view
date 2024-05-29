package server

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func CorrelationInterceptor(log *zap.SugaredLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp any, err error) {
		correlationID := ctx.Value(keyCorrelationID)
		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		ctx = context.WithValue(ctx, keyCorrelationID, correlationID)

		start := time.Now()

		m, err := handler(ctx, req)

		duration := time.Since(start)

		log.Infof("RPC: %s, duration: %v, err: %v", info.FullMethod, duration, err)

		return m, err
	}
}
