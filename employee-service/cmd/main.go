package main

import (
	"context"

	"github.com/Verce11o/resume-view/employee-service/internal/app"
	"github.com/Verce11o/resume-view/employee-service/internal/config"
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

	errCh := make(chan error, 1)

	application.Run(errCh)

	application.Wait(errCh)

	if err := application.Stop(ctx); err != nil {
		log.Errorf("Failed to stop application: %v", err)
	}
}
