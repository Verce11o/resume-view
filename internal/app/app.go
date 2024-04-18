package app

import (
	"context"
	"fmt"
	"github.com/Verce11o/resume-view/internal/config"
	viewgrpc "github.com/Verce11o/resume-view/internal/grpc"
	"github.com/Verce11o/resume-view/internal/repository"
	"github.com/Verce11o/resume-view/internal/service"
	postgresLib "github.com/Verce11o/resume-view/lib/db/postgres"
	"github.com/Verce11o/resume-view/lib/tracer"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
)

type App struct {
	cfg        *config.Config
	log        *zap.SugaredLogger
	grpcServer *grpc.Server
}

func New(ctx context.Context, cfg *config.Config, log *zap.SugaredLogger) *App {
	trace := tracer.InitTracer("view service", cfg.Jaeger.Endpoint)
	db := postgresLib.Run(ctx, cfg)

	repo := repository.NewViewRepository(db, trace)
	services := service.NewViewService(log, trace, repo)

	server := grpc.NewServer(
		grpc.StatsHandler(
			otelgrpc.NewServerHandler(
				otelgrpc.WithTracerProvider(trace.Provider),
				otelgrpc.WithPropagators(propagation.TraceContext{}),
			),
		))

	viewgrpc.Register(log, services, server, trace.Tracer)

	return &App{
		cfg:        cfg,
		log:        log,
		grpcServer: server,
	}
}

func (a *App) Run() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", a.cfg.Server.Port))

	if err != nil {
		log.Fatalf("error while listen grpc server: %v", err)
	}

	if err := a.grpcServer.Serve(l); err != nil {
		log.Fatalf("error while serve grpc server: %v", err)
	}
}

func (a *App) Stop() {
	a.grpcServer.GracefulStop()
}
