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

	logger, err := cfg.Build()
	if err != nil {
		return err
	}

	Logger = logger
	return nil
}
