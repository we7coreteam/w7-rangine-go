package logger

import (
	"go.uber.org/zap"
)

type IFactory interface {
	Channel(channel string) (*zap.Logger, error)
	RegisterDriver(driver string, resolver func(config Config) (IDriver, error))
	RegisterLogger(channel string, loggerResolver func() (*zap.Logger, error))
}
