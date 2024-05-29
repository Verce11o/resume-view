package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Verce11o/resume-view/employee-service/internal/config"
	"github.com/Verce11o/resume-view/employee-service/internal/lib/auth"
	"github.com/Verce11o/resume-view/employee-service/internal/repository/mongodb"
	"github.com/Verce11o/resume-view/employee-service/internal/repository/postgres"
	"github.com/Verce11o/resume-view/employee-service/internal/repository/redis"
	"github.com/Verce11o/resume-view/employee-service/internal/server"
	"github.com/Verce11o/resume-view/employee-service/internal/service"
	mongoLib "github.com/Verce11o/resume-view/shared/db/mongodb"
	postgresLib "github.com/Verce11o/resume-view/shared/db/postgres"
	redisLib "github.com/Verce11o/resume-view/shared/db/redis"
	"go.uber.org/zap"
)

const (
	mainPostgres      = "postgres"
	mainMongodb       = "mongodb"
	mongoMainDatabase = "employees"
)

type App struct {
	cfg     config.Config
	log     *zap.SugaredLogger
	httpSrv *server.HTTP
	grpcSrv *server.GRPC
}

func New(ctx context.Context, cfg config.Config, log *zap.SugaredLogger) (*App, error) {
	employeeRepo, positionRepo, transactor, err := initRepos(ctx, cfg)

	if err != nil {
		return nil, fmt.Errorf("init repos: %w", err)
	}

	redisClient, err := redisLib.New(ctx, redisLib.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		Database: cfg.Redis.Database,
	})

	if err != nil {
		return nil, fmt.Errorf("could not connect to redis: %w", err)
	}

	authenticator := auth.NewAuthenticator(cfg.JWTSignKey, cfg.TokenTTL)

	employeeCache := redis.NewEmployeeCache(redisClient)
	positionCache := redis.NewPositionCache(redisClient)

	employeeService := service.NewEmployeeService(log, employeeRepo, positionRepo, employeeCache,
		transactor, authenticator)
	positionService := service.NewPositionService(log, positionRepo, positionCache)

	httpSrv := server.NewHTTP(log, employeeService, positionService, authenticator, cfg)
	grpcSrv := server.NewGRPC(log, employeeService, positionService, cfg)

	return &App{
		cfg:     cfg,
		log:     log,
		httpSrv: httpSrv,
		grpcSrv: grpcSrv,
	}, nil
}

func (a *App) Run() error {
	a.log.Infof("http server starting on port %s...", a.cfg.HTTPServer.Port)
	a.log.Infof("grpc server starting on port %s...", a.cfg.GRPCServer.Port)

	if err := a.httpSrv.Run(a.httpSrv.InitRoutes()); err != nil {
		a.log.Errorf("error while start http server: %v", err)

		return fmt.Errorf("could not start http server: %w", err)
	}

	if err := a.grpcSrv.Run(); err != nil {
		a.log.Errorf("Error while start grpc server: %v", err)

		return fmt.Errorf("could not start grpc server: %w", err)
	}

	return nil
}

func (a *App) Wait() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func (a *App) Stop(ctx context.Context) error {
	a.grpcSrv.Shutdown()

	if err := a.httpSrv.Shutdown(ctx); err != nil {
		a.log.Errorf("Error while shutting down http server: %v", err)

		return fmt.Errorf("could not stop http server: %w", err)
	}

	return nil
}

func initRepos(ctx context.Context, cfg config.Config) (
	service.EmployeeRepository, service.PositionRepository, service.Transactor, error) {
	switch cfg.MainDatabase {
	case mainPostgres:
		db, err := postgresLib.New(ctx, postgresLib.Config{
			User:     cfg.Postgres.User,
			Password: cfg.Postgres.Password,
			Host:     cfg.Postgres.Host,
			Port:     cfg.Postgres.Port,
			Database: cfg.Postgres.Name,
			SSLMode:  cfg.Postgres.SSLMode,
		})

		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to connect to postgres: %w", err)
		}

		return postgres.NewEmployeeRepository(db), postgres.NewPositionRepository(db), postgres.NewTransactor(db), nil

	case mainMongodb:
		mongo, err := mongoLib.New(ctx, mongoLib.Config{
			Host:       cfg.MongoDB.Host,
			Port:       cfg.MongoDB.Port,
			User:       cfg.MongoDB.User,
			Password:   cfg.MongoDB.Password,
			Database:   cfg.MongoDB.Name,
			ReplicaSet: cfg.MongoDB.ReplicaSet,
		})

		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to connect to mongodb: %w", err)
		}

		db := mongo.Database(mongoMainDatabase)

		return mongodb.NewEmployeeRepository(db), mongodb.NewPositionRepository(db), mongodb.NewTransactor(mongo), nil

	default:
		return nil, nil, nil, fmt.Errorf("unknown database type: %s", cfg.MainDatabase)
	}
}
