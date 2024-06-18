package app

import (
	"context"
	"fmt"
	"net"

	"github.com/Verce11o/resume-view/resume-view/internal/config"
	viewgrpc "github.com/Verce11o/resume-view/resume-view/internal/handler/grpc"
	kafkaHandler "github.com/Verce11o/resume-view/resume-view/internal/handler/kafka"
	"github.com/Verce11o/resume-view/resume-view/internal/repositories"
	"github.com/Verce11o/resume-view/resume-view/internal/services"
	postgresLib "github.com/Verce11o/resume-view/shared/db/postgres"
	kafkaLib "github.com/Verce11o/resume-view/shared/kafka"
	"github.com/Verce11o/resume-view/shared/tracer"
	kafka "github.com/segmentio/kafka-go"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	cfg        *config.Config
	log        *zap.SugaredLogger
	grpcServer *grpc.Server
	consumer   *kafkaHandler.Consumer
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

	kafkaClient, err := kafkaLib.New(ctx, kafkaLib.Config{
		Host: cfg.Kafka.Host,
		Port: cfg.Kafka.Port,
	})

	if err != nil {
		return nil, fmt.Errorf("could not connect to kafka: %w", err)
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

	consumer := kafkaHandler.NewConsumer(log, kafkaClient, cfg.Kafka.Topic, cfg.Kafka.GroupID)

	viewgrpc.Register(log, service, server, trace.Tracer)

	return &App{
		cfg:        cfg,
		log:        log,
		grpcServer: server,
		consumer:   consumer,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", a.cfg.Server.Port))

	if err != nil {
		return fmt.Errorf("failed to listen tcp: %w", err)
	}

	errCh := make(chan error)

	go func() {
		if err = a.consumer.Consume(ctx, func(_ context.Context, message *kafka.Message) error {
			a.log.Debugf("message on offset %d: %s", message.Offset, message.Value)

			return nil
		}); err != nil {
			errCh <- fmt.Errorf("failed to consume: %w", err)

			return
		}
	}()

	go func() {
		if err = a.grpcServer.Serve(l); err != nil {
			errCh <- fmt.Errorf("failed to serve: %w", err)

			return
		}
	}()

	return <-errCh
}

func (a *App) Stop() error {
	if err := a.consumer.Close(); err != nil {
		a.log.Errorf("failed to close consumer: %v", err)

		return fmt.Errorf("failed to close consumer: %w", err)
	}

	a.grpcServer.GracefulStop()

	return nil
}
