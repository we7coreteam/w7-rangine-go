package redis

import (
	"strconv"

	"github.com/go-redis/redis/v8"
)

type RedisFactory struct {
	channelMap map[string]*redis.Client
}

func NewRedisFactory() *RedisFactory {
	return &RedisFactory{
		channelMap: make(map[string]*redis.Client),
	}
}

func (redisFactory *RedisFactory) Channel(channel string) *redis.Client {
	redis, exists := redisFactory.channelMap[channel]
	if !exists {
		panic("redis channel " + channel + " not exists")
	}

	return redis
}

func (redisFactory *RedisFactory) MakeRedis(redisConfig Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     redisConfig.Host + ":" + strconv.Itoa(redisConfig.Port),
		Username: redisConfig.Username,
		Password: redisConfig.Password,
		DB:       redisConfig.Db,
		PoolSize: redisConfig.PoolSize,
	})
}

func (redisFactory *RedisFactory) Register(maps map[string]Config) {
	for key, value := range maps {
		redisFactory.channelMap[key] = redisFactory.MakeRedis(value)
	}
}
