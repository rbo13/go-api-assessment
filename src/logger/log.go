package logger

import (
	"os"

	"go.uber.org/zap"
)

type Log struct {
	*zap.Logger
}

func New(service string) *Log {
	env := os.Getenv("ENV")
	logger, _ := zap.NewProduction(zap.Fields(
		zap.String("env", env),
		zap.String("service", service),
	))

	if env == "" || env == "development" {
		logger, _ = zap.NewDevelopment()
	}

	return &Log{
		logger,
	}
}
