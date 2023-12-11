package test

import (
	"github.com/we7coreteam/w7-rangine-go/src/core/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"log/slog"
	"os"
	"testing"
)

func TestRegisterLogMap(t *testing.T) {
	factory := logger.NewLoggerFactory()
	factory.Register(map[string]logger.Config{
		"test": {
			Driver: "test_file",
			Path:   "./test.log",
			Level:  "debug",
		},
		"test1": {
			Driver: "test_file",
			Path:   "./test1.log",
			Level:  "debug",
		},
	})
	_, err := factory.Channel("test")
	if err == nil {
		t.Error("log channel test driver error")
	}
	_, err = factory.Channel("test1")
	if err == nil {
		t.Error("log channel test1 driver error")
	}
	factory.RegisterDriverResolver("test_file", factory.MakeFileStreamDriver)

	_, err = factory.Channel("test")
	if err != nil {
		t.Error(err)
	}
	_, err = factory.Channel("test1")
	if err != nil {
		t.Error(err)
	}
}

func TestRegisterLogger(t *testing.T) {
	factory := logger.NewLoggerFactory()
	factory.RegisterDriverResolver("test_file", factory.MakeFileStreamDriver)

	factory.RegisterLogger("test", func() (*zap.Logger, error) {
		driver, err := factory.MakeDriver(logger.Config{
			Driver: "test_file",
			Path:   "./test.log",
			Level:  "info",
		})
		if err != nil {
			return nil, err
		}
		return factory.MakeLogger(zap.InfoLevel, driver), nil
	})

	logger, err := factory.Channel("test")
	if err != nil {
		t.Error(err)
	}

	_ = os.Remove("./runtime/logs/test.log")
	logger.Debug("dsfsdfsdf")
	_, err = os.Stat("./runtime/logs/test.log")
	if err == nil || !os.IsNotExist(err) {
		t.Error(err)
	}

	logger.Error("Test")
	_, err = os.Stat("./runtime/logs/test.log")
	if err != nil && os.IsNotExist(err) {
		t.Error(err)
	}
	_ = os.RemoveAll("./runtime")
}

func TestSLogger(t *testing.T) {
	factory := logger.NewLoggerFactory()
	factory.RegisterDriverResolver("test_file", factory.MakeFileStreamDriver)

	factory.RegisterLogger("default", func() (*zap.Logger, error) {
		driver, _ := factory.MakeDriver(logger.Config{
			Driver: "test_file",
			Path:   "./test1.log",
			Level:  "info",
		})
		return factory.MakeLogger(zap.InfoLevel, driver), nil
	})
	defaultLog, err := factory.Channel("default")
	if err == nil {
		defaultSlog := slog.New(zapslog.NewHandler(defaultLog.Core(), nil))
		slog.SetDefault(defaultSlog)
	}

	_ = os.Remove("./runtime/logs/test1.log")
	slog.Debug("dsfsdfsdf")
	_, err = os.Stat("./runtime/logs/test1.log")
	if err == nil || !os.IsNotExist(err) {
		t.Error(err)
	}

	slog.Error("Test")
	_, err = os.Stat("./runtime/logs/test1.log")
	if err != nil && os.IsNotExist(err) {
		t.Error(err)
	}
	_ = os.RemoveAll("./runtime")
}
