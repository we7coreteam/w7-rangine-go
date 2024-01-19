package logger

import (
	"github.com/we7coreteam/w7-rangine-go/src/core/logger/driver"
	"go.uber.org/multierr"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	enc     zapcore.Encoder
	drivers []driver.Driver
}

func NewDefaultLogger(enc zapcore.Encoder, drivers []driver.Driver) zapcore.Core {
	return &Logger{
		enc:     enc,
		drivers: drivers,
	}
}

func (c *Logger) Level() zapcore.Level {
	return zapcore.DebugLevel
}

func (c *Logger) Enabled(level zapcore.Level) bool {
	return true
}

func (c *Logger) With(fields []zapcore.Field) zapcore.Core {
	clone := c.clone()
	for i := range fields {
		fields[i].AddTo(clone.enc)
	}
	return clone
}

func (c *Logger) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return ce.AddCore(ent, c)
}

func (c *Logger) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	for _, driverHandler := range c.drivers {
		err := driverHandler.Write(ent.Level, c.enc, ent, fields)
		if err != nil {
			return err
		}
	}

	if ent.Level > zapcore.ErrorLevel {
		_ = c.Sync()
	}
	return nil
}

func (c *Logger) Sync() error {
	var syncErr error
	for _, driverHandler := range c.drivers {
		err := multierr.Append(syncErr, driverHandler.Sync())
		if err != nil {
			return err
		}
	}

	return syncErr
}

func (c *Logger) clone() *Logger {
	return &Logger{
		enc:     c.enc.Clone(),
		drivers: c.drivers,
	}
}
