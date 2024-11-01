package console

import (
	"github.com/spf13/cobra"
)

type Abstract struct {
	Command
}

func (abstract Abstract) GetDescription() string {
	return ""
}

func (abstract Abstract) Configure(cmd *cobra.Command) {

}
