package logger

import (
	"go.uber.org/zap"
)

type Factory interface {
	Channel(channel string) (*zap.Logger, error)
	RegisterDriver(driver string, resolver func(config Config) (Driver, error))
	RegisterLogger(channel string, loggerResolver func() (*zap.Logger, error))
}
