package app

import (
	"context"
	"github.com/Verce11o/resume-view/employee-service/internal/config"
	"github.com/Verce11o/resume-view/employee-service/internal/server"
	mongoLib "github.com/Verce11o/resume-view/shared/db/mongodb"
	postgresLib "github.com/Verce11o/resume-view/shared/db/postgres"
	redisLib "github.com/Verce11o/resume-view/shared/db/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type App struct {
	cfg config.Config
	log *zap.SugaredLogger
	srv *server.Server
}

func New(ctx context.Context, cfg config.Config, log *zap.SugaredLogger) (*App, error) {
	var db *pgxpool.Pool
	var mongo *mongoDriver.Client
	var err error

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
		return nil, err
	}

	redis, err := redisLib.New(ctx, redisLib.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		Database: cfg.Redis.Database,
	})

	if err != nil {
		return nil, err
	}

	srv := server.NewServer(log, db, mongo.Database("employees"), redis, cfg)

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
		return err
	}
	return nil
}

func (a *App) Wait() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func (a *App) Stop(ctx context.Context) error {
	return a.srv.Shutdown(ctx)
}
