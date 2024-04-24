package grpc

import (
	"context"
	"github.com/Verce11o/resume-view/echo-service/internal/config"
	pb "github.com/Verce11o/resume-view/protos/gen/go"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"strconv"
	"time"
)

func NewViewServiceClient(ctx context.Context, cfg *config.Config) (pb.ViewServiceClient, error) {
	timeout, err := time.ParseDuration(cfg.ClientTimeout)

	if err != nil {
		return nil, err
	}

	retries, err := strconv.Atoi(cfg.RetriesCount)

	if err != nil {
		return nil, err
	}

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.Aborted, codes.DeadlineExceeded, codes.Unavailable),
		grpcretry.WithMax(uint(retries)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	cc, err := grpc.DialContext(ctx, cfg.ViewServiceEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		))

	if err != nil {
		return nil, err
	}

	return pb.NewViewServiceClient(cc), nil
}

func InterceptorLogger() grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
		slog.Log(ctx, slog.Level(lvl), msg, fields)
	})
}
