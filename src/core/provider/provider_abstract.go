package provider

import (
	"github.com/golobby/container/v3/pkg/container"
	"github.com/spf13/viper"
)

type ProviderAbstract struct {
	ProviderInterface
	container container.Container
	config    interface{}
}

func (providerAbstract *ProviderAbstract) SetContainer(container container.Container) {
	providerAbstract.container = container
}

func (providerAbstract *ProviderAbstract) SetConfig(config *viper.Viper) {
	providerAbstract.config = config
}
