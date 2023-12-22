package driver

import "go.uber.org/zap/zapcore"

type Driver interface {
	LevelEnable(zapcore.Level) bool
	Write(buffer []byte, ent zapcore.Entry, fields []zapcore.Field) error
	Sync() error
}
