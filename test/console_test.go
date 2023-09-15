package test

import (
	"github.com/spf13/cobra"
	"github.com/we7coreteam/w7-rangine-go/src/console"
	"testing"
)

func init() {
	cobra.EnableCommandSorting = false
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
