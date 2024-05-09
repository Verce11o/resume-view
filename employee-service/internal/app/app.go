package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Verce11o/resume-view/employee-service/internal/config"
	"github.com/Verce11o/resume-view/employee-service/internal/server"
	mongoLib "github.com/Verce11o/resume-view/shared/db/mongodb"
	postgresLib "github.com/Verce11o/resume-view/shared/db/postgres"
	redisLib "github.com/Verce11o/resume-view/shared/db/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type App struct {
	cfg config.Config
	log *zap.SugaredLogger
	srv *server.Server
}

func New(ctx context.Context, cfg config.Config, log *zap.SugaredLogger) (*App, error) {
	var (
		db    *pgxpool.Pool
		mongo *mongoDriver.Client
		err   error
	)

	if strings.ToLower(cfg.MainDatabase) == "postgres" {
		db, err = postgresLib.New(ctx, postgresLib.Config{
			User:     cfg.Postgres.User,
			Password: cfg.Postgres.Password,
			Host:     cfg.Postgres.Host,
			Port:     cfg.Postgres.Port,
			Database: cfg.Postgres.Name,
			SSLMode:  cfg.Postgres.SSLMode,
		})
	}

	if strings.ToLower(cfg.MainDatabase) == "mongo" {
		mongo, err = mongoLib.New(ctx, mongoLib.Config{
			Host:       cfg.MongoDB.Host,
			Port:       cfg.MongoDB.Port,
			User:       cfg.MongoDB.User,
			Password:   cfg.MongoDB.Password,
			Database:   cfg.MongoDB.Name,
			ReplicaSet: cfg.MongoDB.ReplicaSet,
		})
	}

	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	redis, err := redisLib.New(ctx, redisLib.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		Database: cfg.Redis.Database,
	})

	if err != nil {
		return nil, fmt.Errorf("could not connect to redis: %w", err)
	}

	srv := server.NewServer(log, db, mongo, redis, cfg)

	return &App{
		cfg: cfg,
		log: log,
		srv: srv,
	}, nil
}

func (a *App) Run() error {
	a.log.Infof("server starting on port %s...", a.cfg.Server.Port)

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
