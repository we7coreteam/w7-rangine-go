package driver

import (
	"github.com/we7coreteam/w7-rangine-go/src/core/logger/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Console struct {
	Driver

	levelEnabler zapcore.LevelEnabler
	writer       zapcore.WriteSyncer
}

func NewConsoleDriver(config config.Driver) (Driver, error) {
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zapcore.Level(config.Level))

	return &Console{
		levelEnabler: atomicLevel,
		writer:       zapcore.AddSync(os.Stdout),
	}, nil
}

func (c Console) LevelEnable(level zapcore.Level) bool {
	return c.levelEnabler.Enabled(level)
}

func (c Console) Write(buffer []byte, ent zapcore.Entry, fields []zapcore.Field) error {
	_, err := c.writer.Write(buffer)
	return err
}

func (c Console) Sync() error {
	return nil
}
