package driver

import (
	"errors"
	"github.com/we7coreteam/w7-rangine-go/v2/pkg/support/logger"
	"github.com/we7coreteam/w7-rangine-go/v2/src/core/helper"
	"go.uber.org/zap/zapcore"
)

type Stack struct {
	logger.Driver

	channels       []string
	loggerResolver func(channel string) (zapcore.Core, error)
}

func NewStackDriver(loggerResolver func(channel string) (zapcore.Core, error)) func(driver logger.Config) (logger.Driver, error) {
	return func(config logger.Config) (logger.Driver, error) {
		config.Level = "debug"
		err := helper.ValidateConfig(config)
		if err != nil {
			return nil, errors.New("log config error, reason: " + err.Error())
		}
		if len(config.Channels) == 0 {
			return nil, errors.New("log config error, reason: fields: channels")
		}

		return &Stack{
			channels:       config.Channels,
			loggerResolver: loggerResolver,
		}, nil
	}
}

func (s Stack) Write(level zapcore.Level, enc zapcore.Encoder, ent zapcore.Entry, fields []zapcore.Field) error {
	for _, channel := range s.channels {
		channelLogger, err := s.loggerResolver(channel)
		if err != nil {
			return err
		}

		err = channelLogger.Write(ent, fields)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s Stack) Sync() error {
	for _, channel := range s.channels {
		driver, err := s.loggerResolver(channel)
		if err != nil {
			return err
		}

		err = driver.Sync()
		if err != nil {
			return err
		}
	}

	return nil
}
