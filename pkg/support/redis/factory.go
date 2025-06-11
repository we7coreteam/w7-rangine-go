package redis

import (
	"github.com/redis/go-redis/v9"
)

type FactoryInterface interface {
	Channel(channel string) (redis.Cmdable, error)
	RegisterRedis(channel string, redisResolver func() (redis.Cmdable, error))
}
