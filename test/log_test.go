package test

import (
	"github.com/we7coreteam/w7-rangine-go/src/core/logger"
	"github.com/we7coreteam/w7-rangine-go/src/core/logger/config"
	"os"
	"testing"
)

func TestRegisterLogMap(t *testing.T) {
	factory := logger.NewLoggerFactory()
	factory.Register(config.Config{
		Drivers: map[string]config.Driver{
			"test": {
				Driver: "file",
				Path:   "./test.log",
				Level:  -1,
			},
			"test1": {
				Driver: "console",
				Path:   "./test1.log",
				Level:  2,
			},
		},
		Channels: map[string]config.Channel{
			"test": {
				Drivers: []string{
					"test",
					"test1",
				},
			},
		},
	})

	logger, err := factory.Channel("test")
	if err != nil {
		t.Error(err)
	}

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
}
