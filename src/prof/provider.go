package prof

import (
	"github.com/we7coreteam/w7-rangine-go/src/core/provider"
	"github.com/we7coreteam/w7-rangine-go/src/core/server"
)

type Provider struct {
	provider.Abstract
}

func (provider *Provider) Register() {
	server.RegisterServer(NewProfServer(provider.GetConfig()))
}
