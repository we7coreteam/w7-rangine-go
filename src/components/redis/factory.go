package redis

import (
	"errors"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"strconv"
	"sync"
)

type Factory struct {
	redisResolverMap map[string]func() redis.Cmdable
	redisMap         map[string]redis.Cmdable
	lock             sync.RWMutex
}

func NewRedisFactory() *Factory {
	return &Factory{
		redisMap:         make(map[string]redis.Cmdable),
		redisResolverMap: make(map[string]func() redis.Cmdable),
	}
}

func (factory *Factory) Channel(channel string) (redis.Cmdable, error) {
	factory.lock.RLock()
	redis, exists := factory.redisMap[channel]
	factory.lock.RUnlock()
	if exists {
		return redis, nil
	}

	factory.lock.Lock()
	defer factory.lock.Unlock()

	redis, exists = factory.redisMap[channel]
	if !exists {
		redisResolver, exists := factory.redisResolverMap[channel]
		if !exists {
			return nil, errors.New("redis channel " + channel + " not exists")
		}

		redis = redisResolver()
		factory.redisMap[channel] = redis
	}

	return redis, nil
}

func (factory *Factory) MakeRedis(config Config) redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + strconv.Itoa(config.Port),
		Username: config.Username,
		Password: config.Password,
		DB:       config.Db,
		PoolSize: config.PoolSize,
	})
}

func (factory *Factory) RegisterRedis(channel string, redisResolver func() redis.Cmdable) {
	factory.redisResolverMap[channel] = redisResolver
}

func (factory *Factory) Register(maps map[string]Config) {
	for key, value := range maps {
		func(channel string, config Config) {
			factory.RegisterRedis(key, func() redis.Cmdable {
				err := binding.Validator.ValidateStruct(value)
				if err != nil {
					if validationErrors, ok := err.(validator.ValidationErrors); ok {
						errStr := "redis config error, channel: " + key + ", fields: "
						for _, e := range validationErrors {
							errStr += e.Field() + ";"
						}
						panic(errStr)
					} else {
						panic(err)
					}
				}

				return factory.MakeRedis(value)
			})
		}(key, value)
	}
}
