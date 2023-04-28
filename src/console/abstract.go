package console

import (
	"github.com/spf13/cobra"
	command "github.com/we7coreteam/w7-rangine-go-support/src/console"
)

type Abstract struct {
	command.Command
}

func (abstract Abstract) GetDescription() string {
	return ""
}

func (abstract Abstract) Configure(command *cobra.Command) {

}
