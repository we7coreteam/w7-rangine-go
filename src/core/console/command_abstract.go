package console

import "github.com/spf13/cobra"

type CommandAbstract struct {
	cobra.Command
	CommandInterface
}

func (commandAbstract *CommandAbstract) GetDescription() string {
	return ""
}

func (commandAbstract *CommandAbstract) Configure(command *cobra.Command) {

}
