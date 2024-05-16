package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Verce11o/resume-view/employee-service/internal/config"
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
	cfg config.Config
	log *zap.SugaredLogger
	srv *server.Server
}

func New(ctx context.Context, cfg config.Config, log *zap.SugaredLogger) (*App, error) {
	employeeRepo, positionRepo, err := initRepos(ctx, cfg)

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

	employeeCache := redis.NewEmployeeCache(redisClient)
	positionCache := redis.NewPositionCache(redisClient)

	employeeService := service.NewEmployeeService(log, employeeRepo, employeeCache)
	positionService := service.NewPositionService(log, positionRepo, positionCache)

	srv := server.NewServer(log, employeeService, positionService, cfg)

	return &App{
		cfg: cfg,
		log: log,
		srv: srv,
	}, nil
}

func (a *App) Run() error {
	a.log.Infof("server starting on port %s...", a.cfg.HTTPServer.Port)

	if err := a.srv.Run(a.srv.InitRoutes()); err != nil {
		a.log.Errorf("Error while start server: %v", err)

		return fmt.Errorf("could not start server: %w", err)
	}

	return nil
}

func (a *App) Wait() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func (a *App) Stop(ctx context.Context) error {
	if err := a.srv.Shutdown(ctx); err != nil {
		a.log.Errorf("Error while shutting down server: %v", err)

		return fmt.Errorf("could not stop server: %w", err)
	}

	return nil
}

func initRepos(ctx context.Context, cfg config.Config) (service.EmployeeRepository, service.PositionRepository, error) {
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
			return nil, nil, fmt.Errorf("failed to connect to postgres: %w", err)
		}

		return postgres.NewEmployeeRepository(db), postgres.NewPositionRepository(db), nil

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
			return nil, nil, fmt.Errorf("failed to connect to mongodb: %w", err)
		}

		db := mongo.Database(mongoMainDatabase)

		return mongodb.NewEmployeeRepository(db), mongodb.NewPositionRepository(db), nil

	default:
		return nil, nil, fmt.Errorf("unknown database type: %s", cfg.MainDatabase)
	}
}
