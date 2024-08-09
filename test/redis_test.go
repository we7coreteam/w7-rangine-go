package test

import (
	"context"
	redis2 "github.com/redis/go-redis/v9"
	"github.com/we7coreteam/w7-rangine-go/v2/src/components/redis"
	"testing"
)

func TestRegisterRedisMap(t *testing.T) {
	factory := redis.NewRedisFactory()

	factory.Register(map[string]redis.Config{
		"test": {
			Host: "127.0.0.1",
			Port: 6379,
			Db:   1,
		},
		"test1": {
			Host: "127.0.0.1",
			Port: 6379,
			Db:   2,
		},
	})
	test, err := factory.Channel("test")
	if err != nil {
		t.Error(err)
	}
	test1, err := factory.Channel("test1")
	if err != nil {
		t.Error(err)
	}

	ctx := context.Background()
	result := test.Set(ctx, "test", "test", 0)
	if result.Err() != nil {
		t.Error(result.Err())
	}
	val := test.Get(ctx, "test")
	if val.Val() != "test" {
		t.Error("redis key test val error")
	}
	val1 := test1.Get(ctx, "test")
	if val1.Val() == "test" {
		t.Error("redis key test val error")
	}

	delResult := test.Del(ctx, "test")
	if delResult.Err() != nil {
		t.Error(delResult.Err())
	}
}

func TestRegisterRedis(t *testing.T) {
	factory := redis.NewRedisFactory()
	factory.RegisterRedis("test", func() (redis2.Cmdable, error) {
		return factory.MakeRedis(redis.Config{
			Host: "127.0.0.1",
			Port: 6379,
		})
	})

	redis, err := factory.Channel("test")
	if err != nil {
		t.Error(err)
	}

	ctx := context.Background()
	result := redis.Set(ctx, "test", "test", 0)
	if result.Err() != nil {
		t.Error(result.Err())
	}
	val := redis.Get(ctx, "test")
	if val.Val() != "test" {
		t.Error("redis key test val error")
	}

	delResult := redis.Del(ctx, "test")
	if delResult.Err() != nil {
		t.Error(delResult.Err())
	}
}
