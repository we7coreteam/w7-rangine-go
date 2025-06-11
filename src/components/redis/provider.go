package redis

import (
	"github.com/golobby/container/v3/pkg/container"
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go/v3/pkg/support/redis"
)

type Provider struct {
}

func (provider Provider) Register(config *viper.Viper, container container.Container) {
	var redisConfigMap map[string]Config
	err := config.UnmarshalKey("redis", &redisConfigMap)
	if err != nil {
		panic(err)
	}

	factory := NewRedisFactory()
	factory.Register(redisConfigMap)

	err = container.NamedSingleton("redis-factory", func() redis.FactoryInterface {
		return factory
	})
	if err != nil {
		panic(err)
	}
}
