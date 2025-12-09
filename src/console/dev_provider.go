//go:build !w7_rangine_release

package console

import "github.com/we7coreteam/w7-rangine-go/v2/pkg/support/console"

type Provider struct {
}

func (p Provider) Register(console console.Console) {
	console.RegisterCommand(new(MakeModuleCommand))
	console.RegisterCommand(new(MakeProjectCommand))
	console.RegisterCommand(new(MakeModelCommand))
	console.RegisterCommand(new(MakeCmdCommand))
}
