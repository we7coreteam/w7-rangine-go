package provider

import (
	"github.com/golobby/container/v3/pkg/container"
)

type ProviderInterface interface {
	Register(options interface{})
	SetContainer(container container.Container)
}
