package main

import (
	"context"
	"github.com/Verce11o/resume-view/internal/app"
	"github.com/Verce11o/resume-view/internal/config"
	"github.com/Verce11o/resume-view/lib/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Load()

	log := logger.NewLogger(cfg)

	application, err := app.New(ctx, cfg, log)

	if err != nil {
		log.Errorf("Failed to initialize application: %v", err)
		return
	}

	go func() {
		if err := application.Run(); err != nil {
			log.Errorf("Failed to start application: %v", err)
		}
	}()

	log.Infof("server starting on port %s...", cfg.Server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	application.Stop()

}
