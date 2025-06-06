package test

import (
	"github.com/we7coreteam/w7-rangine-go/v3/src/core/helper"
	"os"
	"testing"
)

func TestConfigEnvParse(t *testing.T) {
	content := "app: ${TEST_APP_ENV}"
	err := os.Setenv("TEST_APP_ENV", "test")
	if err != nil {
		t.Error(err)
	}

	parseStr := string(helper.ParseConfigContentEnv([]byte(content)))
	if parseStr != "app: test" {
		t.Error("env TEST_APP_ENV parse fail")
	}

	err = os.Unsetenv("TEST_APP_ENV")
	if err != nil {
		t.Error(err)
	}
	parseStr1 := string(helper.ParseConfigContentEnv([]byte(content)))
	if parseStr1 != "app: " {
		t.Error("env TEST_APP_ENV parse fail")
	}

	content = "app: ${TEST_APP_ENV-test1}"
	err = os.Setenv("TEST_APP_ENV", "test2")
	if err != nil {
		t.Error(err)
	}

	parseStr = string(helper.ParseConfigContentEnv([]byte(content)))
	if parseStr != "app: test2" {
		t.Error("env TEST_APP_ENV parse fail")
	}

	err = os.Unsetenv("TEST_APP_ENV")
	if err != nil {
		t.Error(err)
	}
	parseStr1 = string(helper.ParseConfigContentEnv([]byte(content)))
	if parseStr1 != "app: test1" {
		t.Error("env TEST_APP_ENV parse fail")
	}
}
