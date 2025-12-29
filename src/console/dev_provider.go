//go:build !w7_rangine_release

package console

import "github.com/we7coreteam/w7-rangine-go/v3/pkg/support/console"

type Provider struct {
}

func (p Provider) Register(console console.ConsoleInterface) {
	console.RegisterCommand(new(MakeModuleCommand))
	console.RegisterCommand(new(MakeProjectCommand))
	console.RegisterCommand(new(MakeModelCommand))
	console.RegisterCommand(new(MakeCmdCommand))
}
