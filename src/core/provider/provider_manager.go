package provider

import (
	"github.com/golobby/container/v3/pkg/container"
	"github.com/spf13/viper"
)

type ProviderManager struct {
	container container.Container
	config    *viper.Viper
}

func NewProviderManager(container container.Container, config *viper.Viper) *ProviderManager {
	return &ProviderManager{
		container: container,
		config:    config,
	}
}

func (providerManager *ProviderManager) RegisterProvider(abstract ProviderInterface) ProviderInterface {
	abstract.SetContainer(providerManager.container)
	abstract.SetConfig(providerManager.config)
	return abstract
}
