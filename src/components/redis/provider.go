package redis

import (
	"github.com/we7coreteam/w7-rangine-go/src/core/provider"
)

type Provider struct {
	provider.Abstract
}

func (provider *Provider) Register() {
	var redisConfigMap map[string]Config
	err := provider.GetConfig().UnmarshalKey("redis", &redisConfigMap)
	if err != nil {
		panic(err)
	}

	factory := NewRedisFactory()
	factory.Register(redisConfigMap)

	err = provider.GetContainer().NamedSingleton("redis-factory", func() *Factory {
		return factory
	})
	if err != nil {
		panic(err)
	}
}
