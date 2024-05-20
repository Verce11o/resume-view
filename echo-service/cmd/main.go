package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Verce11o/resume-view/echo-service/internal/clients/grpc"
	"github.com/Verce11o/resume-view/echo-service/internal/config"
	pb "github.com/Verce11o/resume-view/protos/gen/go"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	defer cancel()

	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: parseLogLevel(cfg.LogLevel)}))
	slog.SetDefault(logger)

	client, err := grpc.NewViewServiceClient(ctx, logger, cfg)

	if err != nil {
		slog.Error("Error creating new view client", "error", err.Error())
	}

	ticker := time.NewTicker(5 * time.Second)

	slog.Info("Starting view client")

	for {
		select {
		case <-ticker.C:
			resp, _ := client.GetResumeViews(ctx, &pb.GetResumeViewsRequest{})
			slog.Info("Total views amount: ", "views", resp.GetTotal())
		case <-ctx.Done():
			slog.Info("Stopping echo-service")

			return
		}
	}
}

func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
