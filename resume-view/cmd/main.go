package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Verce11o/resume-view/resume-view/internal/app"
	"github.com/Verce11o/resume-view/resume-view/internal/config"
	"github.com/Verce11o/resume-view/shared/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Load()

	log := logger.NewLogger(cfg.LogLevel)

	application, err := app.New(ctx, cfg, log)

	if err != nil {
		log.Errorf("Failed to initialize application: %v", err)

		return
	}

	go func() {
		if err := application.Run(ctx); err != nil {
			log.Errorf("Failed to start application: %v", err)
		}
	}()

	log.Infof("server starting on port %s...", cfg.Server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := application.Stop(); err != nil {
		log.Errorf("Failed to stop application: %v", err)
	}
}
