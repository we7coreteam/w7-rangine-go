package test

import (
	"github.com/spf13/cobra"
	"github.com/we7coreteam/w7-rangine-go/v2/src/console"
	"os"
	"testing"
)

func init() {
	cobra.EnableCommandSorting = false
}

type DefaultArgCommand struct {
	console.Abstract
}

func (c DefaultArgCommand) GetName() string {
	return "default:args"
}

func (c DefaultArgCommand) Configure(cmd *cobra.Command) {
	cmd.Flags().String("test-arg", "", "test arg")
}

func (c DefaultArgCommand) Handle(cmd *cobra.Command, args []string) {
	str, _ := cmd.Flags().GetString("test-arg")
	err := os.Setenv("TEST_COMMAND_ARG", str)
	if err != nil {
		panic(err)
	}
}

func TestAddCommand(t *testing.T) {
	consoleManager := console.NewConsole()

	consoleManager.RegisterCommand(new(console.MakeModuleCommand))
	consoleManager.RegisterCommand(new(console.MakeModelCommand))
	if len(consoleManager.GetHandler().Commands()) != 2 {
		t.Error("command register fail")
	}
	if consoleManager.GetHandler().Commands()[0].Use != "make:module" {
		t.Error("command make module register fail")
	}
	if consoleManager.GetHandler().Commands()[1].Use != "make:model" {
		t.Error("command make model register fail")
	}
}

func TestDefaultArgs(t *testing.T) {
	err := os.Unsetenv("TEST_COMMAND_ENV")
	if err != nil {
		t.Error(err)
	}
	oldEnv := os.Getenv("TEST_COMMAND_ENV")
	os.Args = []string{"cobra.test", "-eTEST_COMMAND_ENV=12"}
	console.NewConsole()

	newEnv := os.Getenv("TEST_COMMAND_ENV")
	if oldEnv == newEnv || newEnv != "12" {
		t.Error("command default args err")
	}

	err = os.Unsetenv("TEST_COMMAND_ENV")
	if err != nil {
		t.Error(err)
	}
	oldEnv = os.Getenv("TEST_COMMAND_ENV")
	os.Args = []string{"text", "default:args", "-eTEST_COMMAND_ENV=14", "--config-file=./console_test.go", "--test-arg=test"}
	consoleManager := console.NewConsole()
	consoleManager.RegisterCommand(&DefaultArgCommand{})
	consoleManager.Run()

	newEnv = os.Getenv("TEST_COMMAND_ENV")
	if oldEnv == newEnv || newEnv != "14" {
		t.Error("command default args err")
	}

	configFile := os.Getenv("RANGINE_CONFIG_FILE")
	if configFile != "./console_test.go" {
		t.Error("command args config-file err")
	}

	argEnv := os.Getenv("TEST_COMMAND_ARG")
	if argEnv != "test" {
		t.Error("command args test-arg err")
	}
}
