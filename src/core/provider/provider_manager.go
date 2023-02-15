package provider

import (
	"github.com/golobby/container/v3/pkg/container"
	"github.com/spf13/viper"
)

type ProviderManager struct {
	Container container.Container
	Config    *viper.Viper
}

func (providerManager *ProviderManager) RegisterProvider(abstract ProviderInterface) ProviderInterface {
	abstract.SetContainer(providerManager.Container)
	abstract.SetConfig(providerManager.Config)
	return abstract
}
