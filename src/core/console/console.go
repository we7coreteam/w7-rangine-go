package console

import "github.com/spf13/cobra"

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

func (console *Console) RegisterCommand(command CommandInterface) {
	handler := &cobra.Command{
		Use:   command.GetName(),
		Short: command.GetDescription(),
		Run:   command.Handle,
	}
	command.Configure(handler)

	console.handler.AddCommand(handler)
}

func (console *Console) Run() {
	err := console.handler.Execute()
	if err != nil {
		panic(err)
	}
}
