package logger

import (
	"github.com/Verce11o/resume-view/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	envProd = "prod"
)

func NewLogger(cfg *config.Config) *zap.SugaredLogger {
	var log *zap.Logger

	if cfg.Env == envProd {
		log, _ = zap.NewProduction()
		defer log.Sync()
		return log.Sugar()
	}

	log, _ = zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.DebugLevel),
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "console",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalColorLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},

		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build()

	defer log.Sync()

	return log.Sugar()
}
