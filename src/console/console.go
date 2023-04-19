package console

import (
	"github.com/spf13/cobra"
	command "github.com/we7coreteam/w7-rangine-go-support/src/console"
)

var console = NewConsole()

type Console struct {
	handler *cobra.Command
}

func NewConsole() *Console {
	return &Console{
		handler: &cobra.Command{
			Use: "",
		},
	}
}

func GetConsole() *Console {
	return console
}

func (console *Console) RegisterCommand(command command.Command) {
	handler := &cobra.Command{
		Use:   command.GetName(),
		Short: command.GetDescription(),
		Run:   command.Handle,
	}
	command.Configure(handler)

	console.handler.AddCommand(handler)
}

func (console *Console) Run() {
	console.handler.Execute()
}
