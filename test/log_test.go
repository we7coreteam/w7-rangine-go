package test

import (
	"github.com/we7coreteam/w7-rangine-go/src/core/logger"
	"go.uber.org/zap"
	"os"
	"testing"
)

func TestRegisterLogger(t *testing.T) {
	factory := logger.NewLoggerFactory()
	factory.RegisterDriverResolver("test_file", factory.MakeFileStreamDriver)

	factory.RegisterLogger("test", func() (*zap.Logger, error) {
		driver, _ := factory.MakeDriver(logger.Config{
			Driver: "test_file",
			Path:   "./test.log",
		})
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
