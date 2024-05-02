package main

import (
	"context"
	"github.com/Verce11o/resume-view/employee-service/internal/app"
	"github.com/Verce11o/resume-view/employee-service/internal/config"
	"github.com/Verce11o/resume-view/shared/logger"
	"os"
	"os/signal"
	"syscall"
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
		if err := application.Run(); err != nil {
			log.Errorf("Failed to start application: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := application.Stop(ctx); err != nil {
		log.Errorf("Failed to stop application: %v", err)
	}
}
