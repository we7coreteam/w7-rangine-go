package logger

import "go.uber.org/zap/zapcore"

type Driver interface {
	Write(level zapcore.Level, enc zapcore.Encoder, ent zapcore.Entry, fields []zapcore.Field) error
	Sync() error
}
