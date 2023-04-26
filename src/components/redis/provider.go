package redis

import (
	"github.com/we7coreteam/w7-rangine-go-support/src/facade"
	"github.com/we7coreteam/w7-rangine-go-support/src/provider"
	"github.com/we7coreteam/w7-rangine-go-support/src/redis"
)

type Provider struct {
	provider.Abstract
}

func (provider *Provider) Register() {
	var redisConfigMap map[string]Config
	err := facade.GetConfig().UnmarshalKey("redis", &redisConfigMap)
	if err != nil {
		panic(err)
	}

	factory := NewRedisFactory()
	factory.Register(redisConfigMap)

	err = facade.GetContainer().NamedSingleton("redis-factory", func() redis.Factory {
		return factory
	})
	if err != nil {
		panic(err)
	}
}
