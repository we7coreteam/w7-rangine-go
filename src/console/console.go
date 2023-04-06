package console

import (
	"github.com/spf13/cobra"
)

var console = NewConsole()

type Console struct {
	handler *cobra.Command
}

func NewConsole() *Console {
	console := &Console{
		handler: &cobra.Command{
			Use: "",
		},
	}

	return console
}

func GetConsole() *Console {
	return console
}

func (console *Console) RegisterCommand(command Interface) {
	handler := &cobra.Command{
		Use:   command.GetName(),
		Short: command.GetDescription(),
		Run:   command.Handle,
	}
	command.Configure(handler)

	console.handler.AddCommand(handler)
}

func (console *Console) Run() {
	console.RegisterCommand(new(MakeModuleCommand))
	console.RegisterCommand(new(ServerStartCommand))
	console.RegisterCommand(new(VersionCommand))

	err := console.handler.Execute()
	if err != nil {
		panic(err)
	}
}
