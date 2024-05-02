package app

import (
	"context"
	"github.com/Verce11o/resume-view/employee-service/internal/config"
	"github.com/Verce11o/resume-view/employee-service/internal/server"
	postgresLib "github.com/Verce11o/resume-view/shared/db/postgres"
	"go.uber.org/zap"
)

type App struct {
	cfg config.Config
	log *zap.SugaredLogger
	srv *server.Server
}

func New(ctx context.Context, cfg config.Config, log *zap.SugaredLogger) (*App, error) {

	db, err := postgresLib.Run(ctx, postgresLib.Config{
		User:     cfg.Postgres.User,
		Password: cfg.Postgres.Password,
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		Database: cfg.Postgres.Name,
		SSLMode:  cfg.Postgres.SSLMode,
	})
	if err != nil {
		return nil, err
	}

	srv := server.NewServer(log, db, cfg)

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

func (a *App) Stop(ctx context.Context) error {
	return a.srv.Shutdown(ctx)
}
