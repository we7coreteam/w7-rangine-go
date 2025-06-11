package logger

import (
	"go.uber.org/zap"
)

type FactoryInterface interface {
	Channel(channel string) (*zap.Logger, error)
	RegisterDriver(driver string, resolver func(config Config) (DriverInterface, error))
	RegisterLogger(channel string, loggerResolver func() (*zap.Logger, error))
}
