package test

import (
	"github.com/we7coreteam/w7-rangine-go/v2/src/core/logger"
	"github.com/we7coreteam/w7-rangine-go/v2/src/core/logger/config"
	"os"
	"testing"
)

func TestRegisterLogMap(t *testing.T) {
	factory := logger.NewLoggerFactory()
	factory.Register(map[string]config.Config{
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
