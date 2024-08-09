package console

import (
	"github.com/spf13/cobra"
	"github.com/we7coreteam/w7-rangine-go/v2/pkg/support/console"
)

type Abstract struct {
	console.Command
}

func (abstract Abstract) GetDescription() string {
	return ""
}

func (abstract Abstract) Configure(cmd *cobra.Command) {

}
