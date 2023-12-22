package driver

import (
	"errors"
	"github.com/we7coreteam/w7-rangine-go/src/core/helper"
	"github.com/we7coreteam/w7-rangine-go/src/core/logger/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"strings"
)

type File struct {
	Driver

	levelEnabler zapcore.LevelEnabler
	writer       zapcore.WriteSyncer
}

func NewFileDriver(config config.Driver) (Driver, error) {
	fields := helper.ValidateAndGetErrFields(config)
	if len(fields) > 0 {
		return nil, errors.New("log config error, reason: fields: " + strings.Join(fields, ","))
	}

	if config.MaxSize <= 0 {
		config.MaxSize = 2
	}
	if config.MaxDays <= 0 {
		config.MaxDays = 7
	}
	if config.MaxBackups <= 0 {
		config.MaxBackups = 1
	}
	writer := lumberjack.Logger{
		Filename:   "./runtime/logs/" + config.Path,
		MaxSize:    int(config.MaxSize),
		MaxBackups: int(config.MaxBackups),
		MaxAge:     int(config.MaxDays),
		Compress:   false,
	}

	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zapcore.Level(config.Level))

	return &File{
		levelEnabler: atomicLevel,
		writer:       zapcore.AddSync(&writer),
	}, nil
}

func (f File) LevelEnable(level zapcore.Level) bool {
	return f.levelEnabler.Enabled(level)
}

func (f File) Write(buffer []byte, ent zapcore.Entry, fields []zapcore.Field) error {
	_, err := f.writer.Write(buffer)
	return err
}

func (f File) Sync() error {
	return f.writer.Sync()
}
