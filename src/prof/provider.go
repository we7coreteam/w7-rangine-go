package prof

import (
	"github.com/we7coreteam/w7-rangine-go-support/src/facade"
	"github.com/we7coreteam/w7-rangine-go-support/src/provider"
)

type Provider struct {
	provider.Abstract
}

func (provider *Provider) Register() {
	facade.RegisterServer(NewProfServer(facade.GetConfig()))
}
