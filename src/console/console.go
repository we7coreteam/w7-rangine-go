package console

import (
	"github.com/spf13/cobra"
	command "github.com/we7coreteam/w7-rangine-go-support/src/console"
)

type Console struct {
	rootCommand *RootCommand
	handler     *cobra.Command
}

func NewConsole() *Console {
	rootCommand := &RootCommand{}
	handler := &cobra.Command{
		Use:   rootCommand.GetName(),
		Short: rootCommand.GetDescription(),
	}
	rootCommand.Configure(handler)

	return &Console{
		rootCommand: rootCommand,
		handler:     handler,
	}
}

func (console Console) RegisterCommand(cmd command.Command) {
	handler := &cobra.Command{
		Use:   cmd.GetName(),
		Short: cmd.GetDescription(),
		Run: func(curCmd *cobra.Command, args []string) {
			console.rootCommand.Handle(curCmd, args)
			cmd.Handle(curCmd, args)
		},
	}
	cmd.Configure(handler)
	console.handler.AddCommand(handler)
}

func (console Console) Run() {
	_ = console.handler.Execute()
}
