package config

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogger(env string) error {
	var cfg zap.Config

	if env == "production" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	cfg.EncoderConfig.StacktraceKey = "stacktrace"
	cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}

	logger, err := cfg.Build(zap.AddStacktrace(zap.ErrorLevel))
	if err != nil {
		return err
	}

	Logger = logger
	return nil
}
