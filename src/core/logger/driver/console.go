package driver

import (
	"github.com/we7coreteam/w7-rangine-go/pkg/support/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Console struct {
	logger.Driver

	levelEnabler zapcore.LevelEnabler
	writer       zapcore.WriteSyncer
}

func NewConsoleDriver(config logger.Config) (logger.Driver, error) {
	atomicLevel := zap.NewAtomicLevel()
	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		return nil, err
	}
	atomicLevel.SetLevel(level)

	return &Console{
		levelEnabler: atomicLevel,
		writer:       zapcore.AddSync(os.Stdout),
	}, nil
}

func (c Console) Write(level zapcore.Level, enc zapcore.Encoder, ent zapcore.Entry, fields []zapcore.Field) error {
	if !c.levelEnabler.Enabled(level) {
		return nil
	}
	buf, err := enc.EncodeEntry(ent, fields)
	if err != nil {
		return err
	}
	defer buf.Free()

	_, err = c.writer.Write(buf.Bytes())
	return err
}

func (c Console) Sync() error {
	return c.writer.Sync()
}
