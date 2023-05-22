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
	rootCommandHandler := &RootCommand{}
	handler := &cobra.Command{
		Use:              rootCommandHandler.GetName(),
		Short:            rootCommandHandler.GetDescription(),
		PersistentPreRun: rootCommandHandler.Handle,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}
	rootCommandHandler.Configure(handler)

	return &Console{
		handler: handler,
	}
}

func GetConsole() *Console {
	return console
}

func (console Console) GetFlagConfigFile(name string) (string, error) {
	return console.handler.Flags().GetString("config-file")
}

func (console Console) RegisterCommand(command command.Command) {
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
