package driver

import (
	"errors"
	"github.com/we7coreteam/w7-rangine-go/pkg/support/logger"
	"github.com/we7coreteam/w7-rangine-go/src/core/helper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"strings"
)

type File struct {
	logger.Driver

	levelEnabler zapcore.LevelEnabler
	writer       zapcore.WriteSyncer
}

func NewFileDriver(config logger.Config) (logger.Driver, error) {
	fields := helper.ValidateAndGetErrFields(config)
	if len(fields) > 0 {
		return nil, errors.New("log config error, reason: fields: " + strings.Join(fields, ","))
	}
	if config.Path == "" {
		return nil, errors.New("log config error, reason: fields: path")
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

	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		return nil, err
	}
	atomicLevel := zap.NewAtomicLevelAt(level)

	return &File{
		levelEnabler: atomicLevel,
		writer:       zapcore.AddSync(&writer),
	}, nil
}

func (f File) Write(level zapcore.Level, enc zapcore.Encoder, ent zapcore.Entry, fields []zapcore.Field) error {
	if !f.levelEnabler.Enabled(level) {
		return nil
	}
	buf, err := enc.EncodeEntry(ent, fields)
	if err != nil {
		return err
	}
	defer buf.Free()

	_, err = f.writer.Write(buf.Bytes())
	return err
}

func (f File) Sync() error {
	return f.writer.Sync()
}
