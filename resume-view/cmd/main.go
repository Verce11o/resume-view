package main

import (
	"context"

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
		log.Errorf("failed to initialize application: %v", err)

		return
	}

	errCh := make(chan error, 1)

	application.Run(ctx, errCh)

	application.Wait(cancel, errCh)

	if err := application.Stop(); err != nil {
		log.Errorf("failed to stop application: %v", err)
	}
}
