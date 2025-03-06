package config

import (
	"fmt"
	"go.uber.org/zap"
	"strings"
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
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}

	logger, err := cfg.Build(zap.AddStacktrace(zap.ErrorLevel))
	if err != nil {
		return err
	}

	Logger = logger
	return nil
}

func QueryLogger(query string, args ...interface{}) {
	query = strings.ReplaceAll(query, "\t\t", "")
	query = strings.ReplaceAll(query, ", ", ",\n")
	query = strings.ReplaceAll(query, " from ", "\nfrom ")
	query = strings.ReplaceAll(query, " where ", "\nwhere ")
	query = strings.ReplaceAll(query, " inner join ", "\ninner join ")
	query = strings.ReplaceAll(query, " left join ", "\nleft join ")
	query = strings.ReplaceAll(query, " right join ", "\nright join ")
	query = strings.ReplaceAll(query, " order by ", "\norder by ")
	query = strings.ReplaceAll(query, " limit ", "\nlimit ")

	Logger.Debug(
		fmt.Sprintf("Executing SQL Query:\n%s", query),
		zap.Any("args", args),
	)
}
