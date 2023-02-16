package redis

import (
	"errors"
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

func (redisFactory *RedisFactory) Channel(channel string) (*redis.Client, error) {
	redis, exists := redisFactory.channelMap[channel]
	if !exists {
		return nil, errors.New("redis channel " + channel + " not exists")
	}

	return redis, nil
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

func (redisFactory *RedisFactory) RegisterRedis(channel string, client *redis.Client) {
	redisFactory.channelMap[channel] = client
}

func (redisFactory *RedisFactory) Register(maps map[string]Config) {
	for key, value := range maps {
		redisFactory.RegisterRedis(key, redisFactory.MakeRedis(value))
	}
}
