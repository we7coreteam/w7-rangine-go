package logger

import "go.uber.org/zap/zapcore"

type IDriver interface {
	Write(level zapcore.Level, enc zapcore.Encoder, ent zapcore.Entry, fields []zapcore.Field) error
	Sync() error
}
