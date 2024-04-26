package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewProdLogger() *zap.SugaredLogger {
	log, _ := zap.NewProduction()
	defer log.Sync()
	return log.Sugar()
}

func NewLogger() *zap.SugaredLogger {
	var log *zap.Logger

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
