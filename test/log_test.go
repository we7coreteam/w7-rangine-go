package test

import (
	logSpt "github.com/we7coreteam/w7-rangine-go/pkg/support/logger"
	"github.com/we7coreteam/w7-rangine-go/src/core/logger"
	"os"
	"testing"
)

func TestRegisterLogMap(t *testing.T) {
	factory := logger.NewLoggerFactory()
	factory.Register(map[string]logSpt.Config{
		"test": {
			Driver: "file",
			Path:   "./test.log",
			Level:  "debug",
		},
		"test1": {
			Driver: "console",
			Level:  "error",
		},
		"test2": {
			Driver:   "stack",
			Channels: []string{"test", "test1"},
		},
	})

	logger, err := factory.Channel("test2")
	if err != nil {
		t.Error(err)
	}

	logger.Debug("dsfsdfsdf")
	_, err = os.Stat("./runtime/logs/test.log")
	if err != nil && !os.IsNotExist(err) {
		t.Error(err)
	}

	logger.Error("Test")
	_, err = os.Stat("./runtime/logs/test.log")
	if err != nil && os.IsNotExist(err) {
		t.Error(err)
	}
}
