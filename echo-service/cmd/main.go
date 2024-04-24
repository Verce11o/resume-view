package main

import (
	"context"
	"github.com/Verce11o/resume-view/echo-service/internal/clients/grpc"
	"github.com/Verce11o/resume-view/echo-service/internal/config"
	pb "github.com/Verce11o/resume-view/protos/gen/go"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	cfg := config.Load()
	client, err := grpc.NewViewServiceClient(ctx, logger, cfg)

	if err != nil {
		slog.Error("Error creating new view client: %v", err)
	}

	ticker := time.NewTicker(5 * time.Second)

	slog.Info("Starting view client")
	for {
		select {
		case <-ticker.C:
			resp, _ := client.GetResumeViews(ctx, &pb.GetResumeViewsRequest{})
			slog.Info("Total views: ", "views", resp.GetTotal())
		case <-ctx.Done():
			slog.Info("Stopping echo-service")
			return
		}
	}

}
