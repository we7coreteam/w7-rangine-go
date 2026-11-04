package logger

import (
	"github.com/we7coreteam/w7-rangine-go/v2/pkg/support/logger"
	"go.uber.org/multierr"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	enc     zapcore.Encoder
	fields  []zapcore.Field
	drivers []logger.Driver
}

func NewDefaultLogger(enc zapcore.Encoder, drivers []logger.Driver) zapcore.Core {
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
	clone.fields = append(clone.fields, fields...)
	return clone
}

func (c *Logger) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return ce.AddCore(ent, c)
}

func (c *Logger) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	// Stack drivers forward fields to channel cores, not the parent encoder.
	// Include With context in the forwarded fields so every output retains it.
	if len(c.fields) != 0 {
		combined := make([]zapcore.Field, 0, len(c.fields)+len(fields))
		combined = append(combined, c.fields...)
		fields = append(combined, fields...)
	}
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
		fields:  append([]zapcore.Field(nil), c.fields...),
		drivers: c.drivers,
	}
}
