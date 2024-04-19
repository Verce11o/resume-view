package app

import (
	"context"
	"fmt"
	"github.com/Verce11o/resume-view/internal/config"
	viewgrpc "github.com/Verce11o/resume-view/internal/grpc"
	"github.com/Verce11o/resume-view/internal/repositories"
	"github.com/Verce11o/resume-view/internal/services"
	postgresLib "github.com/Verce11o/resume-view/lib/db/postgres"
	"github.com/Verce11o/resume-view/lib/tracer"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type App struct {
	cfg        *config.Config
	log        *zap.SugaredLogger
	grpcServer *grpc.Server
}

func New(ctx context.Context, cfg *config.Config, log *zap.SugaredLogger) (*App, error) {
	trace, err := tracer.InitTracer(ctx, "view service", cfg.Jaeger.Endpoint)
	if err != nil {
		return nil, err
	}

	db, err := postgresLib.Run(ctx, cfg)
	if err != nil {
		return nil, err
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
		return err
	}

	if err := a.grpcServer.Serve(l); err != nil {
		return err
	}

	return nil
}

func (a *App) Stop() {
	a.grpcServer.GracefulStop()
}
