package console

import "github.com/spf13/cobra"

type Abstract struct {
	cobra.Command
	Interface
}

func (abstract *Abstract) GetDescription() string {
	return ""
}

func (abstract *Abstract) Configure(command *cobra.Command) {

}
