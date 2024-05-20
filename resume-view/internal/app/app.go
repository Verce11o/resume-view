package app

import (
	"context"
	"fmt"
	"net"

	"github.com/Verce11o/resume-view/resume-view/internal/config"
	viewgrpc "github.com/Verce11o/resume-view/resume-view/internal/grpc"
	"github.com/Verce11o/resume-view/resume-view/internal/repositories"
	"github.com/Verce11o/resume-view/resume-view/internal/services"
	postgresLib "github.com/Verce11o/resume-view/shared/db/postgres"
	"github.com/Verce11o/resume-view/shared/tracer"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	cfg        *config.Config
	log        *zap.SugaredLogger
	grpcServer *grpc.Server
}

func New(ctx context.Context, cfg *config.Config, log *zap.SugaredLogger) (*App, error) {
	trace, err := tracer.InitTracer(ctx, "view service", cfg.Jaeger.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to init tracer: %w", err)
	}

	db, err := postgresLib.New(ctx, postgresLib.Config{
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Database: cfg.DB.Name,
		SSLMode:  cfg.DB.SSLMode,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

	repo := repositories.NewViewRepository(db, trace)
	service := services.NewViewService(log, trace, repo)

	server := grpc.NewServer(
		grpc.StatsHandler(
			otelgrpc.NewServerHandler(
				otelgrpc.WithTracerProvider(trace.Provider),
				otelgrpc.WithPropagators(propagation.TraceContext{}),
			),
		))

	viewgrpc.Register(log, service, server, trace.Tracer)

	return &App{
		cfg:        cfg,
		log:        log,
		grpcServer: server,
	}, nil
}

func (a *App) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", a.cfg.Server.Port))

	if err != nil {
		return fmt.Errorf("failed to listen tcp: %w", err)
	}

	if err := a.grpcServer.Serve(l); err != nil {
		return fmt.Errorf("failed to serve grpc: %w", err)
	}

	return nil
}

func (a *App) Stop() {
	a.grpcServer.GracefulStop()
}
