package console

import "github.com/spf13/cobra"

type Console struct {
	Handler *cobra.Command
}

func NewConsole() *Console {
	return &Console{
		Handler: &cobra.Command{Use: "rangine"},
	}
}

func (console *Console) RegisterCommand(command CommandInterface) {
	handler := &cobra.Command{
		Run: command.Handle,
	}
	command.Configure(handler)

	console.Handler.AddCommand(handler)
}

func (console *Console) Run() {
	err := console.Handler.Execute()
	if err != nil {
		panic(err)
	}
}
